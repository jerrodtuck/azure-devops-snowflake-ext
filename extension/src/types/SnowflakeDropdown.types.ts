import { DataType } from './common.types';

export interface SnowflakeDropdownProps {
    fieldName: string;
    dataType: DataType;
    onValueChange: (value: string) => void;
    apiUrl?: string;
    minSearchLength?: number;
    debounceDelay?: number;
    initialValue?: string;
}

export interface DropdownState {
    searchTerm: string;
    selectedValue: string;
    selectedLabel: string;
    isOpen: boolean;
    highlightedIndex: number;
}

// Re-export for convenience
export type { DropdownItem, DataType } from './common.types';