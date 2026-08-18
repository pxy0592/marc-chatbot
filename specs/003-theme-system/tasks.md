# Tasks: 白色与黑色主题系统

## Phase 1: Theme foundation

- [X] T001 Create theme type, initialization, context and persistence in `frontend/src/theme/ThemeProvider.tsx`
- [X] T002 Wrap the application with ThemeProvider in `frontend/src/main.tsx`
- [X] T003 Create accessible reusable toggle in `frontend/src/components/ThemeToggle.tsx`

## Phase 2: Login and dashboard integration

- [X] T004 Add ThemeToggle to `frontend/src/components/AuthGate.tsx`
- [X] T005 Add ThemeToggle to dashboard header in `frontend/src/App.tsx`
- [X] T006 Make QR code foreground inherit theme color in `frontend/src/components/BindingPanel.tsx`

## Phase 3: Full visual replacement

- [X] T007 Replace all root color variables with light and dark semantic tokens in `frontend/src/index.css`
- [X] T008 Replace login page visuals and remove old gradients/scanning effects
- [X] T009 Replace dashboard, cards, navigation and typography visuals
- [X] T010 Replace buttons, inputs, toggles, messages, status and audit visuals
- [X] T011 Validate responsive layouts in both themes

## Phase 4: Tests

- [X] T012 Add initialization and persistence tests in `frontend/src/theme/ThemeProvider.test.tsx`
- [X] T013 Add accessible toggle interaction tests in `frontend/src/components/ThemeToggle.test.tsx`
- [X] T014 Run and preserve all existing component tests

## Phase 5: Documentation and validation

- [X] T015 Update README theme behavior documentation
- [X] T016 Run ESLint, Vitest and production build
- [X] T017 Run Chromium light/dark login and dashboard validation
- [X] T018 Verify refresh persistence and mobile layout
- [X] T019 Record validation results in `specs/003-theme-system/quickstart.md`
