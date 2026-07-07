# Этап сборки фронтенда
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Этап отдачи статики через Nginx (легкий веб-сервер)
FROM nginx:alpine
# Копируем собранные файлы из этапа builder
COPY --from=builder /app/dist /usr/share/nginx/html
# КОПИРУЕМ КОНФИГ СЮДА, а не в builder!
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]