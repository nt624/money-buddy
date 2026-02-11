'use client'

import { useState } from 'react'
import { InitialSetupRequest, FixedCostInput } from '@/lib/types/setup'

type Props = {
  onSubmit: (input: InitialSetupRequest) => Promise<void>
  isSubmitting: boolean
}

type FixedCostWithId = FixedCostInput & { id: number }

export function InitialSetupForm({ onSubmit, isSubmitting }: Props) {
  const [income, setIncome] = useState('')
  const [savingGoal, setSavingGoal] = useState('')
  const [nextId, setNextId] = useState(1)
  const [fixedCosts, setFixedCosts] = useState<FixedCostWithId[]>([
    { id: 0, name: '', amount: 0 }
  ])

  const [errors, setErrors] = useState<{
    income?: string
    savingGoal?: string
    fixedCosts?: string[]
  }>({})

  const validate = (): boolean => {
    const newErrors: typeof errors = {}

    const incomeNumber = Number(income)
    if (!income || isNaN(incomeNumber) || incomeNumber < 1) {
      newErrors.income = '収入は1以上の数値で入力してください'
    }

    const savingGoalNumber = Number(savingGoal)
    if (!savingGoal || isNaN(savingGoalNumber) || savingGoalNumber < 0) {
      newErrors.savingGoal = '貯蓄目標は0以上の数値で入力してください'
    }

    const fixedCostErrors: string[] = []
    fixedCosts.forEach((cost, index) => {
      if (!cost.name.trim()) {
        fixedCostErrors[index] = '固定費の名前を入力してください'
      } else if (cost.amount < 1) {
        fixedCostErrors[index] = '金額は1以上で入力してください'
      }
    })

    if (fixedCostErrors.length > 0) {
      newErrors.fixedCosts = fixedCostErrors
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    await onSubmit({
      income: Number(income),
      savingGoal: Number(savingGoal),
      fixedCosts: fixedCosts.map(({ id, ...cost }) => cost),
    })
  }

  const addFixedCost = () => {
    setFixedCosts([...fixedCosts, { id: nextId, name: '', amount: 0 }])
    setNextId(nextId + 1)
    setErrors({ ...errors, fixedCosts: undefined })
  }

  const removeFixedCost = (id: number) => {
    setFixedCosts(fixedCosts.filter((cost) => cost.id !== id))
    setErrors({ ...errors, fixedCosts: undefined })
  }

  const updateFixedCost = (id: number, field: keyof FixedCostInput, value: string | number) => {
    const updated = fixedCosts.map(cost => {
      if (cost.id !== id) return cost
      return {
        ...cost,
        [field]: field === 'amount' ? Number(value) : value
      }
    })
    setFixedCosts(updated)
    setErrors({ ...errors, fixedCosts: undefined })
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-2xl mx-auto p-6 bg-card rounded-lg shadow-lg space-y-6">
      <div className="text-center space-y-2">
        <h2 className="text-2xl font-bold text-foreground">初期セットアップ</h2>
        <p className="text-muted-foreground">収入、貯蓄目標、固定費を入力してください</p>
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium text-foreground">
          月収（手取り）
          <input
            className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
            type="number"
            value={income}
            onChange={(e) => setIncome(e.target.value)}
            placeholder="例: 300000"
          />
        </label>
        {errors.income && <p className="text-sm text-danger">{errors.income}</p>}
      </div>

      <div className="space-y-2">
        <label className="block text-sm font-medium text-foreground">
          貯蓄目標（月額）
          <input
            className="mt-1 block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
            type="number"
            value={savingGoal}
            onChange={(e) => setSavingGoal(e.target.value)}
            placeholder="例: 50000"
          />
        </label>
        {errors.savingGoal && <p className="text-sm text-danger">{errors.savingGoal}</p>}
      </div>

      <div className="space-y-4">
        <label className="block text-sm font-medium text-foreground">固定費</label>
        <div className="space-y-3">
          {fixedCosts.map((cost, index) => (
            <div key={cost.id} className="p-4 border border-border rounded-lg bg-muted space-y-3">
              <div>
                <input
                  className="block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                  type="text"
                  value={cost.name}
                  onChange={(e) => updateFixedCost(cost.id, 'name', e.target.value)}
                  placeholder="固定費の名前（例: 家賃）"
                />
              </div>
              <div>
                <input
                  className="block w-full px-3 py-2 bg-input border border-border rounded-md text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
                  type="number"
                  value={cost.amount || ''}
                  onChange={(e) => updateFixedCost(cost.id, 'amount', e.target.value)}
                  placeholder="金額"
                />
              </div>
              {errors.fixedCosts && errors.fixedCosts[index] && (
                <p className="text-sm text-danger">{errors.fixedCosts[index]}</p>
              )}
              {fixedCosts.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeFixedCost(cost.id)}
                  className="px-4 py-2 text-sm font-medium rounded bg-danger hover:bg-danger-hover text-danger-foreground transition-colors"
                >
                  削除
                </button>
              )}
            </div>
          ))}
        </div>
        <button
          type="button"
          onClick={addFixedCost}
          className="px-4 py-2 text-sm font-medium rounded bg-success hover:bg-success-hover text-success-foreground transition-colors"
        >
          固定費を追加
        </button>
      </div>

      <button 
        className="w-full px-4 py-3 text-base font-medium rounded bg-primary hover:bg-primary-hover text-primary-foreground disabled:opacity-60 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-ring" 
        type="submit" 
        disabled={isSubmitting}
      >
        {isSubmitting ? '送信中...' : 'セットアップを完了'}
      </button>
    </form>
  )
}
