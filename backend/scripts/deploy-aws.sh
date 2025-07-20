#!/bin/bash

# DoggyClub Backend - AWS ECS Deployment Script
# This script deploys the DoggyClub backend to AWS using ECS and ECR

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
AWS_REGION=${AWS_REGION:-us-east-1}
ECR_REPOSITORY=${ECR_REPOSITORY:-doggyclub-backend}
ECS_CLUSTER=${ECS_CLUSTER:-doggyclub-cluster}
ECS_SERVICE=${ECS_SERVICE:-doggyclub-backend-service}
TASK_DEFINITION=${TASK_DEFINITION:-doggyclub-backend-task}
APP_NAME=${APP_NAME:-doggyclub-backend}

echo -e "${GREEN}🚀 Deploying DoggyClub Backend to AWS ECS${NC}"
echo "Region: $AWS_REGION"
echo "ECR Repository: $ECR_REPOSITORY"
echo "ECS Cluster: $ECS_CLUSTER"
echo "ECS Service: $ECS_SERVICE"

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo -e "${RED}❌ AWS CLI is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Get AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URI="$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

echo "AWS Account ID: $ACCOUNT_ID"
echo "ECR URI: $ECR_URI"

# Create ECR repository if it doesn't exist
echo -e "${YELLOW}📦 Creating ECR repository (if it doesn't exist)...${NC}"
aws ecr create-repository \
    --repository-name $ECR_REPOSITORY \
    --region $AWS_REGION \
    --image-scanning-configuration scanOnPush=true || true

# Get ECR login token
echo -e "${YELLOW}🔐 Logging in to ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_URI

# Build Docker image
echo -e "${YELLOW}🔨 Building Docker image...${NC}"
docker build -t $ECR_REPOSITORY:latest .

# Tag image for ECR
echo -e "${YELLOW}🏷️ Tagging image for ECR...${NC}"
docker tag $ECR_REPOSITORY:latest $ECR_URI/$ECR_REPOSITORY:latest
docker tag $ECR_REPOSITORY:latest $ECR_URI/$ECR_REPOSITORY:$(git rev-parse --short HEAD)

# Push image to ECR
echo -e "${YELLOW}📤 Pushing image to ECR...${NC}"
docker push $ECR_URI/$ECR_REPOSITORY:latest
docker push $ECR_URI/$ECR_REPOSITORY:$(git rev-parse --short HEAD)

# Create ECS task definition
echo -e "${YELLOW}📄 Creating ECS task definition...${NC}"
cat > task-definition.json << EOF
{
    "family": "$TASK_DEFINITION",
    "networkMode": "awsvpc",
    "requiresCompatibilities": ["FARGATE"],
    "cpu": "512",
    "memory": "1024",
    "executionRoleArn": "arn:aws:iam::$ACCOUNT_ID:role/ecsTaskExecutionRole",
    "taskRoleArn": "arn:aws:iam::$ACCOUNT_ID:role/ecsTaskRole",
    "containerDefinitions": [
        {
            "name": "$APP_NAME",
            "image": "$ECR_URI/$ECR_REPOSITORY:latest",
            "portMappings": [
                {
                    "containerPort": 8080,
                    "protocol": "tcp"
                }
            ],
            "essential": true,
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/$TASK_DEFINITION",
                    "awslogs-region": "$AWS_REGION",
                    "awslogs-stream-prefix": "ecs"
                }
            },
            "environment": [
                {
                    "name": "ENV",
                    "value": "production"
                },
                {
                    "name": "PORT",
                    "value": "8080"
                }
            ],
            "secrets": [
                {
                    "name": "DB_PASSWORD",
                    "valueFrom": "arn:aws:secretsmanager:$AWS_REGION:$ACCOUNT_ID:secret:doggyclub/db-password"
                },
                {
                    "name": "JWT_SECRET",
                    "valueFrom": "arn:aws:secretsmanager:$AWS_REGION:$ACCOUNT_ID:secret:doggyclub/jwt-secret"
                },
                {
                    "name": "REDIS_PASSWORD",
                    "valueFrom": "arn:aws:secretsmanager:$AWS_REGION:$ACCOUNT_ID:secret:doggyclub/redis-password"
                }
            ],
            "healthCheck": {
                "command": [
                    "CMD-SHELL",
                    "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"
                ],
                "interval": 30,
                "timeout": 5,
                "retries": 3,
                "startPeriod": 60
            }
        }
    ]
}
EOF

