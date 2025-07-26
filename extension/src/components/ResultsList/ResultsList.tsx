import React, { useRef, useEffect } from 'react';
import { ResultsListProps, ResultsState } from '../../types/ResultsList.types';
import { ResultItem } from './ResultItem';
import { LoadingState } from './LoadingState';
import { ErrorState } from './ErrorState';
import { NoResults } from './NoResults';
import './ResultsList.styles.css';

export const ResultsList: React.FC<ResultsListProps> = ({
  items,
  loading,
  error,
  highlightedIndex,
  selectedValue,
  searchTerm,
  onSelect,
  minSearchLength,
  className = '',
}) => {
  const containerRef = useRef<HTMLDivElement>(null);

  const getState = (): ResultsState => {
    if (loading) return 'loading';
    if (error) return 'error';
    if (searchTerm.length < minSearchLength) return 'too-short';
    if (items.length === 0) return 'no-results';
    return 'results';
  };

  const currentState = getState();

  useEffect(() => {
    if (currentState === 'results' && highlightedIndex >= 0) {
      const container = containerRef.current;
      if (container) {
        const highlightedElement = container.querySelector(
          `.result-item:nth-child(${highlightedIndex + 1})`
        );
        if (highlightedElement) {
          highlightedElement.scrollIntoView({
            block: 'nearest',
            behavior: 'smooth',
          });
        }
      }
    }
  }, [highlightedIndex, currentState]);

  return (
    <div className={`search-results open ${className}`}>
      <div className="results-container" ref={containerRef}>
        {currentState === 'loading' && <LoadingState />}

        {currentState === 'error' && <ErrorState message={error!} />}

        {currentState === 'too-short' && (
          <div className="state-message">
            Type at least {minSearchLength} characters to search...
          </div>
        )}

        {currentState === 'no-results' && (
          <NoResults searchTerm={searchTerm} />
        )}

        {currentState === 'results' && (
          <>
            {items.map((item, index) => (
              <ResultItem
                key={item.value}
                item={item}
                isHighlighted={index === highlightedIndex}
                isSelected={item.value === selectedValue}
                onClick={() => onSelect(item.value, item.label)}
              />
            ))}
          </>
        )}
      </div>

      {currentState === 'results' && items.length > 0 && (
        <div className="keyboard-hint">
          <span>↑↓ Navigate</span>
          <span>Enter Select</span>
          <span>Esc Cancel</span>
        </div>
      )}
    </div>
  );
};