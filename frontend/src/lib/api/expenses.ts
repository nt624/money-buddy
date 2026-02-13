import { CreateExpenseInput, UpdateExpenseInput, Expense, GetExpensesResponse } from "@/lib/types/expense";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function createExpense(
  input: CreateExpenseInput
): Promise<Expense> {
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: 'POST',
    headers,
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    await handleApiError(res, '支出の作成');
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
    await handleApiError(res, '支出の取得');
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
  const headers = await getAuthHeaders(true);
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    await handleApiError(res, '支出の更新');
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
    await handleApiError(res, '支出の削除');
  }
}
