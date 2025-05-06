import { API_URL } from '@/utils/constant.js';

const BASE_URL = API_URL;

const createRequestOptions = (method, data, customHeaders = {}) => {
  const options = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...customHeaders,
    },
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  return options;
};

export const fetchData = async (endpoint, customHeaders = {}) => {
  const response = await fetch(
    `${BASE_URL}${endpoint}`,
    createRequestOptions('GET', null, customHeaders),
  );

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
};

export const postData = async (endpoint, data, customHeaders = {}) => {
  const response = await fetch(
    `${BASE_URL}${endpoint}`,
    createRequestOptions('POST', data, customHeaders),
  );

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
};

export const putData = async (endpoint, data, customHeaders = {}) => {
  const response = await fetch(
    `${BASE_URL}${endpoint}`,
    createRequestOptions('PUT', data, customHeaders),
  );

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
};

export const patchData = async (endpoint, data, customHeaders = {}) => {
  const response = await fetch(
    `${BASE_URL}${endpoint}`,
    createRequestOptions('PATCH', data, customHeaders),
  );

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  return response.json();
};

export const deleteData = async (endpoint, customHeaders = {}) => {
  const response = await fetch(
    `${BASE_URL}${endpoint}`,
    createRequestOptions('DELETE', null, customHeaders),
  );

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
};
