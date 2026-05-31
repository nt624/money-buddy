import { createApiClient } from "@pace/core/api";
import { firebaseRnAuthPort } from "@/auth/firebase-rn";

const API_BASE_URL = process.env.EXPO_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

/**
 * Mobile singleton API client. Token retrieval is delegated to the RN AuthPort;
 * on 401 we sign out, which makes the auth gate redirect to the login screen.
 */
export const api = createApiClient({
  baseUrl: API_BASE_URL,
  getToken: () => firebaseRnAuthPort.getToken(),
  onUnauthorized: () => {
    firebaseRnAuthPort.signOut().catch((error) => {
      console.error("サインアウトに失敗しました:", error);
    });
  },
});
