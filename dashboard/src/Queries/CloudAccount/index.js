import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery } from '@tanstack/react-query';

export function useImageServiceGetQueryById(reqParams, options = {}) {
  return useQuery({
    queryKey: ['resourceAuditGetData', reqParams],
    queryFn: async () => {
      const { serviceGroupId } = reqParams;
      const url = `/resource/${serviceGroupId}`;
      return fetchData(url);
    },
    staleTime: 0,
    cacheTime: 0,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    retry: false,
    ...options,
  });
}

export function useCreateResourceMutation(options = {}) {
  return useMutation({
    mutationFn: async (payload) => {
      const url = '/cloud-accounts';
      return postData(url, payload);
    },
    ...options,
  });
}
