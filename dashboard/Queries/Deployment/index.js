import apiClient from '../../utlis/apiClient';
import { ENVIRONMENT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = ENVIRONMENT_ENDPOINT;

export const getDeployment = async (id) => {
  const updatedUrl = `${url}/${id}/deploymentspace/deployment`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch deployment');
  }
};

export const getDeploymentByName = async (id, name) => {
  const updatedUrl = `${url}/${id}/deploymentspace/deployment/${name}`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch deployment details');
  }
};
