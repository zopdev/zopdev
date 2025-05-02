import apiClient from '../../utlis/apiClient';
import { APPLICATION_ENDPOINT } from '../../utlis/apiEndpoints';

const url = APPLICATION_ENDPOINT;

export const getApplication = async () => {
  try {
    const data = await apiClient.get(url);
    return data;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch applications');
  }
};

export const addEnvironment = async (id, values) => {
  const updatedUrl = `${url}/${id}/environments`;
  try {
    const response = await apiClient.post(updatedUrl, values);
    return response;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to add environment');
  }
};

export const getApplicationById = async (id) => {
  const updatedUrl = `${url}/${id}`;
  try {
    const response = await apiClient.get(updatedUrl);
    return response;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to fetch application');
  }
};

export const addApplication = async (values) => {
  try {
    const response = await apiClient.post(url, values);
    return response;
  } catch (error) {
    throw new Error(error.response?.data?.error?.message || 'Failed to add application');
  }
};
