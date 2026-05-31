import { Category, CreateCategoryInput, UpdateCategoryInput, ReorderCategoryItem } from "../types/category";
import type { HttpContext } from "./http";

export function createCategoriesApi(ctx: HttpContext) {
  return {
    async getCategories(): Promise<Category[]> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/categories`, {
        method: "GET",
        headers,
      })

      if (!res.ok) {
        await ctx.handleApiError(res, "カテゴリの取得");
      }

      const data = await res.json()
      if (!data || !Array.isArray(data.categories)) {
        console.error('[Dev] カテゴリのレスポンスが正しくありません', { data })
        throw new Error('カテゴリのレスポンスが正しくありません')
      }

      return data.categories
    },

    async createCategory(input: CreateCategoryInput): Promise<Category> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/categories`, {
        method: "POST",
        headers: { ...headers, "Content-Type": "application/json" },
        body: JSON.stringify(input),
      })

      if (!res.ok) {
        await ctx.handleApiError(res, "カテゴリの作成");
      }

      const data = await res.json()
      return data.category
    },

    async updateCategory(id: number, input: UpdateCategoryInput): Promise<Category> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/categories/${id}`, {
        method: "PUT",
        headers: { ...headers, "Content-Type": "application/json" },
        body: JSON.stringify(input),
      })

      if (!res.ok) {
        await ctx.handleApiError(res, "カテゴリの更新");
      }

      const data = await res.json()
      return data.category
    },

    async deleteCategory(id: number): Promise<void> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/categories/${id}`, {
        method: "DELETE",
        headers,
      })

      if (!res.ok) {
        await ctx.handleApiError(res, "カテゴリの削除");
      }
    },

    async reorderCategories(items: ReorderCategoryItem[]): Promise<void> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/categories/reorder`, {
        method: "PUT",
        headers: { ...headers, "Content-Type": "application/json" },
        body: JSON.stringify({ items }),
      })

      if (!res.ok) {
        await ctx.handleApiError(res, "カテゴリの並び替え");
      }
    },
  };
}
