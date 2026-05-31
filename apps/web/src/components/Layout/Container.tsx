import { ReactNode } from 'react'

type ContainerProps = {
  children: ReactNode
  className?: string
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full'
}

const maxWidthClasses = {
  sm: 'max-w-screen-sm',
  md: 'max-w-screen-md',
  lg: 'max-w-screen-lg',
  xl: 'max-w-screen-xl',
  '2xl': 'max-w-screen-2xl',
  full: 'max-w-full',
}

export function Container({ children, className = '', maxWidth = 'xl' }: ContainerProps) {
  return (
    <div className={`container mx-auto px-4 sm:px-6 ${maxWidthClasses[maxWidth]} ${className}`}>
      {children}
    </div>
  )
}
