import { ActivityIndicator, Pressable, ScrollView, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useAuth } from "@pace/core/auth";
import { useDashboard, useExpenses, useUser } from "@pace/core/hooks";

function yen(value: number): string {
  return `¥${value.toLocaleString("ja-JP")}`;
}

export default function HomeScreen() {
  const { signOut } = useAuth();
  const { needsSetup, isLoading: userLoading } = useUser();
  const { expenses, selectedMonth, navigateMonth, isLoading: expensesLoading } = useExpenses();
  const { dashboard } = useDashboard({ selectedMonth });

  if (userLoading) {
    return (
      <SafeAreaView className="flex-1 items-center justify-center bg-gray-50">
        <ActivityIndicator />
      </SafeAreaView>
    );
  }

  if (needsSetup) {
    return (
      <SafeAreaView className="flex-1 items-center justify-center bg-gray-50 px-6">
        <Text className="text-base text-gray-700 text-center">
          初期設定が必要です。Web アプリで設定を完了してください。
        </Text>
        <Pressable className="mt-6" onPress={() => signOut()}>
          <Text className="text-emerald-700">ログアウト</Text>
        </Pressable>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView className="flex-1 bg-gray-50">
      <ScrollView contentContainerClassName="p-4">
        <View className="flex-row items-center justify-between mb-4">
          <Text className="text-lg font-bold text-gray-900">Pace Wallet</Text>
          <Pressable onPress={() => signOut()}>
            <Text className="text-sm text-gray-500">ログアウト</Text>
          </Pressable>
        </View>

        {/* 月ナビゲーション */}
        <View className="flex-row items-center justify-between mb-4">
          <Pressable onPress={() => navigateMonth("prev")} className="px-3 py-1">
            <Text className="text-emerald-700 text-lg">‹</Text>
          </Pressable>
          <Text className="text-base font-semibold text-gray-800">
            {selectedMonth.year}年{selectedMonth.month}月
          </Text>
          <Pressable onPress={() => navigateMonth("next")} className="px-3 py-1">
            <Text className="text-emerald-700 text-lg">›</Text>
          </Pressable>
        </View>

        {/* 残額カード */}
        <View className="bg-white rounded-2xl p-5 mb-4 shadow-sm">
          <Text className="text-sm text-gray-500 mb-1">今月の残額</Text>
          <Text className="text-3xl font-bold text-emerald-600">
            {dashboard ? yen(dashboard.remaining) : "—"}
          </Text>
          {dashboard ? (
            <View className="flex-row justify-between mt-4">
              <View>
                <Text className="text-xs text-gray-400">使える額</Text>
                <Text className="text-sm font-medium text-gray-800">
                  {yen(dashboard.variable_budget)}
                </Text>
              </View>
              <View>
                <Text className="text-xs text-gray-400">確定支出</Text>
                <Text className="text-sm font-medium text-gray-800">
                  {yen(dashboard.confirmed_expenses)}
                </Text>
              </View>
              <View>
                <Text className="text-xs text-gray-400">予定支出</Text>
                <Text className="text-sm font-medium text-gray-800">
                  {yen(dashboard.planned_expenses)}
                </Text>
              </View>
            </View>
          ) : null}
        </View>

        {/* 支出リスト */}
        <Text className="text-sm font-semibold text-gray-700 mb-2">支出</Text>
        {expensesLoading ? (
          <ActivityIndicator className="mt-4" />
        ) : expenses.length === 0 ? (
          <Text className="text-sm text-gray-400 mt-2">この月の支出はまだありません</Text>
        ) : (
          expenses.map((expense) => (
            <View
              key={expense.id}
              className="bg-white rounded-xl px-4 py-3 mb-2 flex-row items-center justify-between"
            >
              <View>
                <Text className="text-sm font-medium text-gray-900">
                  {expense.category?.name ?? "未分類"}
                </Text>
                <Text className="text-xs text-gray-400">
                  {expense.spent_at}
                  {expense.status === "planned" ? "・予定" : ""}
                </Text>
              </View>
              <Text className="text-base font-semibold text-gray-900">
                {yen(expense.amount)}
              </Text>
            </View>
          ))
        )}
      </ScrollView>
    </SafeAreaView>
  );
}
