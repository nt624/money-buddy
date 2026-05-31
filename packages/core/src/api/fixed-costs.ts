import {
  FixedCost,
  CreateFixedCostInput,
  UpdateFixedCostInput,
  CreateFixedCostResponse,
  GetFixedCostsResponse,
  UpdateFixedCostResponse,
} from "../types/fixed-cost";
import type { HttpContext } from "./http";

export function createFixedCostsApi(ctx: HttpContext) {
  return {
    async createFixedCost(input: CreateFixedCostInput): Promise<FixedCost> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/fixed-costs`, {
        method: "POST",
        headers,
        body: JSON.stringify(input),
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '固定費の作成');
      }

      const data: CreateFixedCostResponse = await res.json();
      if (!data || !data.fixed_cost) {
        console.error('固定費のレスポンスが正しくありません', { data });
        throw new Error('固定費のレスポンスが正しくありません');
      }

      return data.fixed_cost;
    },

    async getFixedCosts(): Promise<GetFixedCostsResponse> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/fixed-costs`, {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '固定費の取得');
      }

      const data = await res.json();
      if (!data || !Array.isArray(data.fixed_costs)) {
        console.error('固定費のレスポンスが正しくありません', { data });
        throw new Error('固定費のレスポンスが正しくありません');
      }

      return data;
    },

    async updateFixedCost(id: number, input: UpdateFixedCostInput): Promise<FixedCost> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/fixed-costs/${id}`, {
        method: "PUT",
        headers,
        body: JSON.stringify(input),
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '固定費の更新');
      }

      const data: UpdateFixedCostResponse = await res.json();
      if (!data || !data.fixed_cost) {
        console.error('固定費のレスポンスが正しくありません', { data });
        throw new Error('固定費のレスポンスが正しくありません');
      }

      return data.fixed_cost;
    },

    async deleteFixedCost(id: number): Promise<void> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/fixed-costs/${id}`, {
        method: "DELETE",
        headers,
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '固定費の削除');
      }
    },
  };
}
