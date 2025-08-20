# AGENTS.md

## Build, Lint, and Test Commands

### Frontend (the-hub-frontend)
*   **Build:** `nuxt build`
*   **Development:** `nuxt dev`
*   **Start:** `bun --import ./.output/server/sentry.server.config.mjs .output/server/index.mjs`
*   **Lint:** (Likely integrated with build or editor) - check `nuxt.config.ts` and `package.json` for linting commands.  Consider adding ESLint.
*   **Test:** (Likely using Jest or Vitest) - check Nuxt documentation for testing setup. Run tests with `npm run test` or similar. To run a single test, you may need to adjust the test command with a file path or test name (e.g., `npm run test <test_file_path>`).

### Backend (the-hub-backend)
*   **Build:**  (Go build) - `go build`
*   **Lint:**  (Go lint) - `go vet ./...` or similar.  Check for a linter configuration file.
*   **Test:** `go test ./...`

## Code Style Guidelines

### General
*   Follow existing code style.
*   Use consistent indentation (likely spaces).
*   Use descriptive names for variables and functions.
*   Add comments sparingly, focusing on *why* not *what*.

### Frontend (the-hub-frontend)
*   **Imports:** Follow existing import style.
*   **Formatting:** Use Prettier or similar, configured in `nuxt.config.ts` or `.prettierrc.js`.  Uses Tailwind CSS.
*   **Types:** TypeScript is used. Use types for all variables and function parameters.
*   **Naming Conventions:** Follow Vue.js and Nuxt.js conventions.
*   **Error Handling:** Use `try...catch` blocks for asynchronous operations. Use `useToast` for displaying errors.

### Backend (the-hub-backend)
*   **Imports:** Follow Go import style.
*   **Formatting:** Use `go fmt`.
*   **Types:** Use Go's type system.
*   **Naming Conventions:** Follow Go naming conventions (e.g., `CamelCase` for exported names).
*   **Error Handling:** Handle errors explicitly.

## Additional Notes
*   Check for `.cursor/rules/` or `.github/copilot-instructions.md` for additional guidelines.
*   Do not integrate the schedule system for now
