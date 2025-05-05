import { useParams, useRouter } from 'next/navigation';
import { useState } from 'react';
import { addEnvironment } from '../../Queries/Application';

const useAddEnvironment = () => {
  const router = useRouter();
  const params = useParams();
  const [values, setValues] = useState({ name: '' });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (id, values) => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await addEnvironment(id, values);
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
    values,
    setValues,
    handleSubmit,
    isLoading,
    error,
  };
};

export default useAddEnvironment;
