/**
 * Platform-agnostic authentication port.
 *
 * Web implements this with the Firebase Web SDK; mobile with
 * @react-native-firebase. The API client and the offline sync engine both
 * obtain their token through `getToken`, so auth has a single source of truth.
 */
export type AuthUser = {
  uid: string;
};

export type AuthPort = {
  /** Current auth token, or null when unauthenticated. */
  getToken(): Promise<string | null>;
  /** Subscribe to auth state; returns an unsubscribe function. */
  onAuthStateChanged(cb: (user: AuthUser | null) => void): () => void;
  signIn(email: string, password: string): Promise<void>;
  signUp(email: string, password: string): Promise<void>;
  signInWithGoogle(): Promise<void>;
  signOut(): Promise<void>;
};
