import { useCallback, useEffect, useState } from "react";
import {
  getMonthlySettings,
  upsertMonthlySettings,
  deleteMonthlySettings,
} from "@/lib/api/monthly-settings";
import { MonthlySettings, UpsertMonthlySettingsInput } from "@/lib/types/monthly-settings";

type SelectedMonth = { year: number; month: number };
type FallbackUser = { income: number; saving_goal: number };

export function useMonthlySettings(
  selectedMonth: SelectedMonth,
  fallbackUser?: FallbackUser | null
) {
  const [settings, setSettings] = useState<MonthlySettings | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const fetchSettings = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await getMonthlySettings(selectedMonth.year, selectedMonth.month);
      setSettings(data);
    } catch {
      // API 失敗時はユーザーのグローバルデフォルト値に静かにフォールバック
      if (fallbackUser) {
        setSettings({
          year: selectedMonth.year,
          month: selectedMonth.month,
          income: fallbackUser.income,
          saving_goal: fallbackUser.saving_goal,
          is_custom: false,
        });
      }
      // fallbackUser もない場合は settings を null のままにする（エラーは表示しない）
    } finally {
      setIsLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedMonth.year, selectedMonth.month, fallbackUser?.income, fallbackUser?.saving_goal]);

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
      const updated = await upsertMonthlySettings(input);
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
      await deleteMonthlySettings(selectedMonth.year, selectedMonth.month);
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
