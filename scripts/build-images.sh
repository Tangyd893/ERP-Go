# ERP-Go 服务镜像构建
# 用法：./scripts/build-images.sh [tag]

set -euo pipefail

TAG="${1:-latest}"
REPO="${DOCKER_REPO:-erp-go}"
PLATFORM="${DOCKER_PLATFORM:-linux/amd64}"

services=(
  "backend/gateway:gateway:8080"
  "backend/services/iam-service:iam-service:8081"
  "backend/services/tenant-service:tenant-service:8082"
  "backend/services/product-service:product-service:8083"
  "backend/services/channel-service:channel-service:8084"
  "backend/services/order-service:order-service:8085"
  "backend/services/inventory-service:inventory-service:8086"
  "backend/services/warehouse-service:warehouse-service:8087"
  "backend/services/transport-service:transport-service:8088"
  "backend/services/file-service:file-service:8089"
  "backend/services/purchase-service:purchase-service:8091"
  "backend/services/finance-service:finance-service:8092"
  "backend/services/report-service:report-service:8093"
  "backend/services/notification-service:notification-service:8094"
)

for svc in "${services[@]}"; do
  IFS=":" read -r path name port <<< "$svc"
  echo "==> Building ${REPO}/${name}:${TAG}"
  docker build \
    --platform "${PLATFORM}" \
    -f docker/services/Dockerfile \
    --build-arg "SERVICE_PATH=${path}" \
    --build-arg "SERVICE_NAME=${name}" \
    --build-arg "SERVICE_PORT=${port}" \
    -t "${REPO}/${name}:${TAG}" \
    .
done

echo ""
echo "All ${#services[@]} images built with tag ${TAG}"
