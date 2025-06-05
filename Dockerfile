# Stage 1: Build API
FROM golang:1.24 AS api-builder

WORKDIR /app/api

# Copy API source files
COPY api/ .

# Build the API binary
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o api-server

# Stage 2: Build UI
FROM node:18 AS ui-builder

WORKDIR /app/ui

# Copy UI source files
COPY dashboard/ .

# Create configs directory and .env file
RUN mkdir -p configs && \
    echo "HTTP_PORT=3000\nMETRICS_PORT=2122" > configs/.env

# Install dependencies and build
RUN npm install
RUN npm run build

# Stage 3: Final image
FROM zopdev/static-server:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create directory structure for API
RUN mkdir -p /app/api

# Copy API binary and configs
COPY --from=api-builder /app/api/api-server /app/api/

# Copy UI files
COPY --from=ui-builder /app/ui/dist /website
COPY --from=ui-builder /app/ui/configs/.env /configs/.env

# Create startup script with proper line endings
RUN printf '#!/bin/sh\n\
cd /app/api && ./api-server &\n\
echo "Generating runtime env.js..."\n\
cat <<EOF > /website/env.js\n\
window.env = {\n\
  API_BASE_URL: "${API_BASE_URL}"\n\
};\n\
EOF\n\
exec /main\n\
wait' > /start.sh && \
    chmod +x /start.sh && \
    dos2unix /start.sh || true

COPY ./api/configs/ app/api/configs/
# Expose ports
EXPOSE 8000 3000

# Start both services
CMD ["/bin/sh", "/start.sh"]