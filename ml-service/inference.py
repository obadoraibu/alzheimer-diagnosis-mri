import os
import tempfile
import shutil
import warnings
from typing import Tuple, List

import numpy as np
import torch
import torch.nn.functional as F
from torch.utils.data import DataLoader
import nibabel as nib
from nibabel.orientations import aff2axcodes
from PIL import Image
import matplotlib.pyplot as plt
from scipy.ndimage import gaussian_filter
from skimage.measure import label, regionprops

from monai.data import Dataset
from monai.transforms import (
    Compose, LoadImaged, EnsureChannelFirstd, RepeatChanneld,
    Resized, NormalizeIntensityd, ToTensord
)
from monai.networks.nets import DenseNet121

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

idx_to_class = {0: "CN", 1: "MCI", 2: "AD"}

model = DenseNet121(spatial_dims=2, in_channels=3, out_channels=3).to(device)
model.load_state_dict(torch.load("./model.pth", map_location=device))
model.eval()

inference_transforms = Compose([
    LoadImaged(keys=["image"]),
    EnsureChannelFirstd(keys=["image"]),
    RepeatChanneld(keys=["image"], repeats=3),          # 1 → 3 каналов
    Resized(keys=["image"], spatial_size=(224, 224)),
    NormalizeIntensityd(keys=["image"], channel_wise=True),
    ToTensord(keys=["image"]),
])

def load_ras(nii_path: str) -> Tuple[np.ndarray, np.ndarray]:
    """
    Загружает объём и приводит его к канонической ориентации (RAS).
    Возвращает (data, affine).  data имеет форму (X, Y, Z[, T]).
    """
    img = nib.load(nii_path)
    ras_img = nib.as_closest_canonical(img)  # всегда RAS
    data = ras_img.get_fdata(dtype=np.float32)
    return data, ras_img.affine


def check_orientation(nii_path: str) -> None:
    """
    Выводит исходную и каноническую (RAS) ориентацию для контроля.
    """
    img = nib.load(nii_path)
    orig_axcodes = aff2axcodes(img.affine)

    _, ras_aff = load_ras(nii_path)
    ras_axcodes = aff2axcodes(ras_aff)

    print(f"Оси снимка (исходные): {orig_axcodes}")
    print(f"Оси после RAS       : {ras_axcodes}") 
    if ras_axcodes != ("R", "A", "S"):
        warnings.warn("Не удалось привести к RAS полностью — проверьте данные!")


def print_nifti_metadata(nii_path: str) -> None:
    """
    Полный дамп заголовка NIfTI; полезно для отладки.
    """
    img = nib.load(nii_path)
    header = img.header

    print(f"=== Метаданные NIfTI: {nii_path} ===")
    print(f"Shape                : {img.shape}")
    print(f"Data type            : {img.get_data_dtype()}")
    print(f"Affine               :\n{img.affine}")
    print(f"Voxel size           : {header.get_zooms()}")
    print(f"Description          : {header['descrip'].tobytes().decode(errors='ignore').strip()}")
    print(f"qform_code / sform_code : {header['qform_code']} / {header['sform_code']}")
    print(f"cal_max / cal_min    : {header['cal_max']} / {header['cal_min']}")
    print(f"scl_slope / scl_inter: {header['scl_slope']} / {header['scl_inter']}")
    print(f"Time points (dim4)   : {header['dim'][4]}")
    print("=" * 50)

def extract_slices(nii_path: str, num_slices: int = 15) -> torch.Tensor:
    """
    • Загружает объём, приводит к RAS.  
    • Берёт центральные аксиальные срезы (по оси Z) в количестве `num_slices`.
    • Возвращает батч-тензор формы (num_slices, 3, 224, 224).
    """
    data, _ = load_ras(nii_path)            # (X, Y, Z[, T])

    # Если 4-D (например, DWI): берём первый том
    if data.ndim == 4:
        data = data[..., 0]

    depth = data.shape[2]                   # Z
    half = num_slices // 2
    mid = depth // 2
    start = max(0, mid - half)
    end = min(depth, mid + half + 1)
    chosen = list(range(start, end))

    temp_dir = tempfile.mkdtemp()
    try:
        items: List[dict] = []
        for idx, z in enumerate(chosen):
            slice_2d = data[:, :, z]                     
            slice_path = os.path.join(temp_dir, f"s_{idx}.nii.gz")
            nib.save(nib.Nifti1Image(slice_2d, affine=np.eye(4)), slice_path)
            items.append({"image": slice_path})

        dataset = Dataset(data=items, transform=inference_transforms)
        loader = DataLoader(dataset, batch_size=len(items))
        batch = next(iter(loader))
        return batch["image"]                            # (N,3,224,224)
    finally:
        shutil.rmtree(temp_dir, ignore_errors=True)


