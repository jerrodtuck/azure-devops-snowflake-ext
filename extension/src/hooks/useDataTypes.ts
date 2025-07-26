import { useState, useEffect } from 'react';

interface DataTypeInfo {
  id: string;
  name: string;
  description: string;
  icon: string;
}

interface UseDataTypesReturn {
  dataTypes: DataTypeInfo[];
  loading: boolean;
  error: string | null;
  defaultType: string;
}

export function useDataTypes(apiUrl?: string): UseDataTypesReturn {
  const [dataTypes, setDataTypes] = useState<DataTypeInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [defaultType, setDefaultType] = useState<string>('cc');

  useEffect(() => {
    const fetchDataTypes = async () => {
      try {
        const baseUrl = apiUrl || 'http://localhost:8080/api';
        const response = await fetch(`${baseUrl}/config`);
        
        if (!response.ok) {
          throw new Error('Failed to fetch configuration');
        }

        const data = await response.json();
        const dataTypes = data.dataTypes || [];
        setDataTypes(Array.isArray(dataTypes) ? dataTypes : []);
        setDefaultType(data.defaultType || 'cc');
        setError(null);
      } catch (err) {
        console.error('Error fetching data types:', err);
        // Fallback to default types
        setDataTypes([
          { id: 'cc', name: 'Cost Centers', description: 'Company cost centers', icon: '💰' },
          { id: 'wbs', name: 'WBS Elements', description: 'Work Breakdown Structure', icon: '📊' }
        ]);
        setError(err instanceof Error ? err.message : 'Failed to load data types');
      } finally {
        setLoading(false);
      }
    };

    fetchDataTypes();
  }, [apiUrl]);

  return { dataTypes, loading, error, defaultType };
}