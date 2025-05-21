import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function useGetCloudAccounts(reqParams, options = {}) {
  return useQuery({
    queryKey: ['cloudAccountGetData', reqParams],
    queryFn: async () => {
      const url = `/cloud-accounts`;
      return await fetchData(url);
    },
    ...options,
  });
}

export function useGetAuditDetails(reqParams, options = {}) {
  return useQuery({
    queryKey: ['resourceAuditGetData', reqParams],
    queryFn: async () => {
      const url = `/audit/cloud-accounts/${reqParams?.id}/results`;
      return await fetchData(url);
    },
    ...options,
  });
}

export function usePostAuditData() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (req) => {
      let id = req?.id;
      let createResponse = null;

      // Step 1: Conditionally create cloud account
      if (req?.transformedData) {
        const getCloudAccountRes = await postData('/cloud-accounts', req.transformedData);
        id = getCloudAccountRes?.data?.id;

        if (!id) {
          throw new Error('Missing ID in cloud account creation response');
        }

        createResponse = getCloudAccountRes;
      }

      // Step 2: Conditionally trigger audit
      let auditResponse = null;
      if (id && req?.selectedOption) {
        let auditUrl;
        let auditPayload;
        if (req?.selectedOption === 'run-all' || req?.selectedOption === 'all') {
          auditUrl = `/audit/cloud-accounts/${id}/all`;
          auditPayload = {};
        } else {
          auditUrl = `/audit/cloud-accounts/${id}/category/${req?.selectedOption}`;
          auditPayload = req?.selectedOption;
        }

        auditResponse = await postData(auditUrl, auditPayload);
      }

      return {
        createResponse,
        auditResponse,
      };
    },
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ['resourceAuditGetData', { id: variables?.id }],
      });
    },
  });
}
