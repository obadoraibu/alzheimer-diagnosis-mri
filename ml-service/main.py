import os
from redis_queue import listen_queue
from minio_client import download_scan, upload_gradcam
from inference import run_inference, generate_gradcam
from postgres import save_result, create_task, save_failed
from config import settings

SCAN_FOLDER = "./mri_scans"
GRADCAM_FOLDER = "./gradcams"

os.makedirs(SCAN_FOLDER, exist_ok=True)
os.makedirs(GRADCAM_FOLDER, exist_ok=True)

print(f"[Init] MinIO: {settings.MINIO_ENDPOINT}")
print(f"[Init] Redis: {settings.REDIS_PORT}")
print(f"[Init] Postgres DSN: {settings.POSTGRES_DSN}")
print("[Init] Service started and listening for tasks...\n")

for task in listen_queue("mri_tasks"):
    scan_id = task["scan_id"]
    object_name = task["object_name"]

    print(f"\n[Task] Received scan {scan_id} with object: {object_name}")
    create_task(scan_id)
    print(f"[DB] Created task entry for scan_id={scan_id}")

    try:
        downloaded_name = object_name.split("/")[-1]
        scan_path = os.path.join(SCAN_FOLDER, downloaded_name)
        gradcam_path = os.path.join(GRADCAM_FOLDER, f"{scan_id}_gradcam.png")

        print(f"[Download] Downloading scan to {scan_path}...")
        download_scan(object_name, scan_path)

        print(f"[Inference] Running inference on {scan_path}...")
        diagnosis, confidence = run_inference(scan_path)
        print(f"[Inference] Result: diagnosis={diagnosis}, confidence={confidence:.2f}")

        print(f"[GradCAM] Generating GradCAM at {gradcam_path}...")
        generate_gradcam(scan_path, gradcam_path)
        gradcam_object = f"gradcams/{scan_id}.png"

        print(f"[Upload] Uploading GradCAM as {gradcam_object}...")
        upload_gradcam(gradcam_object, gradcam_path)

        print(f"[DB] Saving result to database...")
        save_result(scan_id, diagnosis, confidence, gradcam_object)

        print(f"[Success] Processed scan {scan_id} completely")

    except Exception as e:
        print(f"[Error] Failed to process scan {scan_id}: {e}")
        save_failed(scan_id, str(e))

    finally:
        try:
            os.remove(scan_path)
        except Exception as e:
            print(f"[Warning] Could not delete scan file: {e}")
        try:
            os.remove(gradcam_path)
        except Exception as e:
            print(f"[Warning] Could not delete gradcam file: {e}")


