import { Expense } from "@/lib/types/expense"

type Props = {
  expenses: Expense[]
  onEdit: (expense: Expense) => void
  onDelete: (id: number) => void
}

export function ExpenseList({ expenses, onEdit, onDelete }: Props) {
  return (
    <ul style={{ listStyle: 'none', padding: 0 }}>
      {expenses.map((e) => (
        <li key={e.id} style={{ 
          marginBottom: '1rem', 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '4px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          gap: '1rem'
        }}>
          <div style={{ flex: 1 }}>
            <div style={{ marginBottom: '0.5rem' }}>
              <span style={{
                display: 'inline-block',
                padding: '0.25rem 0.5rem',
                borderRadius: '4px',
                fontSize: '0.875rem',
                fontWeight: 'bold',
                backgroundColor: e.status === 'confirmed' ? '#28a745' : '#ffc107',
                color: e.status === 'confirmed' ? 'white' : '#000',
              }}>
                {e.status === 'confirmed' ? '確定' : '予定'}
              </span>
            </div>
            <div>
              <span>{e.spent_at}</span> / 
              <span>¥{e.amount.toLocaleString()}</span> / 
              <span>カテゴリ: {e.category.name}</span>
              {e.memo && <span> / {e.memo}</span>}
            </div>
          </div>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button
              onClick={() => onEdit(e)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '0.875rem'
              }}
            >
              編集
            </button>
            <button
              onClick={() => onDelete(e.id)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#dc3545',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '0.875rem'
              }}
            >
              削除
            </button>
          </div>
        </li>
      ))}
    </ul>
  )
}
