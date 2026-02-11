import { User, UpdateUserInput } from "@/lib/types/user";
import { API_BASE_URL, getAuthHeaders } from "./client";

export class UserNotFoundError extends Error {
  constructor(message: string = "User not found") {
    super(message);
    this.name = "UserNotFoundError";
  }
}

export async function getMe(): Promise<User> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/user/me`, {
    method: "GET",
    headers,
  });

  if (res.status === 404) {
    throw new UserNotFoundError();
  }

  if (!res.ok) {
    // 開発者向けログ: 詳細なエラー情報を出力
    try {
      const errorData = await res.json();
      console.error('[Dev] ユーザー情報の取得に失敗しました', { 
        status: res.status, 
        error: errorData 
      });
    } catch (parseError) {
      console.error('[Dev] ユーザー情報の取得に失敗しました（JSON parse error）', { 
        status: res.status, 
        parseError 
      });
    }
    // ユーザー向けエラー: 固定文言のみ
    throw new Error('ユーザー情報の取得に失敗しました');
  }

  const data = await res.json();
  return data;
}

export async function updateUser(input: UpdateUserInput): Promise<void> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_BASE_URL}/user/me`, {
    method: "PUT",
    headers,
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    // 開発者向けログ: 詳細なエラー情報を出力
    try {
      const data: unknown = await res.json();
      console.error('[Dev] ユーザー情報の更新に失敗しました', {
        status: res.status,
        error: data
      });
    } catch (parseError) {
      console.error('[Dev] ユーザー情報の更新に失敗しました（JSON parse error）', {
        status: res.status,
        error: parseError,
      });
    }

    // ユーザー向けエラー: 固定文言のみ
    throw new Error('ユーザー情報の更新に失敗しました');
  }
}
