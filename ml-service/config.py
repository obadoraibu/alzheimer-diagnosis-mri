import os
from dotenv import load_dotenv

load_dotenv()

class Settings:
    REDIS_HOST = os.getenv("REDIS_HOST", "localhost")
    REDIS_PORT = int(os.getenv("REDIS_PORT", 6379))

    POSTGRES_HOST = os.getenv("USER_DB_HOST", "localhost")
    POSTGRES_PORT = os.getenv("USER_DB_PORT", "5432")
    POSTGRES_USER = os.getenv("USER_DB_USER", "postgres")
    POSTGRES_PASSWORD = os.getenv("USER_DB_PASSWORD", "postgres")
    POSTGRES_DB = os.getenv("USER_DB_NAME", "mri")

    POSTGRES_DSN = (
        f"host={POSTGRES_HOST} "
        f"port={POSTGRES_PORT} "
        f"user={POSTGRES_USER} "
        f"password={POSTGRES_PASSWORD} "
        f"dbname={POSTGRES_DB}"
    )

    MINIO_ENDPOINT = os.getenv("MINIO_ENDPOINT", "localhost:9000")
    MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
    MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY", "minioadmin")
    MINIO_BUCKET = os.getenv("MINIO_BUCKET", "mri-scans")
    MINIO_SECURE = os.getenv("MINIO_SECURE", "false").lower() == "true"

settings = Settings()