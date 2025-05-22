import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery } from '@tanstack/react-query';

export function useGetCloudResources(id, options = {}) {
  return useQuery({
    queryKey: ['cloudResourcesGetData', id],
    queryFn: async () => {
      const url = `/resources?cloudAccId=${id}`;
      return await fetchData(url);
    },
    ...options,
  });
}

export function usePostResourceState() {
  return useMutation({
    mutationFn: async (payload) => {
      const response = await postData('/resources/state', payload);
      return response;
    },
  });
}
