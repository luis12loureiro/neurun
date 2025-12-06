#!/bin/bash

echo "🚀 Starting Neurun Development Environment"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "❌ Docker is not running. Please start Docker first."
  exit 1
fi

echo "✅ Docker is running"
echo ""

# Start Envoy proxy
echo "📡 Starting Envoy gRPC-Web proxy on port 8080..."
docker run --rm -d \
  --name neurun-envoy \
  -p 8080:8080 \
  -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml \
  --add-host=host.docker.internal:host-gateway \
  envoyproxy/envoy:v1.31-latest

if [ $? -eq 0 ]; then
  echo "✅ Envoy proxy started successfully"
else
  echo "❌ Failed to start Envoy proxy"
  exit 1
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Start the Go backend:  cd apps/workflow && make run"
echo "  2. Start the frontend:    cd apps/frontend && npm start"
echo "  3. Open browser:          http://localhost:4200"
echo ""
echo "To stop Envoy proxy: docker stop neurun-envoy"
