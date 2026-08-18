import { useTheme } from '../theme/ThemeProvider'

function SunIcon() {
  return <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="3.5"/><path d="M12 2v2M12 20v2M4.93 4.93l1.42 1.42M17.65 17.65l1.42 1.42M2 12h2M20 12h2M4.93 19.07l1.42-1.42M17.65 6.35l1.42-1.42"/></svg>
}

function MoonIcon() {
  return <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20.2 15.35A8.5 8.5 0 0 1 8.65 3.8 8.5 8.5 0 1 0 20.2 15.35Z"/></svg>
}

export function ThemeToggle({ compact = false }: { compact?: boolean }) {
  const { theme, toggleTheme } = useTheme()
  const dark = theme === 'dark'
  const target = dark ? '白色' : '黑色'

  return <button
    type="button"
    className={`theme-toggle${compact ? ' compact' : ''}`}
    aria-label={`切换为${target}主题`}
    aria-pressed={dark}
    title={`切换为${target}主题`}
    onClick={toggleTheme}
  >
    <span className="theme-toggle-track" aria-hidden="true">
      <span className="theme-toggle-icon"><SunIcon/></span>
      <span className="theme-toggle-icon"><MoonIcon/></span>
      <span className="theme-toggle-thumb"/>
    </span>
    {!compact && <span className="theme-toggle-label">{dark ? '黑色主题' : '白色主题'}</span>}
  </button>
}
