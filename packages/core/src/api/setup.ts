import { InitialSetupRequest, InitialSetupResponse } from "../types/setup";
import type { HttpContext } from "./http";

export function createSetupApi(ctx: HttpContext) {
  return {
    async submitInitialSetup(input: InitialSetupRequest): Promise<InitialSetupResponse> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/setup`, {
        method: "POST",
        headers,
        body: JSON.stringify(input),
      });

      if (!res.ok) {
        await ctx.handleApiError(res, '初期設定の送信');
      }

      const data = await res.json();
      return data;
    },
  };
}
