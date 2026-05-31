import { createRestDataSource } from "@pace/core/data";
import { api } from "@/lib/apiClient";

/**
 * Phase 1: mobile reads/writes go directly to the backend over REST (online).
 * Phase 2 swaps this for a local-first DataSource backed by packages/sync.
 */
export const restDataSource = createRestDataSource(api);
