'use client'

import { useState, useEffect } from 'react'
import { UpdateUserInput } from '@/lib/types/user'
import { BUSINESS_MAX_AMOUNT } from '@/lib/constants'

type UserFormProps = {
  initialIncome?: number
  initialSavingGoal?: number
  onSubmit: (input: UpdateUserInput) => Promise<void>
  onCancel: () => void
  isSubmitting: boolean
  error: string | null
}

export function UserForm({
  initialIncome = 0,
  initialSavingGoal = 0,
  onSubmit,
  onCancel,
  isSubmitting,
  error
}: UserFormProps) {
  const [income, setIncome] = useState(initialIncome)
  const [savingGoal, setSavingGoal] = useState(initialSavingGoal)

  useEffect(() => {
    setIncome(initialIncome)
    setSavingGoal(initialSavingGoal)
  }, [initialIncome, initialSavingGoal])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    await onSubmit({
      income,
      saving_goal: savingGoal
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-danger/10 border border-danger rounded-lg text-danger text-sm">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="income" className="block text-sm font-medium text-foreground">
          月収（手取り）
        </label>
        <input
          type="number"
          id="income"
          value={income}
          onChange={(e) => setIncome(Number(e.target.value))}
          disabled={isSubmitting}
          min="1"
          max={BUSINESS_MAX_AMOUNT}
          required
          className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
        />
      </div>

      <div className="space-y-2">
        <label htmlFor="savingGoal" className="block text-sm font-medium text-foreground">
          貯金目標額（月）
        </label>
        <input
          type="number"
          id="savingGoal"
          value={savingGoal}
          onChange={(e) => setSavingGoal(Number(e.target.value))}
          disabled={isSubmitting}
          min="0"
          max={BUSINESS_MAX_AMOUNT}
          required
          className="w-full px-3 py-2 bg-input border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring disabled:opacity-50 disabled:cursor-not-allowed"
        />
      </div>

      <div className="flex gap-2 pt-2">
        <button
          type="submit"
          disabled={isSubmitting}
          className="px-4 py-2 text-sm font-medium rounded-lg bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isSubmitting ? '保存中...' : '保存'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={isSubmitting}
          className="px-4 py-2 text-sm font-medium rounded-lg bg-secondary hover:bg-secondary-hover text-secondary-foreground disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          キャンセル
        </button>
      </div>
    </form>
  )
}
