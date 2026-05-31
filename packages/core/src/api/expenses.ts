import { CreateExpenseInput, UpdateExpenseInput, Expense, GetExpensesResponse } from "../types/expense";
import type { HttpContext } from "./http";

export function createExpensesApi(ctx: HttpContext) {
  return {
    async createExpense(input: CreateExpenseInput): Promise<Expense> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/expenses`, {
        method: 'POST',
        headers,
        body: JSON.stringify(input),
      })

      if (!res.ok) {
        await ctx.handleApiError(res, '支出の作成');
      }

      const data = await res.json()
      if (!data || !data.expense) {
        // 開発者向けログ: 詳細なレスポンス内容を出力
        console.error('[Dev] 支出のレスポンスが正しくありません', { data })
        // ユーザー向けエラー: 固定文言のみ
        throw new Error('支出のレスポンスが正しくありません')
      }

      return data.expense
    },

    async getExpenses(params?: { year: number; month: number }): Promise<GetExpensesResponse> {
      const headers = await ctx.authHeaders();
      const url = new URL(`${ctx.baseUrl}/expenses`);
      if (params) {
        url.searchParams.set('year', String(params.year));
        url.searchParams.set('month', String(params.month));
      }
      const res = await fetch(url.toString(), {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '支出の取得');
      }

      const data = await res.json();
      if (!data || !Array.isArray(data.expenses)) {
        // 開発者向けログ: 詳細なレスポンス内容を出力
        console.error('[Dev] 支出のレスポンスが正しくありません', { data })
        // ユーザー向けエラー: 固定文言のみ
        throw new Error('支出のレスポンスが正しくありません')
      }

      return data;
    },

    async updateExpense(id: number, input: UpdateExpenseInput): Promise<Expense> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/expenses/${id}`, {
        method: 'PUT',
        headers,
        body: JSON.stringify(input),
      })

      if (!res.ok) {
        await ctx.handleApiError(res, '支出の更新');
      }

      const data = await res.json()
      if (!data || !data.expense) {
        // 開発者向けログ: 詳細なレスポンス内容を出力
        console.error('[Dev] 支出のレスポンスが正しくありません', { data })
        // ユーザー向けエラー: 固定文言のみ
        throw new Error('支出のレスポンスが正しくありません')
      }

      return data.expense
    },

    async deleteExpense(id: number): Promise<void> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/expenses/${id}`, {
        method: 'DELETE',
        headers,
      })

      if (!res.ok) {
        await ctx.handleApiError(res, '支出の削除');
      }
    },
  };
}
