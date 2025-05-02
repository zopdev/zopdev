import { useEffect, useState } from 'react';
import { getApplicationById } from '../../Queries/Application';
import { useParams } from 'next/navigation';

const useGetDeploymentSpace = () => {
  const params = useParams();

  const [value, setValue] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const getAppById = async () => {
      try {
        const data = await getApplicationById(params?.['application-id']);
        setValue(data);
        return value;
      } catch (e) {
        setError(e);
      } finally {
        setLoading(false);
      }
    };

    getAppById();
  }, []);

  return { value, loading, error };
};

export default useGetDeploymentSpace;
