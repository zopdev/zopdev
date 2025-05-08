import { API_BASE_URL } from '@/utils/constant.js';
import { HttpErrors } from './errors.js';

const BASE_URL = API_BASE_URL;

const createRequestOptions = (method, data, customHeaders = {}) => {
  const headers = {
    'Content-Type': 'application/json',
    ...customHeaders,
  };

  const options = {
    method,
    headers,
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  return options;
};

const getErrorMessage = (status) => {
  const messages = {
    500: 'Internal Server Error. Please try again later.',
    501: 'Not Implemented. This feature is coming soon.',
    502: 'Bad Gateway. Please check your connection.',
    503: 'Service Unavailable. Try again in a few minutes.',
    504: 'Gateway Timeout. The server took too long to respond.',
    505: 'HTTP Version Not Supported.',
  };

  return messages[status] || 'An unexpected server error occurred.';
};

const safeParseJSON = async (response) => {
  try {
    return await response.json();
  } catch {
    return {};
  }
};

const handleErrorResponse = async (response) => {
  const json = await safeParseJSON(response);
  const is5xx = response.status >= 500 && response.status < 600;

  const message = is5xx
    ? getErrorMessage(response.status)
    : json?.error.message || `API error: ${response.status}`;

  throw new HttpErrors({ message, details: json }, response.status);
};

const request = async (method, endpoint, data = null, customHeaders = {}) => {
  try {
    const response = await fetch(
      `${BASE_URL}${endpoint}`,
      createRequestOptions(method, data, customHeaders),
    );

    if (!response.ok) {
      await handleErrorResponse(response);
    }

    return response.status === 204 ? null : response.json();
  } catch (err) {
    if (err.name === 'AbortError') {
      throw new HttpErrors(
        { message: 'Request timed out. Please check your internet connection.' },
        408,
      );
    }

    if (err instanceof TypeError) {
      const msg = (err.message || '').toLowerCase();

      if (msg.includes('failed to fetch')) {
        throw new HttpErrors(
          { message: 'Unable to reach server. This may be due to no internet or CORS issues.' },
          0,
        );
      }

      if (msg.includes('networkerror')) {
        throw new HttpErrors(
          { message: 'Network error. Please check your connection or try again later.' },
          0,
        );
      }

      // fallback
      throw new HttpErrors({ message: 'A network issue occurred.' }, 0);
    }

    // unexpected fetch-related errors
    throw new HttpErrors({ message: err.message || 'Unexpected error occurred.' }, 0);
  }
};

// Public API
export const fetchData = (endpoint, headers = {}) => request('GET', endpoint, null, headers);

export const postData = (endpoint, data, headers = {}) => request('POST', endpoint, data, headers);

export const putData = (endpoint, data, headers = {}) => request('PUT', endpoint, data, headers);

export const patchData = (endpoint, data, headers = {}) =>
  request('PATCH', endpoint, data, headers);

export const deleteData = (endpoint, headers = {}) => request('DELETE', endpoint, null, headers);
