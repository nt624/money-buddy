// Thin web facade over the shared API client (see @/lib/apiClient).
// Implementations live in @pace/core/api.
import { api } from "@/lib/apiClient";

export { UserNotFoundError } from "@pace/core/api";

export const getMe = api.user.getMe;
export const updateUser = api.user.updateUser;
