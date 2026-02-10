import { CreateExpenseInput, UpdateExpenseInput, Expense, GetExpensesResponse } from "@/lib/types/expense";

const API_BASE_URL = "http://localhost:8080";

export async function createExpense(
  input: CreateExpenseInput
): Promise<Expense> {
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(`支出の作成に失敗しました: ${res.status} ${text}`)
  }

  const data = await res.json()
  if (!data || !data.expense) {
    throw new Error(`支出のレスポンスが正しくありません: ${JSON.stringify(data)}`)
  }

  return data.expense
}


export async function getExpenses(): Promise<GetExpensesResponse> {
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: "GET",
  });

  if (!res.ok) {
    throw new Error("支出の取得に失敗しました");
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.expenses)) {
    throw new Error(`支出のレスポンスが正しくありません: ${JSON.stringify(data)}`)
  }

  return data;
}

export async function updateExpense(
  id: number,
  input: UpdateExpenseInput
): Promise<Expense> {
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(`支出の更新に失敗しました: ${res.status} ${text}`)
  }

  const data = await res.json()
  if (!data || !data.expense) {
    throw new Error(`支出のレスポンスが正しくありません: ${JSON.stringify(data)}`)
  }

  return data.expense
}

export async function deleteExpense(id: number): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'DELETE',
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(`支出の削除に失敗しました: ${res.status} ${text}`)
  }
}
