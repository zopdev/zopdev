import { useContext, useEffect } from 'react';
import { getApplication } from '../../Queries/Application';

const { AppContext } = require('../../libs/context');
const { getCloudAccounts } = require('../../Queries/CloudAccount');

export function useInitializeHeader() {
  const { setAppData } = useContext(AppContext);

  const handleAppData = (entity, values) =>
    setAppData((prevValues) => ({ ...prevValues, [entity]: values }));

  const fetchCloudAccounts = async () => {
    try {
      const data = await getCloudAccounts();
      handleAppData('CLOUD_ACCOUNT_DATA', { data: data.data, isSuccess: true });
    } catch (error) {
      handleAppData('CLOUD_ACCOUNT_DATA', { data: [], isSuccess: false });
    }
  };

  const fetchApplications = async () => {
    try {
      const data = await getApplication();
      handleAppData('APPLICATION_DATA', { data: data.data, isSuccess: true });
    } catch (error) {
      handleAppData('APPLICATION_DATA', { data: [], isSuccess: false });
    }
  };

  useEffect(() => {
    fetchCloudAccounts();
    fetchApplications();
  }, []);

  return {};
}
