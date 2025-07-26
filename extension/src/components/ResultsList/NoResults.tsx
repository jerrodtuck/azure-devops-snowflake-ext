import React from 'react';

interface NoResultsProps {
  searchTerm: string;
}

export const NoResults: React.FC<NoResultsProps> = ({ searchTerm }) => {
  return (
    <div className="no-results" role="status">
      <div className="no-results-icon">🔍</div>
      <div className="no-results-text">
        No results found for "{searchTerm}"
      </div>
      <div className="no-results-hint">
        Try different search terms or check your filters
      </div>
    </div>
  );
};