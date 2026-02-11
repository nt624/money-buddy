'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useTheme } from '@/hooks/useTheme'

export function Header() {
  const pathname = usePathname()
  const { theme, toggleTheme, mounted } = useTheme()

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-14 max-w-screen-2xl items-center px-4 sm:px-6">
        {/* Logo / App Name */}
        <Link href="/" className="flex items-center space-x-2">
          <span className="text-xl font-bold text-foreground">Money Buddy</span>
        </Link>

        {/* Navigation */}
        <nav className="flex flex-1 items-center justify-end space-x-4 sm:space-x-6">
          <Link
            href="/"
            className={`text-sm font-medium transition-colors hover:text-primary ${
              pathname === '/' ? 'text-foreground' : 'text-muted-foreground'
            }`}
          >
            ホーム
          </Link>
          <Link
            href="/settings"
            className={`text-sm font-medium transition-colors hover:text-primary ${
              pathname === '/settings' ? 'text-foreground' : 'text-muted-foreground'
            }`}
          >
            設定
          </Link>

          {/* Dark Mode Toggle Button */}
          <button
            onClick={toggleTheme}
            className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-muted h-9 w-9"
            aria-label="ダークモード切り替え"
          >
            {/* SSR時のハイドレーションミスマッチを防ぐ */}
            {mounted && (
              <span className="text-lg">
                {theme === 'dark' ? '☀️' : '🌙'}
              </span>
            )}
          </button>
        </nav>
      </div>
    </header>
  )
}
