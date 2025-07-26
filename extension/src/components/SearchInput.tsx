import React from 'react';
import { SearchInputProps } from '../types/SearchInput.types';
import './SearchInput.styles.css';

export const SearchInput = React.forwardRef<HTMLInputElement, SearchInputProps>(({
    value,
    placeholder = 'Type to search...',
    onValueChange,
    onClear,
    onFocus,
    className = '',
    hasValue = false,
}, ref) => {
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onValueChange(e.target.value);
    };

    const handleClear = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
        onClear();
    };

    return (
        <div className={`search-input-wrapper ${className}`}>
            <input
                ref={ref}
                type="text"
                className={`search-input ${hasValue ? 'has-value' : ''}`}
                value={value}
                onChange={handleInputChange}
                onFocus={onFocus}
                placeholder={placeholder}
                autoComplete="off"
            />

            {/* Clear button - only visible when there's a value */}
            {hasValue && (
                <button
                    className="clear-button"
                    onClick={handleClear}
                    type="button"
                    title="Clear"
                    aria-label="Clear search"
                >
                    <svg width="16" height="16" viewBox="0 0 16 16">
                        <path
                            d="M8 8.707l3.646 3.647.708-.708L8.707 8l3.647-3.646-.708-.708L8 7.293 4.354 3.646l-.708.708L7.293 8l-3.647 3.646.708.708L8 8.707z"
                            fill="currentColor"
                        />
                    </svg>
                </button>
            )}

            {/* Dropdown icon - only visible when no value */}
            {!hasValue && (
                <svg className="dropdown-icon" viewBox="0 0 16 16">
                    <path
                        d="M4 6l4 4 4-4"
                        stroke="currentColor"
                        strokeWidth="2"
                        fill="none"
                    />
                </svg>
            )}
        </div>
    );
});

SearchInput.displayName = 'SearchInput';