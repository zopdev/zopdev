import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL, // Replace with your API URL
  timeout: 10000, // 10 seconds timeout
  headers: {
    'Content-Type': 'application/json',
  },
});

// if (!axiosInstance.defaults.baseURL) {
//   throw new Error(
//     "Base URL is not set in Axios instance. Please provide a valid baseURL."
//   );
// }

// Interceptors for requests (optional)
axiosInstance.interceptors.request.use(
  (config) => {
    // Add any global request modifications (e.g., tokens)
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// Interceptors for responses (optional)
axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle errors globally
    console.error('API error:', error);
    return Promise.reject(error);
  },
);

export default axiosInstance;
