import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { getCronByName } from '../../Queries/CronJob';

const useCronJobDetails = (formData) => {
  const params = useParams();
  const [cronJob, setCronJob] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCronJobDetails = async () => {
      try {
        const data = await getCronByName(params?.['environment-id'], formData?.name);
        setCronJob(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchCronJobDetails();
  }, []);

  return { cronJob, loading, error };
};

export default useCronJobDetails;
