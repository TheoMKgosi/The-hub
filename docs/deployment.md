# Deployment Guide

This document provides comprehensive deployment instructions for The Hub application, including CI/CD pipelines, containerization, and cloud deployment options.

## 🚀 CI/CD Pipeline

### GitHub Actions Workflow

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy The Hub

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.19'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Backend Tests
        run: |
          cd the-hub-backend
          go mod download
          go test ./tests/... -v -coverprofile=coverage.out

      - name: Frontend Tests
        run: |
          cd the-hub-frontend
          npm install
          npm run test

      - name: Upload coverage reports
        uses: codecov/codecov-action@v3
        with:
          file: ./the-hub-backend/coverage.out

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v4

      - name: Setup Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push backend
        uses: docker/build-push-action@v5
        with:
          context: ./the-hub-backend
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/the-hub-backend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Build and push frontend
        uses: docker/build-push-action@v5
        with:
          context: ./the-hub-frontend
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/the-hub-frontend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Deploy to production
        run: |
          echo "Deploying to production server..."
          # Add your deployment commands here
```

### GitLab CI/CD

Create `.gitlab-ci.yml`:

```yaml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_DRIVER: overlay2

test:backend:
  stage: test
  image: golang:1.19
  before_script:
    - cd the-hub-backend
    - go mod download
  script:
    - go test ./tests/... -v -coverprofile=coverage.out
  coverage: '/coverage: \d+\.\d+%/'

test:frontend:
  stage: test
  image: node:18
  before_script:
    - cd the-hub-frontend
    - npm install
  script:
    - npm run test
  cache:
    key: "$CI_COMMIT_REF_SLUG"
    paths:
      - the-hub-frontend/node_modules/

build:backend:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE/backend:$CI_COMMIT_SHA ./the-hub-backend
    - docker push $CI_REGISTRY_IMAGE/backend:$CI_COMMIT_SHA

build:frontend:
  stage: build
  image: node:18
  before_script:
    - cd the-hub-frontend
    - npm install
  script:
    - npm run build
  artifacts:
    paths:
      - the-hub-frontend/.output/
    expire_in: 1 hour

deploy:staging:
  stage: deploy
  script:
    - echo "Deploying to staging..."
  environment:
    name: staging
    url: https://staging.the-hub.com
  only:
    - develop

deploy:production:
  stage: deploy
  script:
    - echo "Deploying to production..."
  environment:
    name: production
    url: https://the-hub.com
  only:
    - main
  when: manual
```

## 🐳 Docker Deployment

### Backend Dockerfile

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./
EXPOSE 8080
CMD ["./main"]
```

### Frontend Dockerfile

```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/.output/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: the_hub
      POSTGRES_USER: thehub
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U thehub -d the_hub"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build:
      context: ./the-hub-backend
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=thehub
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=the_hub
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build:
      context: ./the-hub-frontend
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - backend
    environment:
      - NUXT_PUBLIC_API_BASE_URL=http://backend:8080/api/v1

volumes:
  postgres_data:
  redis_data:
```

### Environment File

Create `.env`:

```bash
# Database
DB_PASSWORD=your_secure_password_here

# JWT
JWT_SECRET=your_super_secret_jwt_key_here

# Docker Compose
COMPOSE_PROJECT_NAME=the-hub
```

### Deployment Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and restart
docker-compose up -d --build

# Scale services
docker-compose up -d --scale backend=3
```

## ☁️ Cloud Deployment

### AWS Deployment

#### ECS Fargate

1. **Create ECR Repositories:**
```bash
aws ecr create-repository --repository-name the-hub-backend
aws ecr create-repository --repository-name the-hub-frontend
```

2. **Build and Push Images:**
```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account>.dkr.ecr.us-east-1.amazonaws.com

# Tag and push images
docker tag the-hub-backend:latest <account>.dkr.ecr.us-east-1.amazonaws.com/the-hub-backend:latest
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/the-hub-backend:latest
```

3. **ECS Task Definition:**
```json
{
  "family": "the-hub-backend",
  "taskRoleArn": "arn:aws:iam::<account>:role/ecsTaskExecutionRole",
  "executionRoleArn": "arn:aws:iam::<account>:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "backend",
      "image": "<account>.dkr.ecr.us-east-1.amazonaws.com/the-hub-backend:latest",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "DB_HOST", "value": "${DB_HOST}"},
        {"name": "JWT_SECRET", "value": "${JWT_SECRET}"}
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/the-hub-backend",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### API Gateway + Lambda

```yaml
# serverless.yml
service: the-hub-backend

provider:
  name: aws
  runtime: go1.x
  stage: prod
  region: us-east-1

functions:
  api:
    handler: main
    events:
      - http:
          path: /{proxy+}
          method: ANY
          cors: true

resources:
  Resources:
    GatewayResponseDefault4XX:
      Type: AWS::ApiGateway::GatewayResponse
      Properties:
        ResponseType: DEFAULT_4XX
        ResponseParameters:
          gatewayresponse.header.Access-Control-Allow-Origin: "'*'"
          gatewayresponse.header.Access-Control-Allow-Headers: "'*'"
```

