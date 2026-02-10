import { useEffect, useState } from "react";
import { getDashboard } from "@/lib/api/dashboard";
import { updateUser } from "@/lib/api/users";
import { Dashboard } from "@/lib/types/dashboard";
import { UpdateUserInput } from "@/lib/types/user";

type UseDashboardOptions = {
  enabled?: boolean; // trueの場合のみ自動fetch、デフォルトtrue
};

export function useDashboard(options: UseDashboardOptions = {}) {
  const { enabled = true } = options;
  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

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

  // Update user settings
  const updateUserSettings = async (input: UpdateUserInput): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);

    try {
      await updateUser(input);
      await fetchDashboard(); // ダッシュボードを再取得
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update user settings");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    dashboard,
    isLoading,
    error,
    refetch,
    updateUserSettings,
    isSubmitting,
  };
}
