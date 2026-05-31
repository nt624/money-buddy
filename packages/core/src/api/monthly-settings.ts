import { MonthlySettings, UpsertMonthlySettingsInput } from "../types/monthly-settings";
import type { HttpContext } from "./http";

export function createMonthlySettingsApi(ctx: HttpContext) {
  return {
    async getMonthlySettings(year: number, month: number): Promise<MonthlySettings> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/monthly-settings?year=${year}&month=${month}`, {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '月設定の取得');
      }

      return res.json();
    },

    async upsertMonthlySettings(input: UpsertMonthlySettingsInput): Promise<MonthlySettings> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/monthly-settings`, {
        method: "PUT",
        headers,
        body: JSON.stringify(input),
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '月設定の保存');
      }

      return res.json();
    },

    async deleteMonthlySettings(year: number, month: number): Promise<void> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/monthly-settings?year=${year}&month=${month}`, {
        method: "DELETE",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '月設定の削除');
      }
    },
  };
}
