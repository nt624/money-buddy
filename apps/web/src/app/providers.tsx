"use client";

import { DataSourceProvider } from "@pace/core/data";
import { AuthProvider } from "@/contexts/AuthContext";
import { restDataSource } from "@/data/rest";

/**
 * Client-side provider stack: auth state + the injected DataSource that the
 * shared hooks read/write through.
 */
export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <DataSourceProvider dataSource={restDataSource}>
        {children}
      </DataSourceProvider>
    </AuthProvider>
  );
}
