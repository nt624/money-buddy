import { Dashboard } from "../types/dashboard";
import type { HttpContext } from "./http";

export function createDashboardApi(ctx: HttpContext) {
  return {
    async getDashboard(params?: { year: number; month: number }): Promise<Dashboard> {
      const headers = await ctx.authHeaders();
      const query = params ? `?year=${params.year}&month=${params.month}` : '';
      const res = await fetch(`${ctx.baseUrl}/dashboard${query}`, {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, 'ダッシュボードの取得');
      }

      const data = await res.json();
      return data;
    },
  };
}
