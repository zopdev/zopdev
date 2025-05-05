import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';
import { addDeploymentConfig } from '../../Queries/DeploymentSpace';

const useAddDeploymentConfig = () => {
  const router = useRouter();
  const params = useParams();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (id, values) => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await addDeploymentConfig(id, values);
      setError(null);
      router.push(`/applications/${params?.['application-id']}/environment`);
      return data;
    } catch (error) {
      setError(error.message);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    handleSubmit,
    isLoading,
    error,
  };
};

export default useAddDeploymentConfig;
