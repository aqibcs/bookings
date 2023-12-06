#!/bin/bash

echo "Building Go application..."
go build -o bookings cmd/web/*.go

if [ $? -eq 0 ]; then
    echo "Build successful. Running application..."
    ./bookings
else
    echo "Build failed. Please check for errors."
fi
