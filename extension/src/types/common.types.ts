/**
 * Common types used across multiple components
 */

// Data type for dropdown selection - now dynamic
export type DataType = string;

// Generic dropdown item structure
export interface DropdownItem {
  value: string;
  label: string;
}

// API response structure
export interface ApiResponse<T> {
  data: T;
  metadata?: {
    exported_at: string;
    row_count: number;
    source: string;
    cached: boolean;
  };
}

// Error state
export interface ErrorState {
  message: string;
  code?: string;
  details?: unknown;
}