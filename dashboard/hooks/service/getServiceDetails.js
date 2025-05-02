import { useState, useEffect } from 'react';
import { getServiceByName } from '../../Queries/Service';
import { useParams } from 'next/navigation';

const useServiceDetails = (formData) => {
  const params = useParams();
  const [service, setService] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchServiceDetails = async () => {
      try {
        const data = await getServiceByName(params?.['environment-id'], formData?.name);
        setService(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchServiceDetails();
  }, []);

  return { service, loading, error };
};

export default useServiceDetails;
