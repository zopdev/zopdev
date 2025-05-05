import { useState } from 'react';
import { addCloudAccount } from '../../Queries/CloudAccount';
import { useRouter } from 'next/navigation';
import { ParseJSON } from '../../utlis/helperFunc';

const formateData = async (values, provider) => {
  const credentials = {
    aws: {
      aws_access_key_id: values.aws_access_key_id,
      aws_secret_access_key: values.aws_secret_access_key,
    },
    azure: {
      appId: values.appId,
      password: values.password,
      tenantId: values.tenantId,
      //   subscriptionId: values.subscriptionId,
    },
    gcp: provider === 'gcp' ? ParseJSON(values.credentials) : {},
  };

  const finalValues = {
    // orgId: tokenInfo?.["tenant-id"],
    name: values.name,
    provider,
    // configs: {
    //   name:
    //     provider !== "aws"
    //       ? provider === "gcp"
    //         ? credentials[provider].project_id
    //         : values.projectName
    //       : values.aws_account_name,
    //   credentials: credentials[provider],
    // },
    credentials: credentials[provider],
  };

  return finalValues;
};

const useAddCloudAccount = () => {
  const router = useRouter();
  const [values, setValues] = useState({
    name: '',
    projectName: '',
    credentials: '',
    aws_account_name: '',
    aws_access_key_id: '',
    aws_secret_access_key: '',
    appRegion: '',
    subscriptionId: '',
    tenantId: '',
    appId: '',
    password: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // Submit form
  const handleSubmit = async (values, provider) => {
    setIsLoading(true);
    setError(null);
    const reqBody = await formateData(values, provider);

    const postCloudAccount = async (values) => {
      try {
        const data = await addCloudAccount(values);
        setError(null); // Clear any previous errors
        router.push('/cloud-accounts'); // Redirect only if no error
        return data;
      } catch (error) {
        setError(error.message);
      } finally {
        setIsLoading(false);
      }
    };

    postCloudAccount(reqBody);
  };

  return {
    values,
    setValues,
    handleSubmit,
    isLoading,
    error,
  };
};

export default useAddCloudAccount;
