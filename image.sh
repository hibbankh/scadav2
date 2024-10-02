#! /usr/bin/bash

#AWS Config
AWS_REGION="ap-southeast-1"
AWS_ID="302432638310"

# Set your Docker image name and ECR repository URL
IMAGE_NAME="dhes-scada"
ECR_REPO_URL="$AWS_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

# Login AWS ECR and docker
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Step 1: Calculate how many images are already pushed in the registry
EXISTING_IMAGE_COUNT=$(aws ecr list-images --repository-name $IMAGE_NAME --region $AWS_REGION | jq '.imageIds | length')

echo $EXISTING_IMAGE_COUNT