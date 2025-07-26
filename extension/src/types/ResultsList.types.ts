import { DropdownItem } from './common.types';

export interface ResultsListProps {
  items: DropdownItem[];
  loading: boolean;
  error: string | null;
  highlightedIndex: number;
  selectedValue: string;
  searchTerm: string;
  onSelect: (value: string, label: string) => void;
  minSearchLength: number;
  className?: string;
}

export type ResultsState = 
  | 'loading'
  | 'error'
  | 'no-results'
  | 'too-short'
  | 'results';

export interface ResultItemProps {
  item: DropdownItem;
  isHighlighted: boolean;
  isSelected: boolean;
  onClick: () => void;
}