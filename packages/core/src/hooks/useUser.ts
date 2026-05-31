import { useEffect, useState } from "react";
import { useDataSource } from "../data";
import { UserNotFoundError } from "../api";
import { User } from "../types/user";

export function useUser() {
  const ds = useDataSource();
  const [user, setUser] = useState<User | null>(null);
  const [needsSetup, setNeedsSetup] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUser = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const userData = await ds.user.getMe();
      setUser(userData);
      setNeedsSetup(false);
    } catch (err) {
      if (err instanceof UserNotFoundError) {
        setNeedsSetup(true);
        setUser(null);
      } else {
        setError(err instanceof Error ? err.message : "エラーが発生しました");
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchUser();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ds]);

  return {
    user,
    needsSetup,
    isLoading,
    error,
    refetchUser: fetchUser,
  };
}
