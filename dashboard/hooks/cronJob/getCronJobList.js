import { useState, useEffect } from 'react';
import { getCronJobs } from '../../Queries/CronJob';
import { useParams } from 'next/navigation';

const useCronJobList = () => {
  const params = useParams();
  const [cronJob, setCronJob] = useState();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCronJob = async () => {
      try {
        const data = await getCronJobs(params?.['environment-id']);
        setCronJob(data);
      } catch (error) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchCronJob();
  }, []);

  return { cronJob, loading, error };
};

export default useCronJobList;
