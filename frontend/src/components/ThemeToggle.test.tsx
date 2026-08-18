import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ThemeToggle } from './ThemeToggle'
import { ThemeProvider } from '../theme/ThemeProvider'

beforeEach(() => {
  localStorage.clear()
  Object.defineProperty(window, 'matchMedia', {
    configurable: true,
    value: vi.fn().mockReturnValue({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }),
  })
})

test('toggles between light and dark and persists the selection', async () => {
  const user = userEvent.setup()
  render(<ThemeProvider><ThemeToggle/></ThemeProvider>)

  const toggle = screen.getByRole('button', { name: '切换为黑色主题' })
  expect(toggle).toHaveAttribute('aria-pressed', 'false')
  expect(screen.getByText('白色主题')).toBeInTheDocument()

  await user.click(toggle)

  expect(screen.getByRole('button', { name: '切换为白色主题' })).toHaveAttribute('aria-pressed', 'true')
  expect(screen.getByText('黑色主题')).toBeInTheDocument()
  expect(document.documentElement).toHaveAttribute('data-theme', 'dark')
  expect(localStorage.getItem('marc-chatbot-theme')).toBe('dark')
})
