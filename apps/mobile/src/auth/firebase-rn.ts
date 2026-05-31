import {
  onAuthStateChanged,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut as firebaseSignOut,
} from "firebase/auth";
import type { AuthPort } from "@pace/core/auth";
import { getFirebaseAuth } from "@/lib/firebase";

/**
 * React Native implementation of the core AuthPort, backed by the Firebase JS
 * SDK with AsyncStorage persistence. Shared API client and (later) sync engine
 * obtain their token through here.
 */
export const firebaseRnAuthPort: AuthPort = {
  async getToken() {
    const user = getFirebaseAuth().currentUser;
    if (!user) return null;
    return user.getIdToken();
  },

  onAuthStateChanged(cb) {
    return onAuthStateChanged(getFirebaseAuth(), (user) => {
      cb(user ? { uid: user.uid } : null);
    });
  },

  async signIn(email, password) {
    await signInWithEmailAndPassword(getFirebaseAuth(), email, password);
  },

  async signUp(email, password) {
    await createUserWithEmailAndPassword(getFirebaseAuth(), email, password);
  },

  async signInWithGoogle() {
    // Google サインインは expo-auth-session / Google SDK を要するため後続フェーズで対応する。
    throw new Error("Google ログインはモバイルでは未対応です");
  },

  async signOut() {
    await firebaseSignOut(getFirebaseAuth());
  },
};
