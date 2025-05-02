import apiClient from '../../utlis/apiClient';
import { ENVIRONMENT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = ENVIRONMENT_ENDPOINT;

export const getCronJobs = async (id) => {
  const updatedUrl = `${url}/${id}/deploymentspace/cronjob`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch cron jobs');
  }
};

export const getCronByName = async (id, name) => {
  const updatedUrl = `${url}/${id}/deploymentspace/cronjob/${name}`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch cron job details');
  }
};
