# Dockerfile
FROM golang:1.18

# Install any additional dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    make \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Run tests
CMD ["./run_tests.sh"]
