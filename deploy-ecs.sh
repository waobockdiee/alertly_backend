#!/bin/bash

# Alertly ECS Deployment Script
# Automatiza el deployment completo a ECS Fargate

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuración
REGION="us-west-2"
STACK_NAME="alertly-ecs-production"
DB_PASSWORD="Po1Ng2O3;"

echo -e "${BLUE}🚀 Starting Alertly ECS Deployment...${NC}"

# 1. Crear el stack de ECS
echo -e "${YELLOW}📦 Step 1: Creating ECS Infrastructure...${NC}"
aws cloudformation deploy \
  --template-file ecs-template.yaml \
  --stack-name $STACK_NAME \
  --region $REGION \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    DBPass="$DB_PASSWORD" \
  --no-fail-on-empty-changeset

# 2. Obtener información del stack
echo -e "${YELLOW}📋 Step 2: Getting stack outputs...${NC}"
ECR_URI=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --region $REGION \
  --query 'Stacks[0].Outputs[?OutputKey==`ECRRepository`].OutputValue' \
  --output text)

LOAD_BALANCER_URL=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --region $REGION \
  --query 'Stacks[0].Outputs[?OutputKey==`LoadBalancerURL`].OutputValue' \
  --output text)

echo -e "${GREEN}✅ ECR Repository: $ECR_URI${NC}"
echo -e "${GREEN}✅ Load Balancer URL: $LOAD_BALANCER_URL${NC}"

# 3. Login a ECR
echo -e "${YELLOW}🔐 Step 3: Logging into ECR...${NC}"
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_URI

# 4. Build y push de la imagen
echo -e "${YELLOW}🔨 Step 4: Building and pushing Docker image...${NC}"
docker build -f Dockerfile.ecs -t alertly-api:latest .
docker tag alertly-api:latest $ECR_URI:latest
docker push $ECR_URI:latest

echo -e "${GREEN}✅ Image pushed successfully!${NC}"

# 5. Forzar nuevo deployment del servicio ECS
echo -e "${YELLOW}🔄 Step 5: Updating ECS service...${NC}"
CLUSTER_NAME=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --region $REGION \
  --query 'Stacks[0].Outputs[?OutputKey==`ClusterName`].OutputValue' \
  --output text)

aws ecs update-service \
  --cluster $CLUSTER_NAME \
  --service alertly-api-service \
  --force-new-deployment \
  --region $REGION > /dev/null

echo -e "${GREEN}✅ Service deployment triggered!${NC}"

# 6. Esperar que el servicio esté estable
echo -e "${YELLOW}⏳ Step 6: Waiting for service to stabilize...${NC}"
aws ecs wait services-stable \
  --cluster $CLUSTER_NAME \
  --services alertly-api-service \
  --region $REGION

# 7. Verificar health check
echo -e "${YELLOW}🔍 Step 7: Verifying API health...${NC}"
sleep 30  # Dar tiempo adicional para que el ALB detecte targets healthy

HEALTH_URL="${LOAD_BALANCER_URL}/health"
echo -e "${BLUE}Testing: $HEALTH_URL${NC}"

for i in {1..10}; do
  if curl -f -s $HEALTH_URL > /dev/null; then
    echo -e "${GREEN}✅ API is healthy and responding!${NC}"
    break
  else
    echo -e "${YELLOW}⏳ Attempt $i/10: Waiting for API to be ready...${NC}"
    sleep 15
  fi
  
  if [ $i -eq 10 ]; then
    echo -e "${RED}❌ API health check failed after 10 attempts${NC}"
    echo -e "${YELLOW}💡 Check ECS service logs for details${NC}"
    exit 1
  fi
done

# 8. Test de performance básico
echo -e "${YELLOW}🏃‍♂️ Step 8: Basic performance test...${NC}"
SIGNUP_URL="${LOAD_BALANCER_URL}/api/account/signup"

echo -e "${BLUE}Testing signup endpoint performance...${NC}"
time curl -X POST "$SIGNUP_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@performance.com",
    "first_name": "Performance",
    "last_name": "Test",
    "password": "testpass123"
  }' \
  -w "\nResponse time: %{time_total}s\n" \
  -s || echo -e "${YELLOW}Note: Signup may fail (expected for existing email)${NC}"

echo ""
echo -e "${GREEN}🎉 DEPLOYMENT COMPLETED SUCCESSFULLY!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🌐 API URL: $LOAD_BALANCER_URL${NC}"
echo -e "${GREEN}🔗 Health Check: $HEALTH_URL${NC}"
echo -e "${GREEN}📊 ECS Console: https://console.aws.amazon.com/ecs/home?region=$REGION#/clusters/$CLUSTER_NAME${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}📝 Next steps:${NC}"
echo -e "${YELLOW}1. Update frontend BASE_URL to: $LOAD_BALANCER_URL${NC}"
echo -e "${YELLOW}2. Test all API endpoints${NC}"
echo -e "${YELLOW}3. Monitor performance and scaling${NC}"
