import React, { useState, useEffect, useRef } from 'react';
import { SearchInput } from '../SearchInput';
import { ResultsList } from '../ResultsList';
import { TypeSelector } from '../TypeSelector';
import { useSnowflakeData } from '../../hooks/useSnowflakeData';
import { useClickOutside } from '../../hooks/useClickOutside';
import { useKeyboardNavigation } from '../../hooks/useKeyboardNavigation';
import { useDataTypes } from '../../hooks/useDataTypes';
import { SnowflakeDropdownProps, DropdownState, DataType } from '../../types/SnowflakeDropdown.types';
import './SnowflakeDropdown.styles.css';

export const SnowflakeDropdown: React.FC<SnowflakeDropdownProps>= ({
    fieldName,
    dataType: initialDataType,
    onValueChange,
    apiUrl,
    minSearchLength = 3,
    debounceDelay = 300,
    initialValue = '',
}) => {
    const [state, setState] = useState<DropdownState>({
        searchTerm: initialValue,
        selectedValue: initialValue,
        selectedLabel: initialValue,
        isOpen: false,
        highlightedIndex: -1,
    });

    const [currentType, setCurrentType] = useState<DataType>(initialDataType);

    const containerRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    // Fetch available data types
    const { dataTypes, loading: typesLoading, defaultType } = useDataTypes(apiUrl);

    // Set default type if initial type is not provided
    useEffect(() => {
        if (!initialDataType && defaultType && currentType !== defaultType) {
            setCurrentType(defaultType);
        }
    }, [defaultType, initialDataType, currentType]);

    const { data, loading, error, search } = useSnowflakeData({
        apiUrl,
        fieldName,
        dataType: currentType,
        minSearchLength,
        debounceDelay,
    });

    useClickOutside(containerRef, () => {
        setState(prev => ({ ...prev, isOpen: false }));
    });

    const { highlightedIndex } = useKeyboardNavigation({
        isOpen: state.isOpen,
        itemsCount: Array.isArray(data) ? data.length : 0,
        onSelect: (index: number) => {
          if (Array.isArray(data) && data[index]) {
            selectItem(data[index].value, data[index].label);
          }
        },
        onClose: () => setState(prev => ({ ...prev, isOpen: false })),
      });

    useEffect(() => {
        setState(prev => ({ 
            ...prev, 
            highlightedIndex,
        }));
        }, [highlightedIndex]);

    const handleSearchChange = (value: string) => {
        setState(prev => ({ 
            ...prev, 
            searchTerm: value, 
            isOpen: value.length >= minSearchLength,
        }));

        if (value !== state.selectedLabel) {
            setState(prev => ({ 
                ...prev, selectedValue: '', 
                selectedLabel: '',
            }));
            onValueChange('');
        }

        if (value.length >= minSearchLength) {
            search(value);
        }
    };

    const selectItem = (value: string, label: string) => {
        setState(prev => ({
            ...prev,
            selectedValue: value,
            selectedLabel: label,
            searchTerm: '',
            isOpen: false,
        }));
        onValueChange(value);
    };

    const handleClear = () => {
        setState({
            searchTerm: '',
            selectedValue: '',
            selectedLabel: '',
            isOpen: false,
            highlightedIndex: -1,
        });
        onValueChange('');
        inputRef.current?.focus();
    };

    const handleTypeChange = (type: DataType) => {
        setCurrentType(type);
        if (state.searchTerm.length >= minSearchLength) {
            search(state.searchTerm);
        }
    };

    const handleFocus = () => {
        if (state.searchTerm.length >= minSearchLength && Array.isArray(data) && data.length > 0) {
            setState(prev => ({ ...prev, isOpen: true }));
        }
    };

    return (
        <div ref={containerRef} className="snowflake-dropdown">
          {/* Type selector for CC vs WBS */}
          <TypeSelector
            currentType={currentType}
            onChange={handleTypeChange}
            dataTypes={dataTypes}
            loading={typesLoading}
          />
  
          {/* Search input */}
          <div className="dropdown-container">
            <SearchInput
              ref={inputRef}
              value={state.searchTerm}
              placeholder={`Search ${currentType === 'cc' ? 'cost centers' : 'WBS elements'}...`}
              onValueChange={handleSearchChange}
              onClear={handleClear}
              onFocus={handleFocus}
              hasValue={state.searchTerm.length > 0}
              minSearchLength={minSearchLength}
            />
  
            {/* Results dropdown */}
            {state.isOpen && (
              <ResultsList
                items={data}
                loading={loading}
                error={error}
                highlightedIndex={state.highlightedIndex}
                selectedValue={state.selectedValue}
                searchTerm={state.searchTerm}
                onSelect={selectItem}
                minSearchLength={minSearchLength}
              />
            )}
          </div>
        </div>
    );
};