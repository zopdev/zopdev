import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { getPods } from '../../Queries/Pods';

const usePodList = () => {
  const params = useParams();
  const [pods, setPods] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPods = async () => {
      try {
        const data = await getPods(params?.['environment-id']);
        setPods(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPods();
  }, []);

  return { pods, loading, error };
};

export default usePodList;
