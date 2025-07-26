import React from 'react';
import { ResultItemProps } from '../../types/ResultsList.types';

export const ResultItem: React.FC<ResultItemProps> = ({
  item,
  isHighlighted,
  isSelected,
  onClick,
}) => {
  const classNames = [
    'result-item',
    isHighlighted && 'highlighted',
    isSelected && 'selected',
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div
      className={classNames}
      onClick={onClick}
      role="option"
      aria-selected={isSelected}
      aria-label={`${item.value} - ${item.label}`}
    >
      <div className="result-value">{item.value}</div>
      {item.value !== item.label && (
        <div className="result-label">{item.label}</div>
      )}
    </div>
  );
};