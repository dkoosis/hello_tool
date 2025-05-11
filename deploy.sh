#!/bin/bash
# direct_deploy.sh - Direct deployment script for hello-tool-base
# Use this if the Makefile/cloudbuild.yaml approach continues to have issues

set -e  # Exit immediately if any command fails

# Configuration
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
SERVICE_NAME="hello-tool-base"
REGION="us-central1"
ARTIFACT_REGISTRY_REPO="my-go-apps"

# Validate project ID
if [ -z "$PROJECT_ID" ]; then
  echo "❌ Error: Google Cloud Project ID not found."
  echo "Set it via 'gcloud config set project YOUR_PROJECT_ID'"
  exit 1
fi

echo "🚀 Deploying $SERVICE_NAME to Google Cloud..."
echo "📋 Project: $PROJECT_ID"
echo "📋 Service: $SERVICE_NAME"
echo "📋 Region: $REGION"
echo "📋 Artifact Registry Repo: $ARTIFACT_REGISTRY_REPO"

# Build the Go application
echo "🔨 Building app for linux/amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
echo "✅ Build successful"

# Build and tag the Docker image
IMAGE_NAME="$REGION-docker.pkg.dev/$PROJECT_ID/$ARTIFACT_REGISTRY_REPO/$SERVICE_NAME:latest"
echo "🔨 Building Docker image: $IMAGE_NAME"
docker build -t "$IMAGE_NAME" .
echo "✅ Docker build successful"

# Push the image to Artifact Registry
echo "📤 Pushing image to Artifact Registry..."
docker push "$IMAGE_NAME"
echo "✅ Image pushed successfully"

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy "$SERVICE_NAME" \
  --image="$IMAGE_NAME" \
  --region="$REGION" \
  --platform=managed \
  --ingress=all \
  --allow-unauthenticated \
  --project="$PROJECT_ID"

echo "✅ Deployment complete!"

# Display service URL
echo "🔗 Service URL: $(gcloud run services describe $SERVICE_NAME --region=$REGION --format='value(status.url)')"