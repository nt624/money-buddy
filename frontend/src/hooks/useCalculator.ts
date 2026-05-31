import { useState } from 'react'

type Operator = '+' | '-' | '*' | '/'

type CalcState = {
  display: string
  operand: number | null
  operator: Operator | null
  waitingForOperand: boolean
}

const INITIAL_STATE: CalcState = {
  display: '0',
  operand: null,
  operator: null,
  waitingForOperand: false,
}

function applyOperator(a: number, op: Operator, b: number): number | 'error' {
  switch (op) {
    case '+': return a + b
    case '-': return a - b
    case '*': return a * b
    case '/': return b === 0 ? 'error' : a / b
  }
}

export function useCalculator(initialValue?: number) {
  const [state, setState] = useState<CalcState>(() => ({
    ...INITIAL_STATE,
    display: initialValue != null && !isNaN(initialValue) ? String(initialValue) : '0',
  }))

  const handleDigit = (digit: string) => {
    setState(prev => {
      if (prev.display === 'エラー') {
        return { ...INITIAL_STATE, display: digit === '.' ? '0.' : digit }
      }
      if (prev.waitingForOperand) {
        return { ...prev, display: digit === '.' ? '0.' : digit, waitingForOperand: false }
      }
      if (digit === '.' && prev.display.includes('.')) return prev
      if (prev.display === '0' && digit !== '.') return { ...prev, display: digit }
      return { ...prev, display: prev.display + digit }
    })
  }

  const handleOperator = (op: Operator) => {
    setState(prev => {
      if (prev.display === 'エラー') return prev
      const current = parseFloat(prev.display)
      if (prev.operand !== null && prev.operator && !prev.waitingForOperand) {
        const result = applyOperator(prev.operand, prev.operator, current)
        if (result === 'error') return { ...INITIAL_STATE, display: 'エラー' }
        return { display: String(result), operand: result, operator: op, waitingForOperand: true }
      }
      return { ...prev, operand: current, operator: op, waitingForOperand: true }
    })
  }

  const handleEquals = () => {
    setState(prev => {
      if (prev.operand === null || prev.operator === null) return prev
      const current = parseFloat(prev.display)
      const result = applyOperator(prev.operand, prev.operator, current)
      if (result === 'error') return { ...INITIAL_STATE, display: 'エラー' }
      const rounded = Math.round(result as number)
      return { display: String(rounded), operand: null, operator: null, waitingForOperand: false }
    })
  }

  const handleClear = () => setState(INITIAL_STATE)

  const handleBackspace = () => {
    setState(prev => {
      if (prev.waitingForOperand) return prev
      const next = prev.display.length > 1 ? prev.display.slice(0, -1) : '0'
      return { ...prev, display: next }
    })
  }

  const parsed = parseFloat(state.display)
  const result = state.display === 'エラー' || isNaN(parsed) ? null : parsed

  return { display: state.display, result, handleDigit, handleOperator, handleEquals, handleClear, handleBackspace }
}
