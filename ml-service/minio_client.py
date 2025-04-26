from minio import Minio
from config import settings

client = Minio(
    settings.MINIO_ENDPOINT,
    access_key=settings.MINIO_ACCESS_KEY,
    secret_key=settings.MINIO_SECRET_KEY,
    secure=settings.MINIO_SECURE
)

def download_scan(object_name: str, filepath: str):
    client.fget_object(settings.MINIO_BUCKET, object_name, filepath)

def upload_gradcam(object_name: str, filepath: str):
    client.fput_object(settings.MINIO_BUCKET, object_name, filepath)