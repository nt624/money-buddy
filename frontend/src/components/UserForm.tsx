'use client'

import { useState, useEffect } from 'react'
import { UpdateUserInput } from '@/lib/types/user'

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
    <form onSubmit={handleSubmit} style={{ marginBottom: '2rem' }}>
      {error && (
        <div style={{ 
          padding: '1rem', 
          marginBottom: '1rem', 
          backgroundColor: '#fee2e2', 
          border: '1px solid #ef4444',
          borderRadius: '0.375rem',
          color: '#991b1b'
        }}>
          {error}
        </div>
      )}

      <div style={{ marginBottom: '1rem' }}>
        <label htmlFor="income" style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 500 }}>
          月収（手取り）
        </label>
        <input
          type="number"
          id="income"
          value={income}
          onChange={(e) => setIncome(Number(e.target.value))}
          disabled={isSubmitting}
          min="0"
          max="1000000000"
          required
          style={{
            width: '100%',
            padding: '0.5rem',
            border: '1px solid #d1d5db',
            borderRadius: '0.375rem',
            fontSize: '1rem'
          }}
        />
      </div>

      <div style={{ marginBottom: '1rem' }}>
        <label htmlFor="savingGoal" style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 500 }}>
          貯金目標額（月）
        </label>
        <input
          type="number"
          id="savingGoal"
          value={savingGoal}
          onChange={(e) => setSavingGoal(Number(e.target.value))}
          disabled={isSubmitting}
          min="0"
          max="1000000000"
          required
          style={{
            width: '100%',
            padding: '0.5rem',
            border: '1px solid #d1d5db',
            borderRadius: '0.375rem',
            fontSize: '1rem'
          }}
        />
      </div>

      <div style={{ display: 'flex', gap: '1rem' }}>
        <button
          type="submit"
          disabled={isSubmitting}
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: isSubmitting ? '#9ca3af' : '#3b82f6',
            color: 'white',
            border: 'none',
            borderRadius: '0.375rem',
            cursor: isSubmitting ? 'not-allowed' : 'pointer',
            fontWeight: 500
          }}
        >
          {isSubmitting ? '保存中...' : '保存'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={isSubmitting}
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: 'transparent',
            color: '#6b7280',
            border: '1px solid #d1d5db',
            borderRadius: '0.375rem',
            cursor: isSubmitting ? 'not-allowed' : 'pointer'
          }}
        >
          キャンセル
        </button>
      </div>
    </form>
  )
}
