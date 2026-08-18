# Implementation Plan: 白色与黑色主题系统

**Branch**: `003-theme-system` | **Date**: 2026-08-18 | **Spec**: [spec.md](./spec.md)

## Summary

新增 React ThemeProvider 和 ThemeToggle，使用根节点 `data-theme` 与语义 CSS custom properties 驱动全页面白色/黑色主题。主题偏好保存在浏览器本地存储，无有效保存值时读取系统配色偏好。完全重写现有绿色霓虹样式为中性的现代管理台视觉，同时保持所有业务组件和 API 行为不变。

## Technical Context

**Language/Version**: TypeScript 5.9、React 19、CSS custom properties

**Primary Dependencies**: React Context、localStorage、matchMedia、Testing Library、Vitest

**Storage**: 浏览器 localStorage 键 `marc-chatbot-theme`

**Testing**: Vitest + Testing Library；ESLint；TypeScript/Vite build；Chromium 登录前后双主题视觉检查

**Target Platform**: 现代桌面和移动浏览器

**Project Type**: React 单页管理台 UI 基础设施和视觉重构

**Performance Goals**: 主题切换同步完成，无网络请求；首次渲染不出现明显主题闪烁

**Constraints**: 仅 light/dark；不改变 API、鉴权或业务状态；二维码颜色随主题；移动端可用

## Constitution Check

- 规格与清单已通过。
- 主题状态集中管理，业务组件不维护重复主题状态。
- 主题偏好仅保存枚举值，不含敏感信息。
- 现有组件测试必须保持通过，并新增主题基础设施测试。
- 浏览器验证必须覆盖登录页、控制台、light、dark 和刷新持久化。

Post-design check: 无门禁违反项。

## Project Structure

```text
frontend/src/
├── theme/
│   ├── ThemeProvider.tsx
│   └── ThemeProvider.test.tsx
├── components/
│   ├── ThemeToggle.tsx
│   └── ThemeToggle.test.tsx
├── App.tsx
├── main.tsx
└── index.css
```

**Structure Decision**: 主题状态放在独立 `theme/` 模块；切换控件作为可复用组件；全局样式完全重写并由 `data-theme` 选择语义变量。
