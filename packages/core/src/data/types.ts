import type { CreateExpenseInput, UpdateExpenseInput, Expense } from "../types/expense";
import type {
  Category,
  CreateCategoryInput,
  UpdateCategoryInput,
  ReorderCategoryItem,
} from "../types/category";
import type { Dashboard } from "../types/dashboard";
import type { User, UpdateUserInput } from "../types/user";
import type {
  FixedCost,
  CreateFixedCostInput,
  UpdateFixedCostInput,
} from "../types/fixed-cost";
import type {
  MonthlySettings,
  UpsertMonthlySettingsInput,
} from "../types/monthly-settings";

export type Period = { year: number; month: number };

/**
 * Platform-agnostic data-access seam consumed by the shared hooks.
 *
 * Web injects a REST-backed implementation (apps/web/src/data/rest.ts).
 * Mobile will inject a local-first implementation backed by the offline store
 * (packages/sync). Methods are promise-based today; reactive `observe*`
 * variants are added in Phase 2 when the mobile local store lands.
 */
export interface ExpensesDataSource {
  list(period: Period): Promise<Expense[]>;
  create(input: CreateExpenseInput): Promise<Expense>;
  update(id: number, input: UpdateExpenseInput): Promise<Expense>;
  remove(id: number): Promise<void>;
}

export interface CategoriesDataSource {
  list(): Promise<Category[]>;
  create(input: CreateCategoryInput): Promise<Category>;
  update(id: number, input: UpdateCategoryInput): Promise<Category>;
  remove(id: number): Promise<void>;
  reorder(items: ReorderCategoryItem[]): Promise<void>;
}

export interface DashboardDataSource {
  get(period?: Period): Promise<Dashboard>;
}

export interface UserDataSource {
  /** Resolves the current user; throws UserNotFoundError when setup is needed. */
  getMe(): Promise<User>;
  update(input: UpdateUserInput): Promise<void>;
}

export interface FixedCostsDataSource {
  list(): Promise<FixedCost[]>;
  create(input: CreateFixedCostInput): Promise<FixedCost>;
  update(id: number, input: UpdateFixedCostInput): Promise<FixedCost>;
  remove(id: number): Promise<void>;
}

export interface MonthlySettingsDataSource {
  get(year: number, month: number): Promise<MonthlySettings>;
  upsert(input: UpsertMonthlySettingsInput): Promise<MonthlySettings>;
  remove(year: number, month: number): Promise<void>;
}

export interface DataSource {
  expenses: ExpensesDataSource;
  categories: CategoriesDataSource;
  dashboard: DashboardDataSource;
  user: UserDataSource;
  fixedCosts: FixedCostsDataSource;
  monthlySettings: MonthlySettingsDataSource;
}
