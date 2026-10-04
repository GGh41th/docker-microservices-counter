# TP DevOps : Conteneurisation & Orchestration Micro-services

## Application "Visit-Counter" (Version Haute Performance & Image Minimale)

Ce projet implémente l'application "Visit-Counter" en **Go**, produisant une image Docker finale de **11.9 Mo** (réduction de 85% par rapport aux ~79 Mo initiaux de Python).

---

## 🛠️ Partie 0 : Configuration de l'environnement

1. **Vérifier Docker :**
   ```bash
   docker version
   ```

2. **Connexion au registre :**
   ```bash
   docker login
   ```

3. **Dossier de travail :**
   ```bash
   mkdir TP_Docker_NOM_PRENOM
   cd TP_Docker_NOM_PRENOM
   ```

---

## 📦 Partie 1 : Dockerfile & Optimisation Extrême

### 1. Structure du projet
- [main.go](file:///c:/Users/rouab/Downloads/tp-devops/main.go) : Serveur HTTP en Go, client Redis (`db-service:6379`), gestion du hostname
- [Dockerfile](file:///c:/Users/rouab/Downloads/tp-devops/Dockerfile) : Multi-stage build (compilation statique + compression UPX + Alpine runtime)
- [.dockerignore](file:///c:/Users/rouab/Downloads/tp-devops/.dockerignore) : Filtre de build

### 2. Techniques d'optimisation appliquées
1. **Multi-Stage Build :**
   - Étape `builder` (`golang:alpine`) pour compiler l'application.
   - Image finale basée sur `alpine:3.20` contenant uniquement le binaire exécutable (aucun compilateur, ni SDK, ni code source).
2. **Compilation statique sans debug :**
   - Flags `CGO_ENABLED=0` et `-ldflags="-s -w"` pour éliminer la table des symboles et les informations de débogage DWARF.
3. **Compression binaire UPX :**
   - Réduction de la taille du binaire de ~6.5 Mo à ~2.0 Mo.
4. **Sécurité :**
   - Utilisateur non-root `appuser` (`USER appuser`).

### 3. Résultat de la taille de l'image
```bash
docker images counter-app
```
**Taille obtenue : 11.9 Mo** (contre 79 Mo au départ).

---

## 🚀 Test d'exécution avec Redis

```bash
# 1. Créer un réseau bridge
docker network create counter-net

# 2. Lancer Redis sous le nom attendu par l'app (db-service)
docker run -d --name db-service --network counter-net redis:alpine

# 3. Lancer l'application web
docker run -d --name counter-web --network counter-net -p 5000:5000 counter-app:1.0

# 4. Tester les requêtes
curl http://localhost:5000
# Résultat : "Bonjour ! Cette page a été vue 1 fois. Je suis le conteneur [ID]"
```
