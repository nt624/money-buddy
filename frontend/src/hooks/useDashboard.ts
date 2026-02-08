import { useEffect, useState } from "react";
import { getDashboard } from "@/lib/api/dashboard";
import { Dashboard } from "@/lib/types/dashboard";

export function useDashboard() {
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

  // Initial fetch on mount
  useEffect(() => {
    fetchDashboard();
  }, []);

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