def save_debug_slices(slices: torch.Tensor, out_dir: str) -> None:
    """
    Сохраняет преобразованные 224×224 PNG-снимки для визуальной проверки.
    """
    os.makedirs(out_dir, exist_ok=True)
    for i, sl in enumerate(slices):
        img_arr = sl[0].cpu().numpy()                    # 1-й канал
        img_arr = np.clip(img_arr * 255, 0, 255).astype(np.uint8)
        Image.fromarray(img_arr).save(os.path.join(out_dir, f"slice_{i+1:02d}.png"))

def run_inference(nii_path: str, save_debug: bool = False) -> Tuple[int, float]:
    """
    Возвращает (pred_class_idx, confidence).  Под отладкой пишет PNG-срезы.
    """
    print_nifti_metadata(nii_path)
    check_orientation(nii_path)

    slices = extract_slices(nii_path).to(device)

    if save_debug:
        name = os.path.splitext(os.path.basename(nii_path))[0]
        save_debug_slices(slices.cpu(), os.path.join("./debug_slices", name))

    with torch.no_grad():
        outputs = model(slices)                          # (N,3)
        probs = F.softmax(outputs, dim=1)                # (N,3)
        mean_probs = probs.mean(dim=0).cpu().numpy()     # (3,)
        pred_idx = int(np.argmax(mean_probs))
        confidence = float(mean_probs[pred_idx])

    return pred_idx, confidence

def generate_gradcam(
    nii_path: str,
    output_path: str,
    target_layer_name: str = "features.denseblock3",
    threshold: float = 0.3,
    apply_blur: bool = True,
    blur_sigma: float = 2.0,
) -> None:
    """
    Сохраняет Grad-CAM наложение для «лучшего» среза в виде PNG.
    """

    activations, gradients = {}, {}

    def get_layer(model_, layer_path: str):
        layer = model_
        for part in layer_path.split("."):
            layer = getattr(layer, part)
        return layer

    def f_hook(_, __, output):          activations["value"] = output.detach()
    def b_hook(_, grad_in, grad_out):   gradients["value"] = grad_out[0].detach()

    layer = get_layer(model, target_layer_name)
    fh = layer.register_forward_hook(f_hook)
    bh = layer.register_full_backward_hook(b_hook)

    slices = extract_slices(nii_path)                    
    best_score, best_slice, best_pred = -1.0, None, None

    with torch.no_grad():
        for sl in slices:
            slb = sl.unsqueeze(0).to(device)              
            out = model(slb)
            prob = F.softmax(out, dim=1)
            score, pred = torch.max(prob, dim=1)          
            if score.item() > best_score:
                best_score, best_slice, best_pred = score.item(), slb, int(pred)

 
    model.zero_grad(set_to_none=True)
    out = model(best_slice)
    out[0, best_pred].backward()

    act = activations["value"][0]                          # (C,H,W)
    grad = gradients["value"][0]                           # (C,H,W)
    weights = grad.mean(dim=(1, 2))                       # (C,)
    cam = (act * weights[:, None, None]).sum(dim=0).relu()

    cam = F.interpolate(cam.unsqueeze(0).unsqueeze(0), size=(224, 224),
                        mode="bilinear", align_corners=False).squeeze()
    cam_np = cam.cpu().numpy()
    cam_np = (cam_np - cam_np.min()) / (cam_np.ptp() + 1e-8)

    if apply_blur:
        cam_np = gaussian_filter(cam_np, sigma=blur_sigma)
    if threshold > 0:
        cam_np[cam_np < threshold] = 0


    orig = best_slice.cpu().squeeze()[0].numpy()           # (224,224)

    intens_thr = np.percentile(orig[orig > 0], 20)
    brain_mask = (orig > intens_thr).astype(float)
    cam_np *= brain_mask

    labeled = label(brain_mask)
    props = regionprops(labeled)
    if props:
        y0, x0, y1, x1 = props[0].bbox
        orig = orig[y0:y1, x0:x1]
        cam_np = cam_np[y0:y1, x0:x1]


    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.figure(figsize=(6, 6))
    plt.imshow(orig, cmap="gray")
    plt.imshow(cam_np, cmap="hot", alpha=0.5)
    plt.axis("off")
    plt.tight_layout()
    plt.savefig(output_path, bbox_inches="tight", pad_inches=0)
    plt.close()

    fh.remove()
    bh.remove()

