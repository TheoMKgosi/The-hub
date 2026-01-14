# AGENTS.md

## Build, Lint, and Test Commands

### Frontend (the-hub-frontend)
*   **Build:** `nuxt build`
*   **Dev:** `nuxt dev`
*   **Start:** `bun --import ./.output/server/sentry.server.config.mjs .output/server/index.mjs`
*   **Lint:** `npx vue-tsc --noEmit`
*   **Test:** `vitest` (watch), `vitest run` (single), `vitest --ui` (UI)
*   **Single Test:** `vitest run <path>` (e.g., `vitest run test/components/ui/Button.test.ts`)

### Backend (the-hub-backend)
*   **Build:** `go build` or `make build`
*   **Lint:** `go vet ./...` && `go fmt`
*   **Test:** `go test ./tests/... -v` or `make test`
*   **Single Test:** `go test -v ./tests/unit/middleware_test.go`
*   **Coverage:** `make test-coverage`
*   **Unit/Integration:** `make test-unit` / `make test-integration`

## Code Style Guidelines

### Frontend (Nuxt)
*   **Imports:** External libs → internal modules
*   **Formatting:** 2-space indent, single quotes, semicolons
*   **Types:** TypeScript interfaces for all props/functions
*   **Naming:** camelCase vars/functions, PascalCase components
*   **Error Handling:** try/catch + useToast composable

### Backend (Go + Gin)
*   **Imports:** Standard lib → third-party → internal
*   **Formatting:** `go fmt` (4-space indent, no semicolons)
*   **Types:** Struct tags for JSON/ORM, explicit types
*   **Naming:** camelCase private, PascalCase exported
*   **Error Handling:** Explicit returns, zap structured logging

## Additional Notes
*   No Cursor rules or Copilot instructions found
*   Do not integrate schedule system for now
