import { useState, useEffect } from 'react';
import { getCloudAccounts } from '../../Queries/CloudAccount';

const useCloudAccountList = () => {
  const [cloudAccounts, setCloudAccounts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCloudAccounts = async () => {
      try {
        const data = await getCloudAccounts();
        setCloudAccounts(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchCloudAccounts();
  }, []);

  return { cloudAccounts, loading, error };
};

export default useCloudAccountList;
