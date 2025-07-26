import React from 'react';

interface ErrorStateProps {
  message: string;
}

export const ErrorState: React.FC<ErrorStateProps> = ({ message }) => {
  return (
    <div className="state-message error" role="alert">
      <div className="error-icon">⚠️</div>
      <div className="error-text">{message}</div>
      <div className="error-hint">Please try again or contact support</div>
    </div>
  );
};