import { useState, useEffect } from 'react';
import { getService } from '../../Queries/Service';
import { useParams } from 'next/navigation';

const useServiceList = () => {
  const params = useParams();
  const [services, setServices] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchServices = async () => {
      try {
        const data = await getService(params?.['environment-id']);
        setServices(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchServices();
  }, []);

  return { services, loading, error };
};

export default useServiceList;
