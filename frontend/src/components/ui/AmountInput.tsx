'use client'

import { useEffect, useRef, useState } from 'react'
import { Calculator } from './Calculator'

type AmountInputProps = {
  value: string | number
  onChange: (value: string) => void
  disabled?: boolean
  placeholder?: string
  min?: number
  max?: number
  id?: string
  className?: string
  inputClassName?: string
}

export function AmountInput({
  value,
  onChange,
  disabled,
  placeholder,
  min,
  max,
  id,
  className,
  inputClassName,
}: AmountInputProps) {
  const [isCalcOpen, setIsCalcOpen] = useState(false)
  const [prevDisabled, setPrevDisabled] = useState(disabled)
  const wrapperRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  if (prevDisabled !== disabled) {
    setPrevDisabled(disabled)
    if (disabled) setIsCalcOpen(false)
  }

  useEffect(() => {
    if (!isCalcOpen) return
    const handler = (e: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setIsCalcOpen(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [isCalcOpen])

  const handleToggle = () => {
    if (!isCalcOpen) inputRef.current?.blur()
    setIsCalcOpen(prev => !prev)
  }

  const handleApply = (val: number) => {
    onChange(String(val))
    setIsCalcOpen(false)
  }

  const initialValue = Number(value) || undefined

  return (
    <div ref={wrapperRef} className={`relative flex items-center gap-1 ${className ?? ''}`}>
      <input
        ref={inputRef}
        type="number"
        id={id}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        placeholder={placeholder}
        min={min}
        max={max}
        className={`flex-1 min-w-0 px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed ${inputClassName ?? ''}`}
      />
      <button
        type="button"
        onClick={handleToggle}
        disabled={disabled}
        aria-label="電卓を開く"
        aria-expanded={isCalcOpen}
        className="shrink-0 p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <svg width="16" height="16" viewBox="2 1 20 22" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
          <rect x="3" y="2" width="18" height="20" rx="2" strokeWidth="2"/>
          {/* + (left-top) */}
          <line x1="6" y1="8" x2="10" y2="8"/>
          <line x1="8" y1="6" x2="8" y2="10"/>
          {/* - (right-top) */}
          <line x1="14" y1="8" x2="18" y2="8"/>
          {/* × (left-bottom) */}
          <line x1="6" y1="14" x2="10" y2="18"/>
          <line x1="10" y1="14" x2="6" y2="18"/>
          {/* = (right-bottom) */}
          <line x1="14" y1="14.5" x2="18" y2="14.5"/>
          <line x1="14" y1="17.5" x2="18" y2="17.5"/>
        </svg>
      </button>
      {isCalcOpen && (
        <div className="absolute top-full right-0 mt-1 z-[60]">
          <Calculator
            initialValue={initialValue}
            onApply={handleApply}
            onClose={() => setIsCalcOpen(false)}
          />
        </div>
      )}
    </div>
  )
}
