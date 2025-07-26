import React from 'react';
import { TypeSelectorProps } from '../../types/TypeSelector.types';
import './TypeSelector.styles.css';

export const TypeSelector: React.FC<TypeSelectorProps> = ({
  currentType,
  onChange,
  disabled = false,
  className = '',
  dataTypes = [],
  loading = false,
}) => {
  // Handle radio button change
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange(event.target.value);
  };

  if (loading) {
    return <div className="type-selector-loading">Loading data types...</div>;
  }

  if (!Array.isArray(dataTypes) || dataTypes.length === 0) {
    return null;
  }

  // Don't show selector if only one type is available
  if (dataTypes.length === 1) {
    return null;
  }

  return (
    <div className={`type-selector ${className}`} role="radiogroup" aria-label="Data type selector">
      {Array.isArray(dataTypes) && dataTypes.map((type) => (
        <label
          key={type.id}
          className={`type-selector-label ${disabled ? 'disabled' : ''}`}
          title={type.description}
        >
          <input
            type="radio"
            name="dataType"
            value={type.id}
            checked={currentType === type.id}
            onChange={handleChange}
            disabled={disabled}
            className="type-selector-input"
            aria-describedby={`${type.id}-description`}
          />
          <span className="type-selector-text">
            {type.icon && <span className="type-selector-icon">{type.icon}</span>}
            {type.name}
          </span>
          {/* Hidden description for screen readers */}
          <span id={`${type.id}-description`} className="sr-only">
            {type.description}
          </span>
        </label>
      ))}
    </div>
  );
};