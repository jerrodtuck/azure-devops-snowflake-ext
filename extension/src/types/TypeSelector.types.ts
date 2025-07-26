export interface DataTypeInfo {
  id: string;
  name: string;
  description: string;
  icon?: string;
}

export interface TypeSelectorProps {
  currentType: string;
  onChange: (type: string) => void;
  disabled?: boolean;
  className?: string;
  dataTypes?: DataTypeInfo[];
  loading?: boolean;
}