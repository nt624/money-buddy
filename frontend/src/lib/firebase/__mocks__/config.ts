// Mock for Firebase config used in Jest tests.
// auth.currentUser is null so getAuthHeaders returns no Authorization header.
export const auth = {
  currentUser: null,
};
