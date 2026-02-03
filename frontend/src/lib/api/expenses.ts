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
    throw new Error(`Failed to create expense: ${res.status} ${text}`)
  }

  const data = await res.json()
  if (!data || !data.expense) {
    throw new Error(`Invalid createExpense response: ${JSON.stringify(data)}`)
  }

  return data.expense
}


export async function getExpenses(): Promise<GetExpensesResponse> {
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: "GET",
  });

  if (!res.ok) {
    throw new Error("Failed to fetch expenses");
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.expenses)) {
    throw new Error(`Invalid getExpenses response: ${JSON.stringify(data)}`)
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
    throw new Error(`Failed to update expense: ${res.status} ${text}`)
  }

  const data = await res.json()
  if (!data || !data.expense) {
    throw new Error(`Invalid updateExpense response: ${JSON.stringify(data)}`)
  }

  return data.expense
}

export async function deleteExpense(id: number): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'DELETE',
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(`Failed to delete expense: ${res.status} ${text}`)
  }
}
