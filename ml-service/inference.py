import os
import tempfile
import shutil
from monai.data import DataLoader, Dataset
import torch
import torch.nn.functional as F
import numpy as np
import nibabel as nib
from PIL import Image
import matplotlib.pyplot as plt
import cv2
from scipy.ndimage import gaussian_filter
from skimage.measure import label, regionprops
from monai.transforms import (
    Compose,
    LoadImaged,
    EnsureChannelFirstd,
    RepeatChanneld,
    Resized,
    NormalizeIntensityd,
    ToTensord
)
from monai.networks.nets import DenseNet121
from nibabel.orientations import aff2axcodes

# === Устройство ===
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# === Соответствие индексов классов ===
idx_to_class = {0: "CN", 1: "MCI", 2: "AD"}

# === Загрузка модели ===
model = DenseNet121(spatial_dims=2, in_channels=3, out_channels=3).to(device)
model.load_state_dict(torch.load("./fullyunfrozen_2.pth", map_location=device))
model.eval()

inference_transforms = Compose([
    LoadImaged(keys=["image"]),
    EnsureChannelFirstd(keys=["image"]),
    RepeatChanneld(keys=["image"], repeats=3),
    Resized(keys=["image"], spatial_size=(224, 224)),
    NormalizeIntensityd(keys=["image"], channel_wise=True),
    ToTensord(keys=["image"]),
])

def extract_slices(nii_path: str, num_slices: int = 15) -> torch.Tensor:

    img = nib.load(nii_path).get_fdata(dtype=np.float32)

    if img.ndim == 2:
        img = img[np.newaxis, ...]
    elif img.ndim == 4:
        img = img[..., 0]

    depth = img.shape[0]
    mid = depth // 2
    half = num_slices // 2
    start = max(0, mid - half)
    end = min(depth, mid + half + 1)

    # === Создаем временную директорию
    temp_dir = tempfile.mkdtemp()

    try:
        # Сохраняем каждый срез как отдельный NIfTI
        data = []
        for idx, i in enumerate(range(start, end)):
            slice_2d = img[i]
            slice_path = os.path.join(temp_dir, f"slice_{idx}.nii.gz")
            nib.save(nib.Nifti1Image(slice_2d, affine=np.eye(4)), slice_path)
            data.append({"image": slice_path})

        dataset = Dataset(data=data, transform=inference_transforms)
        loader = DataLoader(dataset, batch_size=len(data))
        batch = next(iter(loader)) 

        return batch["image"]  # shape: (num_slices, 3, 224, 224)

    finally:
        shutil.rmtree(temp_dir)  # удаляем временные файлы

# === Сохранение промежуточных срезов для отладки ===
def save_debug_slices(slices: torch.Tensor, out_dir: str):
    os.makedirs(out_dir, exist_ok=True)
    for i, sl in enumerate(slices):
        img_arr = sl[0].cpu().numpy()
        img_arr = np.clip(img_arr * 255, 0, 255).astype(np.uint8)
        Image.fromarray(img_arr).save(os.path.join(out_dir, f"slice_{i+1:02d}.png"))

# === Инференс по NIfTI-файлу ===
def run_inference(nii_path: str, save_debug: bool = True) -> tuple[int, float]:
    
    print_nifti_metadata(nii_path)
    check_orientation(nii_path) 
    slices = extract_slices(nii_path).to(device)
    
    if save_debug:
        debug_dir = os.path.join("./debug_slices", os.path.splitext(os.path.basename(nii_path))[0])
        save_debug_slices(slices.cpu(), debug_dir)

    model.eval()
    with torch.no_grad():
        outputs = model(slices)
        print(outputs)
        probs = F.softmax(outputs, dim=1)
        mean_probs = probs.mean(dim=0).cpu().numpy()
        pred_idx = int(np.argmax(mean_probs))
        confidence = float(mean_probs[pred_idx])

    return pred_idx, confidence

def check_orientation(nii_path: str):
    img = nib.load(nii_path)
    axcodes = aff2axcodes(img.affine)
    print(f"Оси снимка: {axcodes}")

    # Простейшая логика: если оси не стандартные, нужно переориентировать
    if axcodes != ('R', 'A', 'S') and axcodes != ('L', 'P', 'S'):
        print("Ориентация нестандартная! Возможно, нужен поворот до аксиального вида.")
    else:
        print("Ориентация в порядке (аксиальный срез).")



