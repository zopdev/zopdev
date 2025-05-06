import React, { useEffect, useState } from 'react';
import Textarea from '../../components/TextArea';
import Button from '@/components/atom/Button/index.jsx';
import Label from '@/components/atom/Loaders/index.jsx';
import Input from '@/components/atom/Input/index.jsx';
import ErrorComponent from '@/components/atom/ErrorComponent/index.jsx';
import { PROVIDER_ICON_MAPPER } from '@/componentMapper.jsx';

export const isValidJSON = (str) => {
  try {
    JSON.parse(str);
  } catch (e) {
    return false;
  }
  const decryptLog = JSON.parse(str);

  return (
    typeof decryptLog === 'object' &&
    !Array.isArray(decryptLog) &&
    Object.keys(decryptLog).length > 0 &&
    Object.keys(decryptLog).every((key) => !!key)
  );
};

const useAddCloudAccount = () => {
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

  return {
    values,
    setValues,
    isLoading,
    error,
  };
};

const CloudForm = ({
  provider,
  value,
  handleTabChange,
  options,
  tabValue,
  audit,
  setIsComplete,
  updateData,
  data,
}) => {
  const { values, setValues, isLoading, error } = useAddCloudAccount();
  const isFormValid = (values, provider) => {
    if (audit === true) {
      if (!values?.name) return false;

      if (provider === 'gcp') {
        return values?.credentials && isValidJSON(values?.credentials);
      }
      if (provider === 'aws') {
        return values?.aws_access_key_id && values?.aws_secret_access_key;
      }
      if (provider === 'azure') {
        return values?.projectName && values?.tenantId && values?.appId && values?.password;
      }
      return false;
    }
  };

  useEffect(() => {
    if (audit === true) {
      const valid = isFormValid(values, provider);
      const updatedData = { ...data, [name]: value };
      updateData(updatedData);
      if (valid) {
        setIsComplete(true);
        updateData(values);
      }
    }
  }, [values, audit, provider]);

  const handleChange = (e) => {
    if (e.target.name === 'credentials' && isValidJSON(e.target.value)) {
      setValues((prevValues) => ({
        ...prevValues,
        [e.target.name]: JSON.stringify(JSON.parse(e.target.value), undefined, 4),
      }));
    } else if (e.target.name === 'name' && provider === 'azure') {
      setValues((prevValues) => ({
        ...prevValues,
        [e.target.name]: enforceCharLimit(e.target.value, 11),
      }));
    } else {
      setValues((prevValues) => ({
        ...prevValues,
        [e.target.name]: e.target.value,
      }));
    }
  };

  useEffect(() => {
    if (
      value === 1 &&
      values?.aws_access_key_id?.length !== 0 &&
      values?.aws_secret_access_key?.length !== 0
    ) {
      const data = {
        aws_access_key_id: values?.aws_access_key_id,
        aws_secret_access_key: values?.aws_secret_access_key,
      };

      handleGetRegion(data, provider);
    }
    if (value === 0 && values?.credentials?.length !== 0 && isValidJSON(values?.credentials)) {
      handleGetRegion({ credentials: values?.credentials }, provider);
    }
    if (
      value === 2 &&
      values?.subscriptionId?.length !== 0 &&
      values?.tenantId?.length !== 0 &&
      values?.appId?.length !== 0 &&
      values?.password?.length !== 0 &&
      values?.projectName?.length !== 0
    ) {
      const data = {
        subscriptionId: values?.subscriptionId,
        tenantId: values?.tenantId,
        appId: values?.appId,
        password: values?.password,
        projectName: values?.projectName,
      };
      handleGetRegion(data, provider);
    }
  }, [value]);

  const handleOnSubmit = (e) => {
    e.preventDefault();
  };

  return (
    <form onSubmit={handleOnSubmit}>
      <div
        className={`grid grid-cols-1 gap-x-10 gap-y-0 sm:max-w-xl sm:grid-cols-6 ${audit ? 'ml-5' : 'ml-5'} `}
      >
        <div className="col-span-full flex gap-4 mb-8 flex-wrap">
          {options?.map((item, idx) => (
            <Button
              key={idx}
              variant={tabValue === idx ? 'primary' : 'secondary'}
              onClick={() => handleTabChange(idx)}
              className={' pointer-events-none '}
            >
              {PROVIDER_ICON_MAPPER[item.label.toLowerCase()]}
              {item.label}
            </Button>
          ))}
        </div>

        <div className="col-span-full">
          <Label htmlFor="name">Name</Label>
          <div className="mt-2">
            <Input
              required
              variant="outlined"
              id="cloud-account-name"
              onChange={handleChange}
              name="name"
              placeholder="Enter name"
              value={data?.name || values?.name}
              helperText={
                values?.name?.length < 255
                  ? ' '
                  : 'Maximum 255 characters are allowed, limit reached'
              }
              inputProps={{ maxLength: 255, minLength: 1 }}
              helperTextClass={' text-yellow-500'}
              className={'focus:outline-none focus:ring-1 focus:ring-primary-500'}
            />
          </div>
        </div>

        {provider === 'gcp' && (
          <div className="col-span-full">
            <div className="flex justify-between items-center flex-wrap">
              <Label htmlFor="credentials">Service Account Credential (JSON)</Label>
            </div>
            <div className="mt-2">
              <Textarea
                rows={6}
                placeholder="Enter GCP credentials"
                multiline="true"
                id="credentials"
                variant="outlined"
                name="credentials"
                className={' focus:outline-none focus:ring-1 focus:ring-primary-500'}
                onChange={handleChange}
                value={data?.credentials || values?.credentials}
                error={values?.credentials ? !isValidJSON(values?.credentials) : false}
                helperText={
                  values?.credentials && !isValidJSON(values?.credentials)
                    ? 'Not a valid JSON'
                    : ' '
                }
                disabled={isLoading}
              />
            </div>
          </div>
        )}
        {provider === 'aws' && (
          <>
            <div className="col-span-full">
              <Label htmlFor="credentials">AWS Access Key ID</Label>
              <div className="mt-2">
                <Input
                  variant="outlined"
                  id="aws_access_key_id"
                  onChange={handleChange}
                  value={values?.aws_access_key_id}
                  placeholder="Enter access key id"
                  name="aws_access_key_id"
                  helperText={' '}
                />
              </div>
            </div>
            <div className="col-span-full">
              <Label htmlFor="credentials">AWS Secret Access Key</Label>
              <div className="mt-2">
                <Input
                  variant="outlined"
                  id="aws_secret_access_key"
                  placeholder="Enter secret access key"
                  onChange={handleChange}
                  value={values?.aws_secret_access_key}
                  name="aws_secret_access_key"
                  helperText={' '}
                />
              </div>
            </div>
          </>
        )}
        {provider === 'azure' && (
          <>
            <div className="col-span-full">
              <div className="flex justify-between items-center flex-wrap">
                <Label htmlFor="credentials">Resource Group</Label>
              </div>

              <div className="mt-2">
                <Input
                  required
                  id="projectName"
                  variant="outlined"
                  placeholder="Enter resource group"
                  onChange={handleChange}
                  name="projectName"
                  value={values?.projectName}
                  helperText={' '}
                />
              </div>
            </div>
            <div className="sm:col-span-3">
              <Label htmlFor="credentials">Tenant Id</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="tenantId"
                  placeholder="Enter tenant id"
                  onChange={handleChange}
                  name="tenantId"
                  value={values?.tenantId}
                  helperText={' '}
                />
              </div>
            </div>
            <div className="sm:col-span-3">
              <Label htmlFor="credentials">Client Id</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="appId"
                  placeholder="Enter application id"
                  onChange={handleChange}
                  name="appId"
                  value={values?.appId}
                  helperText={' '}
                />
              </div>
            </div>
            <div className="sm:col-span-3">
              <Label htmlFor="credentials">Client Password</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="password"
                  placeholder="Enter password"
                  onChange={handleChange}
                  name="password"
                  value={values?.password}
                  helperText={' '}
                />
              </div>
            </div>
          </>
        )}
      </div>

      <div className="text-sm text-red-500 mt-5 grid grid-cols-1 gap-x-6 gap-y-0 sm:max-w-xl sm:grid-cols-6">
        <div className="col-span-full">
          {error && <ErrorComponent errorText={error} className={' !p-2'} />}
        </div>
      </div>
      {!audit && (
        <div className="mt-8 flex gap-4">
          <Button
            type="submit"
            id="connect-provider"
            disabled={
              (values?.credentials ? !isValidJSON(values?.credentials) : false) || isLoading
            }
          >
            Add
            {isLoading && (
              <div
                style={{ borderTopColor: 'transparent' }}
                className="w-4 h-4 border-4 border-blue-200 rounded-full animate-spin"
              />
            )}
          </Button>
        </div>
      )}
    </form>
  );
};

export default CloudForm;
