import apiClient from '../../utlis/apiClient';
import { ENVIRONMENT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = ENVIRONMENT_ENDPOINT;

export const getService = async (id) => {
  const updatedUrl = `${url}/${id}/deploymentspace/service`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch services');
  }
};

export const getServiceByName = async (id, name) => {
  const updatedUrl = `${url}/${id}/deploymentspace/service/${name}`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch service details');
  }
};
