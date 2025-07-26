import { useState, useEffect, useCallback } from 'react';

interface UseKeyboardNavigationProps {
  isOpen: boolean;
  itemsCount: number;
  onSelect: (index: number) => void;
  onClose: () => void;
}

/**
 * Hook for keyboard navigation in dropdown lists
 * @param props - Configuration for keyboard navigation
 * @returns Current highlighted index
 */
export function useKeyboardNavigation({
  isOpen,
  itemsCount,
  onSelect,
  onClose,
}: UseKeyboardNavigationProps) {
  const [highlightedIndex, setHighlightedIndex] = useState(-1);

  useEffect(() => {
    if (!isOpen) {
      setHighlightedIndex(-1);
    }
  }, [isOpen]);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (!isOpen || itemsCount === 0) return;

      switch (event.key) {
        case 'ArrowDown':
          event.preventDefault();
          setHighlightedIndex(prev => {
            const next = prev + 1;
            return next >= itemsCount ? 0 : next;
          });
          break;

        case 'ArrowUp':
          event.preventDefault();
          setHighlightedIndex(prev => {
            const next = prev - 1;
            return next < 0 ? itemsCount - 1 : next;
          });
          break;

        case 'Enter':
          event.preventDefault();
          if (highlightedIndex >= 0 && highlightedIndex < itemsCount) {
            onSelect(highlightedIndex);
          }
          break;

        case 'Escape':
          event.preventDefault();
          onClose();
          break;

        case 'Home':
          event.preventDefault();
          setHighlightedIndex(0);
          break;

        case 'End':
          event.preventDefault();
          setHighlightedIndex(itemsCount - 1);
          break;
      }
    },
    [isOpen, itemsCount, highlightedIndex, onSelect, onClose]
  );

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [handleKeyDown]);

  return { highlightedIndex };
}