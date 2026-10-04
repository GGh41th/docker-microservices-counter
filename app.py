import os
import socket
import time
from flask import Flask
import redis

app = Flask(__name__)

# Configuration de la connexion Redis (nom d'hôte db-service requis par le sujet)
REDIS_HOST = os.getenv('REDIS_HOST', 'db-service')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))

# Client Redis
redis_client = redis.Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    decode_responses=True,
    socket_connect_timeout=2
)

def get_hit_count():
    retries = 5
    while True:
        try:
            return redis_client.incr('hits')
        except redis.exceptions.ConnectionError as exc:
            if retries == 0:
                raise exc
            retries -= 1
            time.sleep(0.5)

@app.route('/')
def index():
    try:
        hits = get_hit_count()
    except Exception as e:
        hits = f"[Erreur de connexion Redis: {e}]"
    
    # Récupération du hostname (qui correspond à l'ID court du conteneur dans Docker)
    container_id = socket.gethostname()
    return f"Bonjour ! Cette page a été vue {hits} fois. Je suis le conteneur {container_id}.\n"

if __name__ == '__main__':
    # Écoute sur toutes les interfaces réseau du conteneur
    app.run(host='0.0.0.0', port=5000)
