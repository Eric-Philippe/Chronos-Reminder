# Build stage
FROM oven/bun:latest AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lockb* ./

# Install all dependencies using bun
RUN bun install

# Copy source code
COPY . .

# Build the application with the provided environment variables
RUN bun run build

# Production stage - minimal image
FROM nginx:alpine

WORKDIR /usr/share/nginx/html

# Remove default nginx config
RUN rm -rf ./*

# Copy built application from builder
COPY --from=builder /app/dist ./

# Copy nginx config - rename to .template to prevent envsubst processing
COPY nginx.conf /etc/nginx/conf.d/default.conf.template

# Disable envsubst by removing the template extension handling
RUN mv /etc/nginx/conf.d/default.conf.template /etc/nginx/conf.d/default.conf && \
    chmod 644 /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Health check - using 127.0.0.1 instead of localhost to avoid IPv6 issues
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://127.0.0.1:3000/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
