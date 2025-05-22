import { ParseJSON } from '@/utils/common';

export function transformResourceAuditPayload({ 0: CloudDetails }) {
  const { name, provider, ...rest } = CloudDetails;

  const credentials = {
    aws: {
      aws_access_key_id: rest.aws_access_key_id,
      aws_secret_access_key: rest.aws_secret_access_key,
    },
    azure: {
      appId: rest.appId,
      password: rest.password,
      tenantId: rest.tenantId,
      //   subscriptionId: values.subscriptionId,
    },
    gcp: provider === 'gcp' ? ParseJSON(rest.credentials) : {},
  };

  const finalValues = {
    name,
    provider,
    credentials: credentials[provider],
  };

  return finalValues;
}
