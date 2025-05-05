import apiClient from '../../utlis/apiClient';
import { ENVIRONMENT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = ENVIRONMENT_ENDPOINT;

export const getPods = async (id) => {
  const updatedUrl = `${url}/${id}/deploymentspace/pod`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch pods');
  }
};

export const getPodByName = async (id, name) => {
  const updatedUrl = `${url}/${id}/deploymentspace/pod/${name}`;
  try {
    const data = await apiClient.get(updatedUrl);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch pod details');
  }
};
