import { User, UpdateUserInput } from "@/lib/types/user";

const API_BASE_URL = "http://localhost:8080";

export class UserNotFoundError extends Error {
  constructor(message: string = "User not found") {
    super(message);
    this.name = "UserNotFoundError";
  }
}

export async function getMe(): Promise<User> {
  const res = await fetch(`${API_BASE_URL}/user/me`, {
    method: "GET",
  });

  if (res.status === 404) {
    throw new UserNotFoundError();
  }

  if (!res.ok) {
    // Try to parse JSON error response from backend
    try {
      const errorData = await res.json();
      if (errorData && typeof errorData.error === 'string') {
        throw new Error(errorData.error);
      }
    } catch (parseError) {
      // If JSON parsing fails, use generic message
      console.error('ユーザー情報の取得に失敗しました', { status: res.status, parseError });
    }
    throw new Error('ユーザー情報の取得に失敗しました');
  }

  const data = await res.json();
  return data;
}

export async function updateUser(input: UpdateUserInput): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/user/me`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    let errorMessage = "ユーザー情報の更新に失敗しました";
    let errorDetail: string | undefined;

    try {
      const data: unknown = await res.json();
      if (data && typeof (data as any).error === "string") {
        errorDetail = (data as any).error;
      }
    } catch (parseError) {
      // JSON パースに失敗した場合も、ステータス情報は必ず保持する
      console.error("updateUser: failed to parse error response JSON", {
        status: res.status,
        error: parseError,
      });
    }

    if (errorDetail) {
      errorMessage += `: ${errorDetail}`;
    }

    // 開発者向けにステータスコードを必ず含める
    errorMessage += ` (status: ${res.status})`;

    throw new Error(errorMessage);
  }
}
