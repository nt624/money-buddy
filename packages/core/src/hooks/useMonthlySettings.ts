import { useCallback, useEffect, useState } from "react";
import { useDataSource } from "../data";
import { MonthlySettings, UpsertMonthlySettingsInput } from "../types/monthly-settings";

type SelectedMonth = { year: number; month: number };
type FallbackUser = { income: number; saving_goal: number };

export function useMonthlySettings(
  selectedMonth: SelectedMonth,
  fallbackUser?: FallbackUser | null
) {
  const ds = useDataSource();
  const [settings, setSettings] = useState<MonthlySettings | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const fetchSettings = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await ds.monthlySettings.get(selectedMonth.year, selectedMonth.month);
      setSettings(data);
    } catch (err) {
      // API 失敗時はユーザーのグローバルデフォルト値にフォールバックしつつエラーを記録
      setError(err instanceof Error ? err.message : "月設定の読み込みに失敗しました");
      if (fallbackUser) {
        setSettings({
          year: selectedMonth.year,
          month: selectedMonth.month,
          income: fallbackUser.income,
          saving_goal: fallbackUser.saving_goal,
          is_custom: false,
        });
      }
    } finally {
      setIsLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ds, selectedMonth.year, selectedMonth.month, fallbackUser?.income, fallbackUser?.saving_goal]);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const saveSettings = async (
    input: UpsertMonthlySettingsInput,
    onSuccess?: () => void
  ): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      const updated = await ds.monthlySettings.upsert(input);
      setSettings(updated);
      onSuccess?.();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "月設定の保存に失敗しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetSettings = async (onSuccess?: () => void): Promise<boolean> => {
    setIsSubmitting(true);
    setError(null);
    try {
      await ds.monthlySettings.remove(selectedMonth.year, selectedMonth.month);
      await fetchSettings();
      onSuccess?.();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "月設定の削除に失敗しました");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const openModal = () => setIsModalOpen(true);
  const closeModal = () => setIsModalOpen(false);

  return {
    settings,
    isLoading,
    isSubmitting,
    error,
    isModalOpen,
    openModal,
    closeModal,
    saveSettings,
    resetSettings,
    refetch: fetchSettings,
  };
}
