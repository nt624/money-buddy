'use client'

import { useState } from 'react'
import styles from './ExpenseForm.module.css'
import { InitialSetupRequest, FixedCostInput } from '@/lib/types/setup'

type Props = {
  onSubmit: (input: InitialSetupRequest) => Promise<void>
  isSubmitting: boolean
}

export function InitialSetupForm({ onSubmit, isSubmitting }: Props) {
  const [income, setIncome] = useState('')
  const [savingGoal, setSavingGoal] = useState('')
  const [fixedCosts, setFixedCosts] = useState<FixedCostInput[]>([
    { name: '', amount: 0 }
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
      fixedCosts: fixedCosts,
    })
  }

  const addFixedCost = () => {
    setFixedCosts([...fixedCosts, { name: '', amount: 0 }])
  }

  const removeFixedCost = (index: number) => {
    setFixedCosts(fixedCosts.filter((_, i) => i !== index))
  }

  const updateFixedCost = (index: number, field: keyof FixedCostInput, value: string | number) => {
    const updated = [...fixedCosts]
    if (field === 'name') {
      updated[index].name = value as string
    } else {
      updated[index].amount = Number(value)
    }
    setFixedCosts(updated)
  }

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <h2>初期セットアップ</h2>
      <p>収入、貯蓄目標、固定費を入力してください</p>

      <div className={styles.field}>
        <label className={styles.label}>
          月収（手取り）
          <input
            className={styles.input}
            type="number"
            value={income}
            onChange={(e) => setIncome(e.target.value)}
            placeholder="例: 300000"
          />
        </label>
        {errors.income && <p className={styles.error}>{errors.income}</p>}
      </div>

      <div className={styles.field}>
        <label className={styles.label}>
          貯蓄目標（月額）
          <input
            className={styles.input}
            type="number"
            value={savingGoal}
            onChange={(e) => setSavingGoal(e.target.value)}
            placeholder="例: 50000"
          />
        </label>
        {errors.savingGoal && <p className={styles.error}>{errors.savingGoal}</p>}
      </div>

      <div className={styles.field}>
        <label className={styles.label}>固定費</label>
        {fixedCosts.map((cost, index) => (
          <div key={index} style={{ marginBottom: '1rem', border: '1px solid #ddd', padding: '1rem', borderRadius: '4px' }}>
            <div style={{ marginBottom: '0.5rem' }}>
              <input
                className={styles.input}
                type="text"
                value={cost.name}
                onChange={(e) => updateFixedCost(index, 'name', e.target.value)}
                placeholder="固定費の名前（例: 家賃）"
              />
            </div>
            <div style={{ marginBottom: '0.5rem' }}>
              <input
                className={styles.input}
                type="number"
                value={cost.amount || ''}
                onChange={(e) => updateFixedCost(index, 'amount', e.target.value)}
                placeholder="金額"
              />
            </div>
            {errors.fixedCosts && errors.fixedCosts[index] && (
              <p className={styles.error}>{errors.fixedCosts[index]}</p>
            )}
            {fixedCosts.length > 1 && (
              <button
                type="button"
                onClick={() => removeFixedCost(index)}
                style={{ 
                  padding: '0.5rem 1rem',
                  backgroundColor: '#dc3545',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                削除
              </button>
            )}
          </div>
        ))}
        <button
          type="button"
          onClick={addFixedCost}
          style={{ 
            padding: '0.5rem 1rem',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            marginTop: '0.5rem'
          }}
        >
          固定費を追加
        </button>
      </div>

      <button className={styles.button} type="submit" disabled={isSubmitting}>
        {isSubmitting ? '送信中...' : 'セットアップを完了'}
      </button>
    </form>
  )
}
