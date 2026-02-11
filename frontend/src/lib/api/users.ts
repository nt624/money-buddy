import { User, UpdateUserInput } from "@/lib/types/user";
import { API_BASE_URL, getAuthHeaders, handleApiError } from "./client";

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
    await handleApiError(res, 'ユーザー情報の取得');
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
    await handleApiError(res, 'ユーザー情報の更新');
  }
}