### Vercel Deployment (Frontend)

1. **Connect Repository:**
   - Import your GitHub repository to Vercel
   - Configure project settings

2. **Build Configuration:**
```json
{
  "buildCommand": "npm run build",
  "outputDirectory": ".output/public",
  "installCommand": "npm install",
  "framework": "nuxt"
}
```

3. **Environment Variables:**
```
NUXT_PUBLIC_API_BASE_URL=https://your-api-domain.com/api/v1
NUXT_PUBLIC_APP_NAME=The Hub
```

4. **Deploy:**
```bash
vercel --prod
```

### Railway Deployment

1. **Connect Repository:**
   - Link your GitHub repository
   - Configure services

2. **Database Setup:**
   - Add PostgreSQL database
   - Configure environment variables

3. **Deploy Services:**
```bash
# Backend service
railway up --service backend

# Frontend service
railway up --service frontend
```

## 🔧 Manual Deployment

### Backend Deployment

```bash
# On your server
ssh user@your-server

# Install Go
wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone <repository_url>
cd the-hub-backend
go mod download
go build -o the-hub-backend

# Create systemd service
sudo tee /etc/systemd/system/the-hub-backend.service > /dev/null <<EOF
[Unit]
Description=The Hub Backend
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/the-hub-backend
ExecStart=/home/ubuntu/the-hub-backend/the-hub-backend
Restart=always
RestartSec=5
Environment=DB_HOST=localhost
Environment=JWT_SECRET=your_secret

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl enable the-hub-backend
sudo systemctl start the-hub-backend
```

### Frontend Deployment

```bash
# Build static site
cd the-hub-frontend
npm install
npm run generate

# Serve with nginx
sudo cp -r .output/public/* /var/www/html/

# Nginx configuration
sudo tee /etc/nginx/sites-available/the-hub > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/html;
    index index.html;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/the-hub /etc/nginx/sites-enabled/
sudo systemctl reload nginx
```

## 🔒 Security Considerations

### SSL/TLS Setup

```bash
# Let's Encrypt SSL
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com

# Force HTTPS
sudo tee /etc/nginx/sites-available/the-hub > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # ... rest of config
}
EOF
```

### Environment Variables

```bash
# Production environment
DB_HOST=prod-db-host
DB_PASSWORD=secure_password
JWT_SECRET=very_secure_random_string
REDIS_URL=redis://prod-redis-host:6379

```

### Firewall Configuration

```bash
# UFW firewall rules
sudo ufw allow ssh
sudo ufw allow 'Nginx Full'
sudo ufw allow 5432  # PostgreSQL (if needed)
sudo ufw --force enable
```

## 📊 Monitoring & Logging

### Application Monitoring

```bash
# Install Prometheus
sudo apt install prometheus

# Configure prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'the-hub-backend'
    static_configs:
      - targets: ['localhost:8080']
```

### Log Aggregation

```bash
# Install ELK stack or use cloud logging
# AWS CloudWatch, Google Cloud Logging, etc.

# Structured logging configuration
{
  "level": "info",
  "message": "User logged in",
  "user_id": 123,
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "the-hub-backend"
}
```

### Health Checks

```bash
# Application health endpoint
GET /health

# Database health check
GET /health/database

# External service checks
GET /health/external
```

## 🚀 Blue-Green Deployment

```bash
# Create blue environment
docker-compose -f docker-compose.blue.yml up -d

# Test blue environment
curl -f http://blue-environment/health

# Switch traffic to blue
# Update load balancer/nginx config

# Keep green environment as rollback
docker-compose -f docker-compose.green.yml down
```

## 📈 Scaling Strategies

### Horizontal Scaling

```bash
# Docker Swarm
docker swarm init
docker stack deploy -c docker-compose.yml the-hub

# Kubernetes
kubectl apply -f k8s/
kubectl scale deployment the-hub-backend --replicas=5
```

### Load Balancing

```nginx
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    listen 80;
    location /api {
        proxy_pass http://backend;
    }
}
```

## 🔄 Rollback Strategy

```bash
# Quick rollback
docker-compose down
git checkout previous-version-tag
docker-compose up -d --build

# Database rollback
pg_dump prod_db > backup.sql
psql prod_db < backup.sql
```

## 📝 Deployment Checklist

- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] SSL certificates installed
- [ ] Firewall configured
- [ ] Monitoring setup
- [ ] Backup strategy in place
- [ ] Rollback plan ready
- [ ] Team notified of deployment

## 🆘 Troubleshooting

### Common Issues

**Container Won't Start:**
```bash
docker logs <container_id>
docker-compose logs
```

**Database Connection Failed:**
```bash
# Check database connectivity
psql -h localhost -U user -d database

# Check environment variables
docker exec -it <container> env
```

**Application Errors:**
```bash
# Check application logs
docker-compose logs backend

# Health check
curl http://localhost:8080/health
```

### Performance Issues

```bash
# Monitor resource usage
docker stats

# Database performance
EXPLAIN ANALYZE SELECT * FROM users;

# Application profiling
go tool pprof http://localhost:8080/debug/pprof/profile
```

---

**Happy deploying! 🚀**