import { CreateExpenseInput, UpdateExpenseInput, Expense, GetExpensesResponse } from "@/lib/types/expense";
import { API_BASE_URL, getAuthHeaders } from "./client";

export async function createExpense(
  input: CreateExpenseInput
): Promise<Expense> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: 'POST',
    headers,
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    const text = await res.text()
    console.error('支出の作成に失敗しました', { status: res.status, body: text })
    throw new Error('支出の作成に失敗しました')
  }

  const data = await res.json()
  if (!data || !data.expense) {
    // 開発者向けログ: 詳細なレスポンス内容を出力
    console.error('[Dev] 支出のレスポンスが正しくありません', { data })
    // ユーザー向けエラー: 固定文言のみ
    throw new Error('支出のレスポンスが正しくありません')
  }

  return data.expense
}


export async function getExpenses(): Promise<GetExpensesResponse> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: "GET",
    headers,
  });

  if (!res.ok) {
    throw new Error("支出の取得に失敗しました");
  }

  const data = await res.json();
  if (!data || !Array.isArray(data.expenses)) {
    // 開発者向けログ: 詳細なレスポンス内容を出力
    console.error('[Dev] 支出のレスポンスが正しくありません', { data })
    // ユーザー向けエラー: 固定文言のみ
    throw new Error('支出のレスポンスが正しくありません')
  }

  return data;
}

export async function updateExpense(
  id: number,
  input: UpdateExpenseInput
): Promise<Expense> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    const text = await res.text()
    console.error('支出の更新に失敗しました', { status: res.status, body: text })
    throw new Error('支出の更新に失敗しました')
  }

  const data = await res.json()
  if (!data || !data.expense) {
    // 開発者向けログ: 詳細なレスポンス内容を出力
    console.error('[Dev] 支出のレスポンスが正しくありません', { data })
    // ユーザー向けエラー: 固定文言のみ
    throw new Error('支出のレスポンスが正しくありません')
  }

  return data.expense
}

export async function deleteExpense(id: number): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'DELETE',
    headers,
  })

  if (!res.ok) {
    const text = await res.text()
    console.error('支出の削除に失敗しました', { status: res.status, body: text })
    throw new Error('支出の削除に失敗しました')
  }
}
