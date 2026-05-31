import { Category, CreateCategoryInput, UpdateCategoryInput, ReorderCategoryItem } from "@/lib/types/category";
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
    console.error('[Dev] カテゴリのレスポンスが正しくありません', { data })
    throw new Error('カテゴリのレスポンスが正しくありません')
  }

  return data.categories
}

export async function createCategory(input: CreateCategoryInput): Promise<Category> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/categories`, {
    method: "POST",
    headers: { ...headers, "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    await handleApiError(res, "カテゴリの作成");
  }

  const data = await res.json()
  return data.category
}

export async function updateCategory(id: number, input: UpdateCategoryInput): Promise<Category> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/categories/${id}`, {
    method: "PUT",
    headers: { ...headers, "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })

  if (!res.ok) {
    await handleApiError(res, "カテゴリの更新");
  }

  const data = await res.json()
  return data.category
}

export async function deleteCategory(id: number): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/categories/${id}`, {
    method: "DELETE",
    headers,
  })

  if (!res.ok) {
    await handleApiError(res, "カテゴリの削除");
  }
}

export async function reorderCategories(items: ReorderCategoryItem[]): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/categories/reorder`, {
    method: "PUT",
    headers: { ...headers, "Content-Type": "application/json" },
    body: JSON.stringify({ items }),
  })

  if (!res.ok) {
    await handleApiError(res, "カテゴリの並び替え");
  }
}
