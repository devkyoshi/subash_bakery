#!/bin/bash
set -e

# Install Docker
apt-get update -y
apt-get install -y ca-certificates curl gnupg lsb-release

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

systemctl enable docker
systemctl start docker

# Create data directories backed by the SSD persistent disk
mkdir -p /data/mongodb /data/redis /data/rabbitmq
mkdir -p /opt/erp-infra

# MongoDB init script — creates the app database and user
cat > /opt/erp-infra/mongo-init.js << EOF
db = db.getSiblingDB('erp_db');
db.createUser({
  user: 'erp_user',
  pwd: '${mongo_password}',
  roles: [{ role: 'readWrite', db: 'erp_db' }]
});
db.createCollection('organizations');
EOF

# docker-compose.yml with injected credentials
cat > /opt/erp-infra/docker-compose.yml << EOF
services:
  mongodb:
    image: mongo:7.0
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: ${mongo_password}
      MONGO_INITDB_DATABASE: erp_db
    volumes:
      - /data/mongodb:/data/db
      - /opt/erp-infra/mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro
    ports:
      - "27017:27017"

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --appendonly yes
    volumes:
      - /data/redis:/data
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3-management-alpine
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: ${rabbitmq_password}
    volumes:
      - /data/rabbitmq:/var/lib/rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
EOF

cd /opt/erp-infra
docker compose up -d

# Systemd service for auto-restart on reboot
cat > /etc/systemd/system/erp-infra.service << 'UNIT'
[Unit]
Description=ERP Infrastructure (MongoDB, Redis, RabbitMQ)
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/erp-infra
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
UNIT

systemctl enable erp-infra.service
