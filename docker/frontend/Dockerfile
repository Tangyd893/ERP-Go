# Multi-stage SPA image for Vite React frontend apps.
# Build from repo root:
#   docker build -f docker/frontend/Dockerfile --build-arg APP=admin-web-react -t erp-admin-web-react .
#
# Run:
#   docker run --rm -p 8088:80 erp-admin-web-react

ARG APP=admin-web-react
ARG NODE_IMAGE=node:20-alpine
ARG NGINX_IMAGE=nginx:1.27-alpine

FROM ${NODE_IMAGE} AS build
ARG APP
WORKDIR /app

COPY package.json package-lock.json ./
COPY frontend ./frontend

RUN npm ci
RUN npm -w "frontend/${APP}" run build

FROM ${NGINX_IMAGE} AS runtime
ARG APP
COPY docker/frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/frontend/${APP}/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
