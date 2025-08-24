# AGENTS.md

## Build, Lint, and Test Commands

### Frontend (the-hub-frontend)
*   **Build:** `nuxt build`
*   **Development:** `nuxt dev`
*   **Start:** `bun --import ./.output/server/sentry.server.config.mjs .output/server/index.mjs`
*   **Lint:** `npx vue-tsc --noEmit` (TypeScript checking)
*   **Test:** `vitest` (watch), `vitest run` (single run), `vitest --ui` (UI mode)
*   **Single Test:** `vitest run <test_file_path>` (e.g., `vitest run test/components/ui/Button.test.ts`)

### Backend (the-hub-backend)
*   **Build:** `go build` or `make build`
*   **Lint:** `go vet ./...` and `go fmt`
*   **Test:** `go test ./tests/... -v` or `make test`
*   **Single Test:** `go test -v ./tests/unit/middleware_test.go`
*   **Test Coverage:** `make test-coverage`
*   **Unit Tests:** `make test-unit`
*   **Integration Tests:** `make test-integration`

## Code Style Guidelines

### Frontend (the-hub-frontend)
*   **Imports:** Group by external libraries first, then internal modules
*   **Formatting:** 2-space indentation, single quotes, semicolons
*   **Types:** Use TypeScript types for all variables and function parameters
*   **Naming:** camelCase for variables/functions, PascalCase for components
*   **Error Handling:** `try...catch` for async operations, `useToast` for user feedback

### Backend (the-hub-backend)
*   **Imports:** Standard library first, then third-party, then internal packages
*   **Formatting:** Use `go fmt` (4-space indentation, no semicolons)
*   **Types:** Use Go's type system with struct tags for JSON/ORM
*   **Naming:** camelCase for unexported, PascalCase for exported identifiers
*   **Error Handling:** Explicit error returns, use `zap` logger for logging

## Additional Notes
*   No Cursor rules or Copilot instructions found
*   Do not integrate the schedule system for now
