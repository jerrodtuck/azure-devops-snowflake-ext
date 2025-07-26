import { useState, useEffect, useCallback, useRef } from 'react';
import { useDebounce } from './useDebounce';
import { DropdownItem, DataType } from '../types/common.types';

interface UseSnowflakeDataProps {
  apiUrl?: string;
  dataType: DataType;
  fieldName: string;
  minSearchLength: number;
  debounceDelay: number;
}

interface UseSnowflakeDataReturn {
  data: DropdownItem[];
  loading: boolean;
  error: string | null;
  search: (term: string) => void;
}

const cache = new Map<string, { data: DropdownItem[]; timestamp: number }>();
const CACHE_DURATION = 60 * 60 * 1000;

/**
 * Hook for fetching and managing Snowflake data
 */
export function useSnowflakeData({
  apiUrl,
  dataType,
  minSearchLength,
  debounceDelay,
}: UseSnowflakeDataProps): UseSnowflakeDataReturn {
  const [data, setData] = useState<DropdownItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  
  const abortControllerRef = useRef<AbortController | null>(null);
  
  const debouncedSearchTerm = useDebounce(searchTerm, debounceDelay);

  const search = useCallback((term: string) => {
    setSearchTerm(term);
  }, []);

  useEffect(() => {
    if (debouncedSearchTerm.length < minSearchLength) {
      setData([]);
      setError(null);
      return;
    }

    const fetchData = async () => {
      const cacheKey = `${dataType}:${debouncedSearchTerm}`;
      const cached = cache.get(cacheKey);
      
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        setData(cached.data);
        setLoading(false);
        return;
      }

      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      const abortController = new AbortController();
      abortControllerRef.current = abortController;

      setLoading(true);
      setError(null);

      try {
        const baseUrl = apiUrl || 'http://localhost:8080/api';
        const url = `${baseUrl}/search/${dataType}?q=${encodeURIComponent(debouncedSearchTerm)}`;

        const response = await fetch(url, {
          signal: abortController.signal,
          headers: {
            'Accept': 'application/json',
          },
        });

        if (!response.ok) {
          throw new Error(`Search failed: ${response.statusText}`);
        }

        const result = await response.json();
        
        const items = result.data || result || [];
        const validItems = Array.isArray(items) ? items : [];
        
        cache.set(cacheKey, {
          data: validItems,
          timestamp: Date.now(),
        });

        setData(validItems);
        setError(null);
      } catch (err) {
        if (err instanceof Error && err.name === 'AbortError') {
          return;
        }
        
        setError(err instanceof Error ? err.message : 'An error occurred');
        setData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [debouncedSearchTerm, dataType, apiUrl, minSearchLength]);

  return { data, loading, error, search };
}