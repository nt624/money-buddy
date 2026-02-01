import { useEffect, useState } from "react";
import { getMe, UserNotFoundError } from "@/lib/api/users";
import { User } from "@/lib/types/user";

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [needsSetup, setNeedsSetup] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUser = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const userData = await getMe();
      setUser(userData);
      setNeedsSetup(false);
    } catch (err) {
      if (err instanceof UserNotFoundError) {
        setNeedsSetup(true);
        setUser(null);
      } else {
        setError(err instanceof Error ? err.message : "unknown error");
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  return {
    user,
    needsSetup,
    isLoading,
    error,
    refetchUser: fetchUser,
  };
}
