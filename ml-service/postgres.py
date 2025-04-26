import psycopg2
from datetime import datetime, timezone
from config import settings

conn = psycopg2.connect(settings.POSTGRES_DSN)

def create_task(scan_id: int):
    with conn.cursor() as cur:
        cur.execute("""
            INSERT INTO mri_analysis_results (scan_id, started_at, status)
            VALUES (%s, %s, 'processing')
        """, (scan_id, datetime.now(timezone.utc)))

        cur.execute("""
            UPDATE mri_scans
            SET status = 'processing'
            WHERE id = %s
        """, (scan_id,))
        
        conn.commit()

def save_result(scan_id: int, diagnosis: int, confidence: float, gradcam_url: str):
    with conn.cursor() as cur:
        cur.execute("""
            UPDATE mri_analysis_results
            SET 
                diagnosis = %s,
                confidence = %s,
                gradcam_url = %s,
                status = 'done',
                completed_at = %s
            WHERE scan_id = %s
        """, (
            diagnosis,
            confidence,
            gradcam_url,
            datetime.now(timezone.utc),
            scan_id
        ))

        cur.execute("""
            UPDATE mri_scans
            SET status = 'done'
            WHERE id = %s
        """, (scan_id,))
        
        conn.commit()

def save_failed(scan_id: int, error_message: str):
    conn.rollback()
    with conn.cursor() as cur:
        cur.execute("""
            UPDATE mri_analysis_results
            SET 
                status = 'failed',
                error_message = %s,
                completed_at = %s
            WHERE scan_id = %s
        """, (
            error_message,
            datetime.now(timezone.utc),
            scan_id
        ))

        cur.execute("""
            UPDATE mri_scans
            SET status = 'failed'
            WHERE id = %s
        """, (scan_id,))
        
        conn.commit()
