import json
import redis
from config import settings

redis_client = redis.Redis(host=settings.REDIS_HOST, port=settings.REDIS_PORT, decode_responses=True)

def listen_queue(queue_name: str):
    while True:
        _, task = redis_client.blpop(queue_name)
        yield json.loads(task)