def print_nifti_metadata(nii_path: str):
    img = nib.load(nii_path)
    header = img.header

    print(f"=== Метаданные NIfTI файла: {nii_path} ===")
    print(f"Форма данных (shape): {img.shape}")
    print(f"Тип данных (dtype): {img.get_data_dtype()}")
    print(f"Аффинная матрица (affine):\n{img.affine}")
    print(f"Размер вокселей (pixdim): {header.get_zooms()}")
    print(f"Описание (descrip): {header['descrip'].tobytes().decode(errors='ignore').strip()}")
    print(f"Код пространственной ориентации (qform_code): {header['qform_code']}")
    print(f"Код пространственной ориентации (sform_code): {header['sform_code']}")
    print(f"Калибровка максимума (cal_max): {header['cal_max']}")
    print(f"Калибровка минимума (cal_min): {header['cal_min']}")
    print(f"Параметры интенсивности (scl_slope, scl_inter): ({header['scl_slope']}, {header['scl_inter']})")
    print(f"Время сканирования (toffset): {header['toffset']}")
    print(f"Количество временных точек (dim[4]): {header['dim'][4]}")
    print("=" * 50)


def generate_gradcam(
    nii_path: str,
    output_path: str,
    target_layer_name: str = "features.denseblock3",
    threshold: float = 0.3,
    apply_blur: bool = True,
    blur_sigma: float = 2.0
):
    activations = {}
    gradients = {}

    def get_layer(model, layer_path: str):
        parts = layer_path.split('.')
        layer = model
        for part in parts:
            layer = getattr(layer, part)
        return layer

    def forward_hook(module, inp, out):
        activations['value'] = out.detach()

    def backward_hook(module, grad_in, grad_out):
        gradients['value'] = grad_out[0].detach()

    layer = get_layer(model, target_layer_name)
    fh = layer.register_forward_hook(forward_hook)
    bh = layer.register_full_backward_hook(backward_hook)

    slices = extract_slices(nii_path)
    best_score = -float('inf')
    best_slice = None
    best_pred = None

    model.eval()
    for sl in slices:
        sl_batch = sl.unsqueeze(0).to(device)
        out = model(sl_batch)
        prob = F.softmax(out, dim=1)
        score, pred = torch.max(prob, dim=1)
        if score.item() > best_score:
            best_score = score.item()
            best_slice = sl_batch
            best_pred = int(pred)

    model.zero_grad()
    out = model(best_slice)
    target_score = out[0, best_pred]
    target_score.backward()

    act = activations['value'][0]
    grad = gradients['value'][0]
    weights = grad.mean(dim=(1, 2))
    cam = (act * weights[:, None, None]).sum(dim=0).relu()

    cam = cam.unsqueeze(0).unsqueeze(0)
    cam = F.interpolate(cam, size=(224, 224), mode='bilinear', align_corners=False)
    cam = cam.squeeze()
    cam = (cam - cam.min()) / (cam.max() - cam.min() + 1e-8)
    cam_np = cam.cpu().numpy()

    if apply_blur:
        cam_np = gaussian_filter(cam_np, sigma=blur_sigma)
    if threshold > 0:
        cam_np[cam_np < threshold] = 0

    orig = best_slice.cpu().squeeze().numpy().transpose(1, 2, 0)[:, :, 0]

    # === Маска мозга на основе интенсивности ===
    intensity_threshold = np.percentile(orig[orig > 0], 20)  # адаптивный порог
    brain_mask = (orig > intensity_threshold).astype(np.float32)
    cam_np *= brain_mask

    # === Обрезаем по маске мозга (bounding box) ===
    labeled = label(brain_mask)
    props = regionprops(labeled)
    if props:
        y0, x0, y1, x1 = props[0].bbox
        orig = orig[y0:y1, x0:x1]
        cam_np = cam_np[y0:y1, x0:x1]

    # === Сохраняем результат ===
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.figure(figsize=(6, 6))
    plt.imshow(orig, cmap='gray')
    plt.imshow(cam_np, cmap='hot', alpha=0.5)
    plt.axis('off')
    plt.tight_layout()
    plt.savefig(output_path, bbox_inches='tight', pad_inches=0)
    plt.close()

    fh.remove()
    bh.remove()
