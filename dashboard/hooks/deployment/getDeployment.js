import { useState, useEffect } from 'react';
import { getService } from '../../Queries/Service';
import { useParams } from 'next/navigation';
import { getDeployment } from '../../Queries/Deployment';

const useDeploymentList = () => {
  const params = useParams();
  const [deployment, setDeployment] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDeployment = async () => {
      try {
        const data = await getDeployment(params?.['environment-id']);
        setDeployment(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchDeployment();
  }, []);

  return { deployment, loading, error };
};

export default useDeploymentList;
