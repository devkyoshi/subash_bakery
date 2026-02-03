# Advanced ERP System Backend

An AI-ready, workflow-driven, multi-tenant ERP system built with Golang, MongoDB, and Python.

## Features

- **Multi-Tenant Architecture**: Organization → Company → Location hierarchy
- **Workflow Engine**: Customizable business workflows
- **Dynamic Forms**: JSON Schema-based form builder
- **AI Integration**: Natural language queries, automation, and analytics
- **License Management**: User-based, device-based, and usage-based licensing
- **Subscription Plans**: Flexible plan management
- **Complete ERP Modules**: Products, Inventory, PO, GRN, Suppliers, CRM, Sales, Bills

## Tech Stack

- **Backend**: Golang (Gin framework)
- **Database**: MongoDB
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **AI Services**: Python (FastAPI)
- **Containerization**: Docker + Docker Compose

## Quick Start

### Prerequisites
- Go 1.21+
- Python 3.11+
- Docker & Docker Compose
- MongoDB 7.0+
- Redis 7.0+

### Local Development

1. Clone the repository
```bash
git clone <repo-url>
cd erp-system-backend
```

2. Start infrastructure services
```bash
docker-compose up -d mongodb redis rabbitmq
```

3. Install dependencies
```bash
make install
```

4. Run migrations
```bash
make migrate
```

5. Start services
```bash
make dev
```

### Railway Deployment (Production)

Deploy all services to Railway with automated scripts!

**🚂 Quick Deploy (3 Steps)**
```bash
# 1. Setup (one time)
bash scripts/railway-init.sh

# 2. Configure Railway dashboard (see guide below)

# 3. Deploy
bash scripts/deploy-railway.sh
```

**📚 Documentation**
- [ENVIRONMENT_VARIABLES.md](./ENVIRONMENT_VARIABLES.md) - All env vars

**Features:**
- ✅ Automated deployment scripts
- ✅ Smart change detection (deploy only what changed)
- ✅ Keeps Dockerfiles for local dev
- ✅ Monorepo-friendly setup
- ✅ One-command deployment

## Project Structure

```
erp-system-backend/
├── services/          # Microservices
├── shared/            # Shared libraries
├── infrastructure/    # Docker, K8s configs
├── scripts/           # Utility scripts
└── docs/              # Documentation
```

## Services

- **api-gateway**: API Gateway with routing and auth
- **auth-service**: Authentication & authorization
- **org-service**: Organization hierarchy management
- **license-service**: License management
- **subscription-service**: Subscription plans
- **workflow-service**: Workflow engine
- **form-service**: Dynamic form builder
- **product-service**: Product & inventory management
- **purchase-service**: PO & GRN
- **supplier-service**: Supplier management
- **crm-service**: Customer relationship management
- **sales-service**: Sales & invoices
- **bill-service**: Bill management
- **notification-service**: Multi-channel notifications
- **ai-service**: AI agent (Python)

## API Documentation

Comprehensive API documentation is available in multiple formats:

### Interactive Swagger UI
```bash
# Serve the Swagger UI locally
make swagger-serve

# Then open http://localhost:8000/swagger-ui.html in your browser
```

### Documentation Files
- **OpenAPI Spec**: [`swagger.yaml`](./swagger.yaml) - Complete OpenAPI 3.0 specification
- **API Guide**: [`API_DOCS_README.md`](./API_DOCS_README.md) - Quick start and examples
- **Org Service**: [`services/org-service/API_DOCUMENTATION.md`](./services/org-service/API_DOCUMENTATION.md) - Detailed service docs

### Viewing Options

**Option 1: Local Swagger UI**
```bash
make swagger-serve
# Visit http://localhost:8000/swagger-ui.html
```

**Option 2: Online Swagger Editor**
1. Go to [https://editor.swagger.io/](https://editor.swagger.io/)
2. File > Import File > Select `swagger.yaml`

**Option 3: Postman**
1. Open Postman
2. Import > Select `swagger.yaml`
3. All endpoints available as collection

**Option 4: VS Code**
1. Install "Swagger Viewer" extension
2. Open `swagger.yaml`
3. Right-click > "Preview Swagger"

### Quick Examples

**Register & Login:**
```bash
# Register new user
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Pass123","first_name":"John","last_name":"Doe"}'

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Pass123"}'
```

**Create Organization:**
```bash
curl -X POST http://localhost:8082/api/v1/organizations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Corp","legal_name":"Acme Corporation","domain":"acme.com","email":"info@acme.com","billing_email":"billing@acme.com"}'
```

See [`API_DOCS_README.md`](./API_DOCS_README.md) for complete examples and best practices.

## License

Proprietary
