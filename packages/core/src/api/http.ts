/**
 * Platform-agnostic HTTP context for the API layer.
 *
 * No Firebase or DOM coupling lives here: the caller injects how to obtain an
 * auth token (`getToken`) and what to do on 401 (`onUnauthorized`). Web wires
 * these to Firebase Web SDK + `window.location`; mobile wires them to
 * @react-native-firebase + navigation.
 */
export type ApiClientConfig = {
  baseUrl: string;
  /** Returns the current auth token, or null when unauthenticated. */
  getToken: () => Promise<string | null>;
  /** Called once when the server responds 401 (token invalid/expired). */
  onUnauthorized?: () => void;
};

export type HttpContext = {
  baseUrl: string;
  /** Builds request headers, optionally including Content-Type for a JSON body. */
  authHeaders: (hasBody?: boolean) => Promise<Record<string, string>>;
  /** Throws a user-facing Error; triggers onUnauthorized on 401. */
  handleApiError: (response: Response, operation: string) => Promise<never>;
};

export function createHttpContext(config: ApiClientConfig): HttpContext {
  async function authHeaders(hasBody: boolean = false): Promise<Record<string, string>> {
    const headers: Record<string, string> = {};

    // リクエストボディがある場合のみContent-Typeを付与
    if (hasBody) {
      headers["Content-Type"] = "application/json";
    }

    // トークンが取得できる場合のみAuthorizationヘッダーを付与
    try {
      const token = await config.getToken();
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }
    } catch (error) {
      console.error("Failed to get auth token:", error);
      // トークン取得失敗時はヘッダーなしで続行（バックエンドが401を返す）
    }

    return headers;
  }

  async function handleApiError(response: Response, operation: string): Promise<never> {
    // 401 Unauthorized - トークンが無効または期限切れ
    if (response.status === 401) {
      console.warn("認証エラーが発生しました。ログアウトします。");
      config.onUnauthorized?.();
      throw new Error("認証の有効期限が切れました。再度ログインしてください。");
    }

    // その他のエラー
    let errorMessage = `${operation}に失敗しました`;
    try {
      const errorData = await response.json();
      if (errorData.error) {
        errorMessage = errorData.error;
      }
    } catch {
      // JSONパースエラーは無視
    }

    throw new Error(errorMessage);
  }

  return {
    baseUrl: config.baseUrl,
    authHeaders,
    handleApiError,
  };
}
