// Mock for Firebase config used in Jest tests.
// getFirebaseAuth().currentUser is null so getAuthHeaders returns no Authorization header.
export function getFirebaseAuth() {
  return { currentUser: null };
}

export function getFirebaseApp() {
  return {};
}
