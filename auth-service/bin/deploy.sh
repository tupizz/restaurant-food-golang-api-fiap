#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting build and deployment process...${NC}"

# Step 1: Install dependencies if needed
if [ "$1" == "--install" ] || [ "$1" == "-i" ]; then
  echo -e "${YELLOW}Installing dependencies...${NC}"
  npm install
fi

# Step 2: Build TypeScript project
echo -e "${YELLOW}Building TypeScript project...${NC}"
npm run build

# Check if build was successful
if [ $? -ne 0 ]; then
  echo -e "${RED}Build failed! Aborting deployment.${NC}"
  exit 1
fi

# Step 3: Generate Terraform route configurations
echo -e "${YELLOW}Generating Terraform route configurations...${NC}"
npx ts-node ./bin/generete-routes.ts

# Check if route generation was successful
if [ $? -ne 0 ]; then
  echo -e "${RED}Route generation failed! Aborting deployment.${NC}"
  exit 1
fi

# Step 4: Initialize Terraform (if needed)
echo -e "${YELLOW}Initializing Terraform...${NC}"
cd terraform
terraform init

# Step 5: Plan Terraform changes
echo -e "${YELLOW}Planning Terraform changes...${NC}"
terraform plan -out=tfplan

# Step 6: Apply Terraform changes
echo -e "${YELLOW}Applying Terraform changes...${NC}"
terraform apply tfplan

# Check if Terraform apply was successful
if [ $? -ne 0 ]; then
  echo -e "${RED}Deployment failed!${NC}"
  exit 1
fi

# Step 7: Output the API Gateway URL
echo -e "${GREEN}Deployment successful!${NC}"
echo -e "${YELLOW}API Gateway URL:${NC}"
terraform output api_gateway_url

cd ..
echo -e "${GREEN}Build and deployment process completed successfully!${NC}"