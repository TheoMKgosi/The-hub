FROM oven/bun:alpine AS builder

RUN apk add --no-cache libc6-compat
WORKDIR /src

# Copy package files first to cache the 'bun install' step
COPY package.json bun.lockb* ./
RUN bun install

# Copy the rest of the project
COPY . .

# Expose the Nuxt dev port
EXPOSE 3000

# Run Nuxt dev and bind to 0.0.0.0
CMD ["bun", "run", "dev", "--host", "0.0.0.0", "--port", "3000"]
