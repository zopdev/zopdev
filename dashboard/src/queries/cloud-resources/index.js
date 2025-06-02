import { fetchData, postData, putData } from '@/services/api.js';
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
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ cloudAccId, ...payload }) => {
      const response = await postData(`/cloud-account/${cloudAccId}/resources/state`, payload);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['cloudResourcesGetData'],
      });
    },
  });
}

export function useGetResourceGroup(id, options = {}) {
  return useQuery({
    queryKey: ['ResourceGroupGetData', id],
    queryFn: async () => {
      const url = `/cloud-account/${id}/resource-groups`;
      const response = await fetchData(url);
      return response?.data;
    },
    ...options,
  });
}

export function usePostResourceGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ cloudAccId, ...details }) => {
      const response = await postData(`/cloud-account/${cloudAccId}/resource-groups`, details);
      return response?.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['ResourceGroupGetData'],
      });
    },
  });
}

export function usePutResourceGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ cloudAccId, id, ...details }) => {
      const response = await putData(`/cloud-account/${cloudAccId}/resource-groups/${id}`, details);
      return response?.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['ResourceGroupGetData'],
      });
    },
  });
}

export function usePostResourceGroupSync() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ cloudAccId }) => {
      const response = await postData(`/cloud-account/${cloudAccId}/resources/sync`);
      return response?.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['cloudResourcesGetData'],
      });
    },
  });
}
