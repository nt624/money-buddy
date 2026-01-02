'use client'

import { useExpenses } from '@/hooks/useExpenses'
import { ExpenseForm } from '@/components/ExpenseForm'

export default function Home() {
  const { expenses, addExpense, isSubmitting, isLoading, error } =
    useExpenses()

  return (
    <main>
      <h1>支出入力</h1>

      <ExpenseForm onSubmit={addExpense} isSubmitting={isSubmitting} />

      {isLoading && <p>読み込み中...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <h2>支出一覧</h2>
      <ul>
        {expenses.map((e) => (
          <li key={e.id}>
            {e.spent_at} / {e.amount}円
          </li>
        ))}
      </ul>
    </main>
  )
}
