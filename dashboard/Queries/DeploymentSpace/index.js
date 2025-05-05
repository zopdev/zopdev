import apiClient from '../../utlis/apiClient';
import { ENVIRONMENT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = ENVIRONMENT_ENDPOINT;

export const addDeploymentConfig = async (id, values) => {
  const updatedUrl = `${url}/${id}/deploymentspace`;
  try {
    const response = await apiClient.post(updatedUrl, values);
    return response;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to add deployment space');
  }
};