# Create CloudWatch log group
echo -e "${YELLOW}📊 Creating CloudWatch log group...${NC}"
aws logs create-log-group \
    --log-group-name "/ecs/$TASK_DEFINITION" \
    --region $AWS_REGION || true

# Register task definition
echo -e "${YELLOW}📝 Registering ECS task definition...${NC}"
aws ecs register-task-definition \
    --cli-input-json file://task-definition.json \
    --region $AWS_REGION

# Create ECS cluster if it doesn't exist
echo -e "${YELLOW}🏗️ Creating ECS cluster (if it doesn't exist)...${NC}"
aws ecs create-cluster \
    --cluster-name $ECS_CLUSTER \
    --capacity-providers FARGATE \
    --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1 \
    --region $AWS_REGION || true

# Create or update ECS service
echo -e "${YELLOW}🔄 Creating/updating ECS service...${NC}"
cat > service-definition.json << EOF
{
    "serviceName": "$ECS_SERVICE",
    "cluster": "$ECS_CLUSTER",
    "taskDefinition": "$TASK_DEFINITION",
    "desiredCount": 2,
    "launchType": "FARGATE",
    "networkConfiguration": {
        "awsvpcConfiguration": {
            "subnets": [
                "subnet-12345678",
                "subnet-87654321"
            ],
            "securityGroups": [
                "sg-12345678"
            ],
            "assignPublicIp": "ENABLED"
        }
    },
    "loadBalancers": [
        {
            "targetGroupArn": "arn:aws:elasticloadbalancing:$AWS_REGION:$ACCOUNT_ID:targetgroup/doggyclub-backend-tg/1234567890123456",
            "containerName": "$APP_NAME",
            "containerPort": 8080
        }
    ],
    "deploymentConfiguration": {
        "maximumPercent": 200,
        "minimumHealthyPercent": 50
    },
    "enableExecuteCommand": true
}
EOF

# Check if service exists
if aws ecs describe-services --cluster $ECS_CLUSTER --services $ECS_SERVICE --region $AWS_REGION --query 'services[0].status' --output text | grep -q "ACTIVE"; then
    echo "Updating existing service..."
    aws ecs update-service \
        --cluster $ECS_CLUSTER \
        --service $ECS_SERVICE \
        --task-definition $TASK_DEFINITION \
        --region $AWS_REGION
else
    echo "Creating new service..."
    echo -e "${RED}⚠️ Note: You need to update the subnet IDs, security group IDs, and target group ARN in service-definition.json${NC}"
    # aws ecs create-service --cli-input-json file://service-definition.json --region $AWS_REGION
fi

# Wait for deployment to complete
echo -e "${YELLOW}⏳ Waiting for deployment to complete...${NC}"
aws ecs wait services-stable \
    --cluster $ECS_CLUSTER \
    --services $ECS_SERVICE \
    --region $AWS_REGION

# Get service status
echo -e "${YELLOW}📊 Checking service status...${NC}"
aws ecs describe-services \
    --cluster $ECS_CLUSTER \
    --services $ECS_SERVICE \
    --region $AWS_REGION \
    --query 'services[0].{Status:status,Running:runningCount,Desired:desiredCount}'

# Clean up temporary files
rm -f task-definition.json service-definition.json

echo -e "${GREEN}✅ Deployment completed!${NC}"
echo ""
echo "ECS Cluster: $ECS_CLUSTER"
echo "ECS Service: $ECS_SERVICE"
echo "ECR Repository: $ECR_URI/$ECR_REPOSITORY"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Set up Application Load Balancer"
echo "2. Configure Route 53 for domain"
echo "3. Set up AWS Secrets Manager for sensitive config"
echo "4. Configure RDS PostgreSQL and ElastiCache Redis"
echo "5. Set up CloudWatch alarms and monitoring"
echo ""
echo "6. Update security groups and subnets in service-definition.json"
echo "7. Set up proper IAM roles for ECS tasks"