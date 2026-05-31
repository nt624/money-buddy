import { useState } from "react";
import {
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  Text,
  TextInput,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useAuth } from "@pace/core/auth";

export default function LoginScreen() {
  const { signIn, signUp } = useAuth();
  const [isSignUp, setIsSignUp] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async () => {
    setError("");
    setLoading(true);
    try {
      if (isSignUp) {
        await signUp(email, password);
      } else {
        await signIn(email, password);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView className="flex-1 bg-white">
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : undefined}
        className="flex-1 justify-center px-6"
      >
        <Text className="text-2xl font-bold text-gray-900 mb-1">Pace Wallet</Text>
        <Text className="text-sm text-gray-500 mb-8">
          {isSignUp ? "アカウントを作成" : "ログイン"}
        </Text>

        <TextInput
          className="border border-gray-300 rounded-lg px-4 py-3 mb-3 text-gray-900"
          placeholder="メールアドレス"
          placeholderTextColor="#9ca3af"
          autoCapitalize="none"
          keyboardType="email-address"
          value={email}
          onChangeText={setEmail}
        />
        <TextInput
          className="border border-gray-300 rounded-lg px-4 py-3 mb-3 text-gray-900"
          placeholder="パスワード"
          placeholderTextColor="#9ca3af"
          secureTextEntry
          value={password}
          onChangeText={setPassword}
        />

        {error ? <Text className="text-red-600 text-sm mb-3">{error}</Text> : null}

        <Pressable
          className="bg-emerald-600 rounded-lg py-3 items-center active:opacity-80"
          disabled={loading}
          onPress={onSubmit}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold">
              {isSignUp ? "登録する" : "ログイン"}
            </Text>
          )}
        </Pressable>

        <Pressable className="mt-4 items-center" onPress={() => setIsSignUp((v) => !v)}>
          <Text className="text-emerald-700 text-sm">
            {isSignUp ? "すでにアカウントをお持ちの方" : "アカウントを作成する"}
          </Text>
        </Pressable>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}
