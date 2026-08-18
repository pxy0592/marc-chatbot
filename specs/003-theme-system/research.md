# Research: 白色与黑色主题系统

## Theme state

**Decision**: React Context 提供 `theme`、`setTheme` 和 `toggleTheme`，Provider 包裹完整应用。

**Rationale**: 登录页和控制台都需要相同状态；集中管理避免 prop drilling 和组件状态不一致。

## Initial preference

**Decision**: 先读取有效 localStorage 值；不存在时读取 `prefers-color-scheme: dark`；异常时回退 light。

**Rationale**: 恢复主动选择，同时为首次访问提供符合设备偏好的默认值。

## DOM integration

**Decision**: 在 `document.documentElement` 设置 `data-theme="light|dark"` 和 `color-scheme`。

**Rationale**: CSS、表单原生控件和二维码可共享主题，且无需组件级颜色判断。

## Visual direction

**Decision**: 抛弃现有深绿、mint、径向渐变、扫描线和霓虹效果，改为黑白中性色基底、蓝色交互强调、12px 圆角、细边框和低对比阴影。

## Alternatives considered

- 单纯添加 dark class 并保留原主题：不满足“抛弃当前主题”。
- 第三方主题库：当前页面规模不需要额外依赖。
- 三态 light/dark/system：用户只要求白色和黑色，第三态增加交互复杂度。
