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
    const text = await res.text();
    throw new Error(`ユーザー情報の取得に失敗しました: ${res.status} ${text}`);
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
    const error = await res.json().catch(() => ({ error: "不明なエラー" }));
    throw new Error(error.error || `ユーザー情報の更新に失敗しました: ${res.status}`);
  }
}
