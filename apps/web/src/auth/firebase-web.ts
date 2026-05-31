import {
  onAuthStateChanged,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut as firebaseSignOut,
  GoogleAuthProvider,
  signInWithPopup,
} from "firebase/auth";
import type { AuthPort } from "@pace/core/auth";
import { getFirebaseAuth } from "@/lib/firebase/config";

/**
 * Web implementation of the core AuthPort, backed by the Firebase Web SDK.
 * This is the single source of truth for auth tokens used by both the API
 * client and (later) the sync engine.
 */
export const firebaseWebAuthPort: AuthPort = {
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
    const provider = new GoogleAuthProvider();
    await signInWithPopup(getFirebaseAuth(), provider);
  },

  async signOut() {
    await firebaseSignOut(getFirebaseAuth());
  },
};
