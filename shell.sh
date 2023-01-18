# Build image dari main workspace
docker build -t svc-product:1.0 -f ./app/svc-product/dockerfiles/Dockerfile ./app/svc-product

# Checking stage builder, untuk image name pastikan menggunakan image id terbaru
docker create -it --rm --name stage-builder 11f06801523c
docker start stage-builder
docker exec -it stage-builder sh
docker stop stage-builder

# Running service
docker create --name svc-product -p 8002:8002 --env-file ./app/svc-product/config/.env svc-product:1.0