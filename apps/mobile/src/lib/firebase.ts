import { initializeApp, getApps, getApp } from "firebase/app";
import { initializeAuth, type Auth, type Persistence } from "firebase/auth";
import * as firebaseAuth from "firebase/auth";
import AsyncStorage from "@react-native-async-storage/async-storage";

// `getReactNativePersistence` is exported from firebase/auth's React Native
// build (resolved by Metro) but is absent from the web type definitions, so we
// bridge to it through a typed cast.
const getReactNativePersistence = (
  firebaseAuth as unknown as {
    getReactNativePersistence: (storage: unknown) => Persistence;
  }
).getReactNativePersistence;

const firebaseConfig = {
  apiKey: process.env.EXPO_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.EXPO_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.EXPO_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.EXPO_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.EXPO_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.EXPO_PUBLIC_FIREBASE_APP_ID,
};

// Lazy singleton. Auth persists across launches via AsyncStorage so the user
// stays logged in offline.
let _auth: Auth | undefined;

export function getFirebaseAuth(): Auth {
  if (!_auth) {
    const app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApp();
    _auth = initializeAuth(app, {
      persistence: getReactNativePersistence(AsyncStorage),
    });
  }
  return _auth;
}
