#!/bin/bash

# DoggyClub Backend - Heroku Deployment Script
# This script deploys the DoggyClub backend to Heroku

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME=${1:-doggyclub-backend}
REGION=${2:-us}
STACK=${3:-heroku-22}

echo -e "${GREEN}🚀 Deploying DoggyClub Backend to Heroku${NC}"
echo "App Name: $APP_NAME"
echo "Region: $REGION"
echo "Stack: $STACK"

# Check if Heroku CLI is installed
if ! command -v heroku &> /dev/null; then
    echo -e "${RED}❌ Heroku CLI is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if user is logged in to Heroku
if ! heroku auth:whoami &> /dev/null; then
    echo -e "${RED}❌ Not logged in to Heroku. Please run 'heroku login' first.${NC}"
    exit 1
fi

# Create Heroku app if it doesn't exist
echo -e "${YELLOW}📱 Creating Heroku app (if it doesn't exist)...${NC}"
heroku apps:create $APP_NAME --region $REGION --stack $STACK || true

# Add PostgreSQL addon
echo -e "${YELLOW}🐘 Adding PostgreSQL addon...${NC}"
heroku addons:create heroku-postgresql:mini -a $APP_NAME || true

# Add Redis addon
echo -e "${YELLOW}🔴 Adding Redis addon...${NC}"
heroku addons:create heroku-redis:mini -a $APP_NAME || true

# Set environment variables
echo -e "${YELLOW}⚙️ Setting environment variables...${NC}"
heroku config:set \
  ENV=production \
  PORT=\$PORT \
  GO_VERSION=1.21 \
  -a $APP_NAME

# Set required config vars (these should be set manually with real values)
echo -e "${YELLOW}📝 Setting placeholder config vars (update these with real values)...${NC}"
heroku config:set \
  JWT_SECRET=change_this_to_a_secure_secret_key \
  CLOUDFLARE_R2_ACCESS_KEY=your_r2_access_key \
  CLOUDFLARE_R2_SECRET_KEY=your_r2_secret_key \
  CLOUDFLARE_R2_BUCKET=doggyclub-files \
  CLOUDFLARE_R2_ENDPOINT=https://your-account.r2.cloudflarestorage.com \
  FIREBASE_PROJECT_ID=your-firebase-project \
  STRIPE_SECRET_KEY=sk_test_your_stripe_key \
  STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret \
  GOOGLE_MAPS_API_KEY=your_google_maps_key \
  ENABLE_ENCOUNTER_DETECTION=true \
  ENABLE_PUSH_NOTIFICATIONS=true \
  ENABLE_PREMIUM_FEATURES=true \
  -a $APP_NAME

# Create Procfile for Heroku
echo -e "${YELLOW}📄 Creating Procfile...${NC}"
cat > Procfile << EOF
web: ./main
release: echo "Running migrations..." && ./main migrate
EOF

# Create app.json for Heroku review apps
echo -e "${YELLOW}📄 Creating app.json...${NC}"
cat > app.json << EOF
{
  "name": "DoggyClub Backend",
  "description": "Backend API for the DoggyClub dog social platform",
  "repository": "https://github.com/your-username/doggyclub-backend",
  "logo": "https://your-domain.com/logo.png",
  "keywords": [
    "go",
    "api",
    "social",
    "dogs",
    "mobile"
  ],
  "image": "heroku/go:1.21",
  "addons": [
    {
      "plan": "heroku-postgresql:mini"
    },
    {
      "plan": "heroku-redis:mini"
    }
  ],
  "env": {
    "ENV": {
      "value": "staging"
    },
    "GO_VERSION": {
      "value": "1.21"
    },
    "JWT_SECRET": {
      "generator": "secret"
    },
    "ENABLE_ENCOUNTER_DETECTION": {
      "value": "true"
    },
    "ENABLE_PUSH_NOTIFICATIONS": {
      "value": "true"
    },
    "ENABLE_PREMIUM_FEATURES": {
      "value": "false"
    }
  },
  "formation": {
    "web": {
      "quantity": 1,
      "size": "basic"
    }
  },
  "buildpacks": [
    {
      "url": "heroku/go"
    }
  ]
}
EOF

# Deploy to Heroku
echo -e "${YELLOW}🚀 Deploying to Heroku...${NC}"
git add Procfile app.json
git commit -m "Add Heroku deployment configuration" || true
heroku git:remote -a $APP_NAME
git push heroku main

# Run migrations
echo -e "${YELLOW}🗃️ Running database migrations...${NC}"
heroku run ./main migrate -a $APP_NAME

# Check deployment status
echo -e "${YELLOW}📊 Checking deployment status...${NC}"
heroku ps -a $APP_NAME

# Show app info
echo -e "${GREEN}✅ Deployment completed!${NC}"
echo ""
echo "App URL: https://$APP_NAME.herokuapp.com"
echo "Admin URL: https://dashboard.heroku.com/apps/$APP_NAME"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Update config vars with real values:"
echo "   heroku config:set JWT_SECRET=your_real_secret -a $APP_NAME"
echo "   heroku config:set CLOUDFLARE_R2_ACCESS_KEY=your_real_key -a $APP_NAME"
echo "   # ... and other secrets"
echo ""
echo "2. Add custom domain (optional):"
echo "   heroku domains:add api.doggyclub.app -a $APP_NAME"
echo ""
echo "3. Set up SSL certificate:"
echo "   heroku certs:auto:enable -a $APP_NAME"
echo ""
echo "4. Monitor logs:"
echo "   heroku logs --tail -a $APP_NAME"
echo ""
echo "5. Scale dynos if needed:"
echo "   heroku ps:scale web=2 -a $APP_NAME"