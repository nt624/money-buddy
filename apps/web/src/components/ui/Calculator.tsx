'use client'

import { useCalculator } from '@/hooks/useCalculator'

type Operator = '+' | '-' | '*' | '/'

type CalculatorProps = {
  initialValue?: number
  onApply: (value: number) => void
  onClose: () => void
}

export function Calculator({ initialValue, onApply, onClose }: CalculatorProps) {
  const { display, result, handleDigit, handleOperator, handleEquals, handleClear, handleBackspace } =
    useCalculator(initialValue)

  const base = 'flex items-center justify-center min-h-[44px] rounded-lg text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-ring'
  const num = `${base} bg-muted hover:bg-muted-foreground/20 text-foreground`
  const op = `${base} bg-primary/10 hover:bg-primary/20 text-primary`
  const eq = `${base} bg-primary hover:bg-primary-hover text-primary-foreground`
  const ac = `${base} bg-danger/10 hover:bg-danger/20 text-danger`
  const back = `${base} bg-muted hover:bg-muted-foreground/20 text-muted-foreground`

  const digit = (d: string) => () => handleDigit(d)
  const operator = (o: Operator) => () => handleOperator(o)

  return (
    <div className="w-60 bg-card border border-border rounded-xl shadow-lg p-3 space-y-2">
      <div className="px-3 py-2 bg-input rounded-lg text-right font-mono text-xl text-foreground overflow-hidden whitespace-nowrap min-h-[44px] flex items-center justify-end">
        {display}
      </div>

      <div className="grid grid-cols-4 gap-1">
        <button type="button" onClick={handleClear} className={ac}>AC</button>
        <button type="button" onClick={handleBackspace} className={back}>←</button>
        <button type="button" onClick={operator('/')} className={op}>÷</button>
        <button type="button" onClick={operator('*')} className={op}>×</button>

        <button type="button" onClick={digit('7')} className={num}>7</button>
        <button type="button" onClick={digit('8')} className={num}>8</button>
        <button type="button" onClick={digit('9')} className={num}>9</button>
        <button type="button" onClick={operator('-')} className={op}>−</button>

        <button type="button" onClick={digit('4')} className={num}>4</button>
        <button type="button" onClick={digit('5')} className={num}>5</button>
        <button type="button" onClick={digit('6')} className={num}>6</button>
        <button type="button" onClick={operator('+')} className={op}>+</button>

        <button type="button" onClick={digit('1')} className={num}>1</button>
        <button type="button" onClick={digit('2')} className={num}>2</button>
        <button type="button" onClick={digit('3')} className={num}>3</button>
        <button type="button" onClick={handleEquals} className={`${eq} row-span-2`}>=</button>

        <button type="button" onClick={digit('0')} className={`${num} col-span-2`}>0</button>
        <button type="button" onClick={digit('.')} className={num}>.</button>
      </div>

      <div className="flex gap-2 pt-1">
        <button
          type="button"
          onClick={onClose}
          className="flex-1 py-2 text-sm font-medium rounded-lg bg-secondary hover:bg-secondary-hover text-secondary-foreground transition-colors"
        >
          キャンセル
        </button>
        <button
          type="button"
          onClick={() => result !== null && onApply(result)}
          disabled={result === null}
          className="flex-1 py-2 text-sm font-medium rounded-lg bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          適用
        </button>
      </div>
    </div>
  )
}
