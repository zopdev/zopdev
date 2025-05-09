import { fetchData, postData } from '@/services/api.js';
import { useMutation, useQuery } from '@tanstack/react-query';

export function useGetCloudAccountsWithAudit(reqParams, options = {}) {
  return useQuery({
    queryKey: ['cloudAccountsWithAudit', reqParams],
    queryFn: async () => {
      const url = `/cloud-accounts`;
      const cloudAccountsResponse = await fetchData(url, options);

      if (!cloudAccountsResponse?.data || !Array.isArray(cloudAccountsResponse.data)) {
        return cloudAccountsResponse;
      }

      const accountsWithAudit = await Promise.all(
        cloudAccountsResponse.data.map(async (account) => {
          try {
            const auditUrl = `/audit/cloud-accounts/${account.id}/results`;
            const auditResponse = await fetchData(auditUrl, options);

            return {
              ...account,
              auditDetails: {
                data: auditResponse?.data || null,
                error: null,
                status: true, // success
              },
            };
          } catch (error) {
            console.error(`Failed to fetch audit data for account ${account.id}:`, error);
            return {
              ...account,
              auditDetails: {
                data: null,
                error: {
                  message: error.message || 'Failed to fetch audit data',
                  status: error.response?.status,
                  details: error.response?.data,
                },
                status: false, // failure
              },
            };
          }
        }),
      );

      return {
        ...cloudAccountsResponse,
        data: accountsWithAudit,
      };
    },
    staleTime: 0,
    cacheTime: 0,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    retry: false,
    ...options,
  });
}

export function usePostAuditData() {
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
  });
}
