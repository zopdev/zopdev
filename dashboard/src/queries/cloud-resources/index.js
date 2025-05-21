import { fetchData } from '@/services/api.js';
import { useQuery } from '@tanstack/react-query';

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
