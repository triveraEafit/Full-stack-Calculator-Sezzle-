# ---- Build stage ----
FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# Baked in at build time. Empty string means "same origin" — in the
# docker-compose setup, the browser only ever talks to the frontend
# container; nginx proxies /api/* to the backend container internally,
# so the browser never makes a cross-origin request and CORS doesn't
# come into play in this deployment path.
ENV VITE_API_BASE_URL=""
RUN npm run build

# ---- Run stage ----
FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
