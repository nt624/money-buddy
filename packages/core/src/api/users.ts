import { User, UpdateUserInput } from "../types/user";
import type { HttpContext } from "./http";

export class UserNotFoundError extends Error {
  constructor(message: string = "User not found") {
    super(message);
    this.name = "UserNotFoundError";
  }
}

export function createUsersApi(ctx: HttpContext) {
  return {
    async getMe(): Promise<User> {
      const headers = await ctx.authHeaders();
      const res = await fetch(`${ctx.baseUrl}/user/me`, {
        method: "GET",
        headers,
      });

      if (res.status === 404) {
        throw new UserNotFoundError();
      }

      if (!res.ok) {
        await ctx.handleApiError(res, 'ユーザー情報の取得');
      }

      const data = await res.json();
      return data;
    },

    async updateUser(input: UpdateUserInput): Promise<void> {
      const headers = await ctx.authHeaders(true);
      const res = await fetch(`${ctx.baseUrl}/user/me`, {
        method: "PUT",
        headers,
        body: JSON.stringify(input),
      });

      if (!res.ok) {
        await ctx.handleApiError(res, 'ユーザー情報の更新');
      }
    },
  };
}
