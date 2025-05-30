import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery } from '@tanstack/react-query';

export function useGetCloudResources(id, options = {}) {
  return useQuery({
    queryKey: ['cloudResourcesGetData', id],
    queryFn: async () => {
      const url = `/cloud-account/${id}/resources`;
      return await fetchData(url);
    },
    ...options,
  });
}

export function usePostResourceState() {
  return useMutation({
    mutationFn: async ({ cloudAccId, ...payload }) => {
      const response = await postData(`/cloud-account/${cloudAccId}/resources/state`, payload);
      return response;
    },
  });
}
