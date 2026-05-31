'use client'

import { useState } from 'react'
import { UpdateUserInput } from '@pace/core/types/user'
import { BUSINESS_MAX_AMOUNT } from '@/lib/constants'
import { AmountInput } from '@/components/ui/AmountInput'

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
  const [prevInitialIncome, setPrevInitialIncome] = useState(initialIncome)
  const [prevInitialSavingGoal, setPrevInitialSavingGoal] = useState(initialSavingGoal)
  const [income, setIncome] = useState(initialIncome)
  const [savingGoal, setSavingGoal] = useState(initialSavingGoal)

  // Sync state when initial values change (React-recommended pattern)
  if (prevInitialIncome !== initialIncome || prevInitialSavingGoal !== initialSavingGoal) {
    setPrevInitialIncome(initialIncome)
    setPrevInitialSavingGoal(initialSavingGoal)
    setIncome(initialIncome)
    setSavingGoal(initialSavingGoal)
  }

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
        <AmountInput
          id="income"
          value={income}
          onChange={(v) => setIncome(Number(v))}
          disabled={isSubmitting}
          min={1}
          max={BUSINESS_MAX_AMOUNT}
          className="w-full"
        />
      </div>

      <div className="space-y-2">
        <label htmlFor="savingGoal" className="block text-sm font-medium text-foreground">
          貯金目標額（月）
        </label>
        <AmountInput
          id="savingGoal"
          value={savingGoal}
          onChange={(v) => setSavingGoal(Number(v))}
          disabled={isSubmitting}
          min={0}
          max={BUSINESS_MAX_AMOUNT}
          className="w-full"
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
