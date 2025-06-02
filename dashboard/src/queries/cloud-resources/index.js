import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function useGetCloudResources(id, options = {}) {
  return useQuery({
    queryKey: ['cloudResourcesGetData', id],
    queryFn: async () => {
      const url = `/cloud-account/${id}/resources`;
      const response = await fetchData(url);
      return response?.data;
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

export function useGetResourceGroup(id, options = {}) {
  return useQuery({
    queryKey: ['ResourceGroupGetData', id],
    queryFn: async () => {
      const url = `/cloud-account/${id}/resource-groups`;
      return await fetchData(url);
    },
    ...options,
  });
}

export function usePostResourceGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ cloudAccId, resourceIds, ...details }) => {
      const response = await postData(`/cloud-account/${cloudAccId}/resource-groups`, details);
      const assignPromises = resourceIds.map((resourceId) =>
        postData(`/resource-groups/${response?.data?.id}/resources/${resourceId?.id}`, {}),
      );
      await Promise.all(assignPromises);
      return response?.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['ResourceGroupGetData'],
      });
    },
  });
}
