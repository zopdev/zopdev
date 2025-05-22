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
  // const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload) => {
      const response = await postData('/resources/state', payload);
      return response;
    },
    // onSuccess: (_data, variables) => {
    //   // Invalidate or refetch any related queries if needed
    //   queryClient.invalidateQueries({
    //     queryKey: ['resourceState', { name: variables?.name }],
    //   });
    // },
  });
}
