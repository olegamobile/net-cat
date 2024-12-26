#!/bin/bash

PROGRAM_NAME="TCPChat"

echo "Building $PROGRAM_NAME..."
go build -o "$PROGRAM_NAME" main.go

if [ $? -eq 0 ]; then
  echo "Build successful!"
  echo "[USAGE]: ./TCPChat \$port"
else
  echo "Build failed!"
fi