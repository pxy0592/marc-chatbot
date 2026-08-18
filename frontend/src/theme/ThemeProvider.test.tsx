import { render, screen } from '@testing-library/react'
import { ThemeProvider, useTheme } from './ThemeProvider'

function ThemeProbe() {
  const { theme } = useTheme()
  return <span>{theme}</span>
}

function setSystemDark(matches: boolean) {
  Object.defineProperty(window, 'matchMedia', {
    configurable: true,
    value: vi.fn().mockReturnValue({ matches, addEventListener: vi.fn(), removeEventListener: vi.fn() }),
  })
}

beforeEach(() => {
  localStorage.clear()
  delete document.documentElement.dataset.theme
  document.documentElement.style.colorScheme = ''
  setSystemDark(false)
})

test('restores a saved dark theme and applies it to the document root', () => {
  localStorage.setItem('marc-chatbot-theme', 'dark')
  render(<ThemeProvider><ThemeProbe/></ThemeProvider>)

  expect(screen.getByText('dark')).toBeInTheDocument()
  expect(document.documentElement).toHaveAttribute('data-theme', 'dark')
  expect(document.documentElement.style.colorScheme).toBe('dark')
})

test('uses the system preference when no valid theme was saved', () => {
  localStorage.setItem('marc-chatbot-theme', 'unsupported')
  setSystemDark(true)
  render(<ThemeProvider><ThemeProbe/></ThemeProvider>)

  expect(screen.getByText('dark')).toBeInTheDocument()
  expect(localStorage.getItem('marc-chatbot-theme')).toBe('dark')
})
