#!/bin/bash

# Build script for MultiClip macOS app

echo "Building MultiClip for macOS..."

# Install dependencies
go mod tidy

# Build the application
go build -ldflags="-s -w" -o multiclip main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "Run with: ./multiclip"
    echo ""
    echo "To install system-wide:"
    echo "  sudo cp multiclip /usr/local/bin/"
    echo ""
    echo "To run at startup, add to Login Items in System Preferences"
else
    echo "❌ Build failed!"
    exit 1
fi