# Deployment Guide

This document provides comprehensive deployment instructions for The Hub application, including CI/CD pipelines and cloud deployment options.

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

       - name: Setup Go
         uses: actions/setup-go@v4
         with:
           go-version: '1.24'

       - name: Setup Node.js
         uses: actions/setup-node@v4
         with:
           node-version: '18'
           cache: 'npm'

       - name: Build backend
         run: |
           cd the-hub-backend
           go mod download
           go build -o the-hub-backend

       - name: Build frontend
         run: |
           cd the-hub-frontend
           npm install
           npm run build

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

test:backend:
   stage: test
   image: golang:1.24
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
   image: golang:1.24
   before_script:
     - cd the-hub-backend
     - go mod download
   script:
     - go build -o the-hub-backend
   artifacts:
     paths:
       - the-hub-backend/the-hub-backend
     expire_in: 1 hour

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

## 🗄️ Database Migrations

The application uses golang-migrate for database schema management. Migrations are automatically run during deployment.

#### Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check current migration version
make migrate-version

# Show migration status
make migrate-status
```

#### Migration Safety

- **Always backup** your database before running migrations in production
- Migrations are **idempotent** - safe to run multiple times
- Each migration has both **up** and **down** scripts for rollbacks
- Test migrations on staging environment first

#### Deployment Migration Flow

1. **Pre-deployment**: Create database backup
2. **Migration**: Run `migrate up` to apply schema changes
3. **Verification**: Check migration status and run health checks
4. **Rollback**: If issues occur, use `migrate down` to rollback

## ☁️ Cloud Deployment

### AWS Deployment

#### Elastic Beanstalk

1. **Create Application:**
```bash
aws elasticbeanstalk create-application --application-name the-hub-backend
```

2. **Create Environment:**
```bash
aws elasticbeanstalk create-environment \
  --application-name the-hub-backend \
  --environment-name production \
  --solution-stack-name "64bit Amazon Linux 2 v3.4.4 running Go 1.x"
```

3. **Deploy Application:**
```bash
# Package your application
cd the-hub-backend
go mod download
go build -o application

# Create deployment package
zip deployment.zip application Procfile .ebextensions/

# Deploy
eb deploy production
```

4. **Environment Configuration:**
```yaml
# .ebextensions/environment.config
option_settings:
  aws:elasticbeanstalk:application:environment:
    DB_HOST: "${DB_HOST}"
    JWT_SECRET: "${JWT_SECRET}"
    PORT: 5000
  aws:autoscaling:launchconfiguration:
    InstanceType: t3.micro
    IamInstanceProfile: aws-elasticbeanstalk-ec2-role
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
# Create blue environment (staging)
git checkout -b blue-deployment
# Deploy to blue environment using your deployment method

# Test blue environment
curl -f https://blue.your-domain.com/health

# Switch traffic to blue
# Update DNS/load balancer to point to blue environment

# Keep green environment as rollback
# Blue becomes green, prepare new blue for next deployment
```

## 📈 Scaling Strategies

### Horizontal Scaling

```bash
# AWS Auto Scaling
aws autoscaling update-auto-scaling-group \
  --auto-scaling-group-name the-hub-backend-asg \
  --min-size 2 \
  --max-size 10 \
  --desired-capacity 5

# Kubernetes
kubectl apply -f k8s/
kubectl scale deployment the-hub-backend --replicas=5

# Manual scaling with load balancer
# Add more backend instances and update load balancer configuration
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
git checkout previous-version-tag
# Rebuild and redeploy using your deployment method

# Database rollback
pg_dump prod_db > backup.sql
psql prod_db < backup.sql

# Application rollback
# Stop current version, start previous version
sudo systemctl stop the-hub-backend
sudo systemctl start the-hub-backend-previous
```

## 📝 Deployment Checklist

- [ ] Environment variables configured
- [ ] Database backup created (before migrations)
- [ ] Database migrations run successfully
- [ ] Migration version verified
- [ ] SSL certificates installed
- [ ] Firewall configured
- [ ] Monitoring setup
- [ ] Backup strategy in place
- [ ] Rollback plan ready (including migration rollback)
- [ ] Team notified of deployment

## 🆘 Troubleshooting

### Common Issues

**Application Won't Start:**
```bash
# Check application logs
sudo journalctl -u the-hub-backend -f

# Check systemd status
sudo systemctl status the-hub-backend
```

**Database Connection Failed:**
```bash
# Check database connectivity
psql -h localhost -U user -d database

# Check environment variables
sudo systemctl show the-hub-backend | grep Environment
```

**Application Errors:**
```bash
# Check application logs
sudo journalctl -u the-hub-backend -n 100

# Health check
curl http://localhost:8080/health
```

### Performance Issues

```bash
# Monitor resource usage
top
htop

# Database performance
EXPLAIN ANALYZE SELECT * FROM users;

# Application profiling
go tool pprof http://localhost:8080/debug/pprof/profile
```

---

**Happy deploying! 🚀**