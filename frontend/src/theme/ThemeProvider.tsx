/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useLayoutEffect, useMemo, useState, type ReactNode } from 'react'

export type Theme = 'light' | 'dark'

const STORAGE_KEY = 'marc-chatbot-theme'

type ThemeContextValue = {
  theme: Theme
  setTheme: (theme: Theme) => void
  toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextValue | null>(null)

function isTheme(value: string | null): value is Theme {
  return value === 'light' || value === 'dark'
}

function getInitialTheme(): Theme {
  try {
    const saved = window.localStorage.getItem(STORAGE_KEY)
    if (isTheme(saved)) return saved
  } catch {
    // Storage can be unavailable in private or restricted browser contexts.
  }

  try {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  } catch {
    // matchMedia is not available in every embedded browser.
  }

  return 'light'
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(getInitialTheme)

  useLayoutEffect(() => {
    document.documentElement.dataset.theme = theme
    document.documentElement.style.colorScheme = theme
    try {
      window.localStorage.setItem(STORAGE_KEY, theme)
    } catch {
      // The in-memory theme still works when persistence is unavailable.
    }
  }, [theme])

  const value = useMemo<ThemeContextValue>(() => ({
    theme,
    setTheme,
    toggleTheme: () => setTheme(current => current === 'light' ? 'dark' : 'light'),
  }), [theme])

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
}

export function useTheme() {
  const context = useContext(ThemeContext)
  if (!context) throw new Error('useTheme must be used within ThemeProvider')
  return context
}
