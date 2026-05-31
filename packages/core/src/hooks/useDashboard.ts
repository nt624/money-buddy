import { useEffect, useState } from "react";
import { useDataSource } from "../data";
import { Dashboard } from "../types/dashboard";
import { UpdateUserInput } from "../types/user";

type SelectedMonth = { year: number; month: number };

type UseDashboardOptions = {
  enabled?: boolean; // trueの場合のみ自動fetch、デフォルトtrue
  selectedMonth?: SelectedMonth;
};

export function useDashboard(options: UseDashboardOptions = {}) {
  const ds = useDataSource();
  const { enabled = true, selectedMonth } = options;
  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchDashboard = async (month?: SelectedMonth) => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await ds.dashboard.get(month);
      setDashboard(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setIsLoading(false);
    }
  };

  // selectedMonth または enabled が変化したときに再取得
  useEffect(() => {
    if (enabled) {
      fetchDashboard(selectedMonth);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, selectedMonth?.year, selectedMonth?.month]);

  const refetch = () => {
    fetchDashboard(selectedMonth);
  };

  const updateUserSettings = async (input: UpdateUserInput): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);

    try {
      await ds.user.update(input);
      await fetchDashboard(selectedMonth);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "ユーザー設定の更新に失敗しました");
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
