import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { addApplication } from '../../Queries/Application';

// import { addApplication } from '../../Queries/Application';

const useAddApplication = () => {
  const router = useRouter();
  const [values, setValues] = useState({
    name: '',
    environments: [],
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (values) => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await addApplication(values);
      setError(null);
      router.push('/applications');
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

export default useAddApplication;
