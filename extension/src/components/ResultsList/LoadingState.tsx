import React from 'react';

export const LoadingState: React.FC = () => {
  return (
    <div className="loading-container" role="status" aria-live="polite">
      <div className="loading-spinner" aria-hidden="true"></div>
      <div className="loading-text">Searching...</div>
    </div>
  );
};