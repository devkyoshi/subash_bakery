# Stage 1: Build
FROM node:20-alpine as builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# Stage 2: Serve
FROM node:20-alpine

WORKDIR /app

RUN npm install -g serve

COPY --from=builder /app/dist ./dist

# Railway injects PORT automatically
ENV PORT=3000
EXPOSE 3000

CMD ["sh", "-c", "serve -s dist -l $PORT"]
