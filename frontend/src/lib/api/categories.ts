import { Category } from "@/lib/types/category";

export async function getCategories(): Promise<Category[]> {
  const res = await fetch("http://localhost:8080/categories")

  if (!res.ok) {
    const text = await res.text()
    console.error('[Dev] カテゴリの取得に失敗しました', { status: res.status, body: text })
    throw new Error("カテゴリの取得に失敗しました")
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