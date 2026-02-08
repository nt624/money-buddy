import { useEffect, useState } from "react";
import { getDashboard } from "@/lib/api/dashboard";
import { Dashboard } from "@/lib/types/dashboard";

type UseDashboardOptions = {
  enabled?: boolean; // trueの場合のみ自動fetch、デフォルトtrue
};

export function useDashboard(options: UseDashboardOptions = {}) {
  const { enabled = true } = options;
  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboard = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getDashboard();
      setDashboard(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "unknown error");
    } finally {
      setIsLoading(false);
    }
  };

  // Initial fetch on mount (enabledがtrueの場合のみ)
  useEffect(() => {
    if (enabled) {
      fetchDashboard();
    }
  }, [enabled]);

  // Refetch function for manual updates
  const refetch = () => {
    fetchDashboard();
  };

  return {
    dashboard,
    isLoading,
    error,
    refetch,
  };
}
