import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { getDeploymentByName } from '../../Queries/Deployment';

const useDeploymentDetails = (formData) => {
  const params = useParams();
  const [deployment, setDeployment] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDeploymentDetails = async () => {
      try {
        const data = await getDeploymentByName(params?.['environment-id'], formData?.name);
        setDeployment(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchDeploymentDetails();
  }, []);

  return { deployment, loading, error };
};

export default useDeploymentDetails;
