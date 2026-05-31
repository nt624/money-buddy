import { ApiClientConfig, createHttpContext } from "./http";
import { createExpensesApi } from "./expenses";
import { createCategoriesApi } from "./categories";
import { createDashboardApi } from "./dashboard";
import { createFixedCostsApi } from "./fixed-costs";
import { createMonthlySettingsApi } from "./monthly-settings";
import { createSetupApi } from "./setup";
import { createUsersApi } from "./users";

/**
 * Builds a fully-wired API client. Each platform constructs one instance with
 * its own token getter and 401 handler (see ApiClientConfig). The returned
 * object groups the REST calls by resource; every method is a plain closure
 * over the config, so individual methods can be safely destructured/re-exported.
 */
export function createApiClient(config: ApiClientConfig) {
  const ctx = createHttpContext(config);
  return {
    expenses: createExpensesApi(ctx),
    categories: createCategoriesApi(ctx),
    dashboard: createDashboardApi(ctx),
    fixedCosts: createFixedCostsApi(ctx),
    monthlySettings: createMonthlySettingsApi(ctx),
    setup: createSetupApi(ctx),
    user: createUsersApi(ctx),
  };
}

export type ApiClient = ReturnType<typeof createApiClient>;
