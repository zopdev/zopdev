import axiosInstance from './axiosInstance';

// Generic API methods
const apiClient = {
  get: async (url, params = {}) => {
    const response = await axiosInstance.get(url, { params });
    return response.data;
  },

  post: async (url, data = {}) => {
    const response = await axiosInstance.post(url, data);
    return response.data;
  },

  put: async (url, data = {}) => {
    const response = await axiosInstance.put(url, data);
    return response.data;
  },

  patch: async (url, data = {}) => {
    const response = await axiosInstance.patch(url, data);
    return response.data;
  },

  delete: async (url) => {
    const response = await axiosInstance.delete(url);
    return response.data;
  },
};

export default apiClient;
