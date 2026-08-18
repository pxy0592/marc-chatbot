# Quickstart: Theme validation

```bash
cd frontend
npm test -- --run
npm run lint
npm run build
npm run dev
```

Validate in Chromium:

1. Open login page and switch light/dark.
2. Reload and verify persistence.
3. Log in to mock backend.
4. Switch theme from the dashboard header.
5. Check binding, group, message, form, error, empty and audit states.
6. Check mobile-width layout.
7. Verify no original green gradient or neon styling remains.

## Validation record — 2026-08-18

- ThemeProvider initialization and persistence tests: PASS.
- ThemeToggle accessible interaction test: PASS.
- Full frontend test suite: PASS — 7 tests in 6 files.
- ESLint: PASS.
- TypeScript and Vite production build: PASS.
- Legacy theme scan: PASS — previous mint variables, dark-green backgrounds, radial gradients and scanning patterns are absent.
- Chromium login page: PASS in white and black themes.
- Refresh persistence: PASS — black theme remained selected after reload.
- Chromium authenticated dashboard: PASS in white and black themes.
- Mobile viewport 390x844: PASS — theme, refresh and primary content remained usable without overlap.
- Browser console: no application errors; only development React/Vite informational entries.
