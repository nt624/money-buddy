import { AuthProvider } from "@pace/core/auth";
import { DataSourceProvider } from "@pace/core/data";
import { firebaseRnAuthPort } from "@/auth/firebase-rn";
import { restDataSource } from "@/data/rest";

/**
 * Provider stack: shared auth state + the injected (REST) DataSource that the
 * shared hooks read through.
 */
export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider authPort={firebaseRnAuthPort}>
      <DataSourceProvider dataSource={restDataSource}>{children}</DataSourceProvider>
    </AuthProvider>
  );
}
