"use client";

import { createContext, useContext } from "react";
import type { DataSource } from "./types";

const DataSourceContext = createContext<DataSource | null>(null);

export function DataSourceProvider({
  dataSource,
  children,
}: {
  dataSource: DataSource;
  children: React.ReactNode;
}) {
  return (
    <DataSourceContext.Provider value={dataSource}>
      {children}
    </DataSourceContext.Provider>
  );
}

export function useDataSource(): DataSource {
  const ctx = useContext(DataSourceContext);
  if (!ctx) {
    throw new Error("useDataSource must be used within a DataSourceProvider");
  }
  return ctx;
}
