// Thin web facade over the shared API client (see @/lib/apiClient).
// Implementations live in @pace/core/api.
import { api } from "@/lib/apiClient";

export const createFixedCost = api.fixedCosts.createFixedCost;
export const getFixedCosts = api.fixedCosts.getFixedCosts;
export const updateFixedCost = api.fixedCosts.updateFixedCost;
export const deleteFixedCost = api.fixedCosts.deleteFixedCost;
