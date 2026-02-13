import { Category } from "@/lib/types/category";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

export async function getCategories(): Promise<Category[]> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/categories`, {
    method: "GET",
    headers,
  })

  if (!res.ok) {
    await handleApiError(res, "カテゴリの取得");
  }

  const data = await res.json()
  if (!data || !Array.isArray(data.categories)) {
    // 開発者向けログ: 詳細なレスポンス内容を出力
    console.error('[Dev] カテゴリのレスポンスが正しくありません', { data })
    // ユーザー向けエラー: 固定文言のみ
    throw new Error('カテゴリのレスポンスが正しくありません')
  }

  return data.categories
}