import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { getPodByName } from '../../Queries/Pods';

const usePodDetails = (formData) => {
  const params = useParams();
  const [pod, setPod] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPodDetails = async () => {
      try {
        const data = await getPodByName(params?.['environment-id'], formData?.name);
        setPod(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPodDetails();
  }, []);

  return { pod, loading, error };
};

export default usePodDetails;
