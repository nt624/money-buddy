import { createRestDataSource } from "@pace/core/data";
import { api } from "@/lib/apiClient";

/**
 * Web's REST DataSource: the shared factory wired to the web API client.
 */
export const restDataSource = createRestDataSource(api);
