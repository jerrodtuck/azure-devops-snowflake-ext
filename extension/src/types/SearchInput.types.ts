export interface SearchInputProps {
    value: string;
    placeholder?: string;
    minSearchLength?: number;
    onValueChange: (value: string) => void;
    onClear: () => void;
    onFocus?: () => void;
    className?: string;
    hasValue?: boolean;
}