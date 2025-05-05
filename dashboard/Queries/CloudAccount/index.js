import apiClient from '../../utlis/apiClient';
import { CLOUD_ACCOUNT_ENDPOINT } from '../../utlis/apiEndpoints';

const url = CLOUD_ACCOUNT_ENDPOINT;

export const getCloudAccounts = async () => {
  try {
    const data = await apiClient.get(url);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch cloud accounts');
  }
};

export const addCloudAccount = async (values) => {
  try {
    const response = await apiClient.post(url, values);
    return response;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to add cloud account');
  }
};
