# Prompts Used

This project was built with Claude as an AI pair-programmer, working
incrementally: design first, then one small implementation step at a time,
with review/verification between steps rather than generating everything at
once. Below is the sequence of prompts used, condensed to their essential
content.

1. **Design first, no code.** Provided the full assignment brief and asked
   for a design only — folder structure, API design, request/response
   schemas, backend and frontend architecture, testing strategy, Docker
   strategy, and an implementation order — explicitly optimized for
   maintainability, clean code, and completion within 2–4 hours.

2. **Simplify the design for actual scope.** Asked to strip the initial
   design down: standard library `net/http` instead of a framework, no
   service-layer abstraction, minimal error-handling package, routing kept
   simple, `/operations` endpoint made optional/likely dropped, and a final
   scope decision on which optional operations (power, sqrt, percentage) to
   include based on whether they fit the API contract cleanly.

3. **Scaffold only.** Requested exact terminal commands (PowerShell) to
   create the `backend/` and `frontend/` directory structure and initialize
   Go modules and the Vite React+TypeScript project — no application code
   yet.

4. **Backend: calculator package.** Requested the pure arithmetic package
   (`add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`) with sentinel
   errors for division-by-zero, negative square root, and NaN/Infinity
   results, plus table-driven tests — explicitly no HTTP code, no
   interfaces, standard library only.

5. **Backend: HTTP handler.** Requested the `handler` package wrapping the
   calculator package: JSON decode/validation, operand-count checks,
   mapping calculator errors to `400`/`422` responses with typed error
   codes, plus `httptest`-based table-driven tests.

6. **Backend: small handler fix.** Follow-up request to change the
   unsupported-method response from `400` to `405` with an `Allow: POST`
   header, and update the corresponding test — an isolated, scoped change.

7. **Backend: server wiring.** Requested `main.go` and a small CORS
   `middleware.go`: route registration, graceful shutdown on
   SIGINT/SIGTERM, `PORT` env var support, and CORS for the Vite dev
   server — manual testing only, no test files yet at this step.

8. **Frontend implementation.** Requested the full React/TypeScript
   frontend consuming the existing API: types, API client, a calculator
   state hook, presentational components, CSS styling, and Vitest/RTL
   tests — same simplicity standard as the backend, no state management
   library.

9. **Bug report and fix (state management).** Reported that 6/7 hook tests
   were failing with a consistent pattern (digits overwriting each other,
   wrong operands reaching the API). Asked for the underlying bug to be
   diagnosed and fixed properly rather than changing the tests to match
   broken behavior. This led to replacing the `useState`-based hook with a
   `useReducer` + `useEffect` design.

10. **Visual design pass.** Asked for a significant visual redesign of the
    (functionally complete) UI — clean/modern/professional, no UI library,
    responsive, clear visual hierarchy between number/operator/equals/clear
    buttons, visible focus states — without touching any application logic
    or existing tests.

11. **Bug report and fix (layout height).** Reported the calculator's
    height was changing based on how long the displayed number was, and
    asked for the exact cause and fix. Diagnosed as the display wrapping to
    two lines for long results; fixed with a fixed-height display panel, an
    auto-shrinking font size, and `white-space: nowrap`.

12. **Bug report and fix (layout width).** Reported the calculator's width
    was also changing with input length, and explicitly asked for the exact
    component and CSS rule responsible, with no redesign. Diagnosed as the
    unstyled `#root` element (the real flex child of `<body>`) sizing itself
    to its content's natural width; fixed with one added CSS rule
    constraining `#root`.

Each step's output was reviewed before moving to the next; several steps
(4–8, 11–12) were verified by actually running the resulting tests/build
rather than trusting generated output at face value.
