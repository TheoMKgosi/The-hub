FROM oven/bun:alpine AS builder

RUN apk add --no-cache libc6-compat
WORKDIR /src

# Copy package files first to cache the 'bun install' step
COPY package.json bun.lockb* ./
RUN bun install --frozen-lockerfile


# Copy the rest of the project
COPY . .

# Pass public environment variables needed during compilation
ARG NUXT_PUBLIC_API_BASE
ENV NUXT_PUBLIC_API_BASE=$NUXT_PUBLIC_API_BASE

RUN bun run postinstall
RUN bun run build

FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

# Copy ONLY the optimized production engine from the builder stage
COPY --from=builder /src/.output ./.output

EXPOSE 3000

# Run the standalone server natively using Bun
CMD ["node", ".output/server/index.mjs"]
