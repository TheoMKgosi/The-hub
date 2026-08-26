FROM oven/bun:alpine AS builder
RUN apk add --no-cache libc6-compat
WORKDIR /src

COPY package.json bun.lockb* ./
RUN bun install --frozen-lockfile

COPY . .

ARG API_BASE
ENV API_BASE=$API_BASE
ENV HOST=0.0.0.0
ENV PORT=3000
RUN bun run postinstall
RUN bun run build

FROM node:24-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

COPY --from=builder /src/.output ./.output

EXPOSE 3000
CMD ["HOST=0.0.0.0", "node", ".output/server/index.mjs"]
