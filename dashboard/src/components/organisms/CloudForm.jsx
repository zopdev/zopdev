import React, { useEffect, useState } from 'react';
import Button from '@/components/atom/Button/index.jsx';
import Label from '@/components/atom/Loaders/index.jsx';
import Input from '@/components/atom/Input/index.jsx';
import { PROVIDER_ICON_MAPPER } from '@/utils/componentMapper.jsx';
import Textarea from '@/components/atom/Textarea/index.jsx';
// import { enforceCharLimit } from '@/utils/common.js';
// import Tooltip from '@/components/atom/Tooltip/index.jsx';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay/index.jsx';
import { InformationCircleIcon } from '@heroicons/react/24/outline/index.js';
import CloudAccountCreationGuide from '@/components/molecules/SetupGuides/CloudAccountSetupGuide.jsx';
import NameValidation from '@/components/molecules/Validations/NameValidation.jsx';
import AwsAccountCreationGuide from '../molecules/SetupGuides/AwsAccountCreationGuid';

// Provider-specific form configurations
const CloudForm = ({
  providers,
  activeProvider,
  handleProviderChange,
  data,
  updateData,
  audit,
  setIsComplete,
  provider,
  isEdit = false, // Added isEdit prop with default false
}) => {
  // Initialize state with data for the current provider
  // const [values, setValues] = useState(() => {
  //   // Use existing data or start with empty object
  //   return data || {};
  // });

  const useAddCloudAccount = (initialValues = {}) => {
    const [values, setValues] = useState({
      name: '',
      credentials: '',
      ...initialValues,
    });

    return { values, setValues };
  };
  const { values, setValues } = useAddCloudAccount(data);
  const isValidJSON = (str) => {
    try {
      const parsed = JSON.parse(str);
      return (
        typeof parsed === 'object' &&
        !Array.isArray(parsed) &&
        Object.keys(parsed).length > 0 &&
        Object.keys(parsed).every((key) => !!key)
      );
    } catch {
      return false;
    }
  };
  // Check if form is valid based on provider-specific validations
  const isFormValid = () => {
    switch (provider) {
      case 'gcp':
        return values?.name && values?.credentials && isValidJSON(values?.credentials);

      case 'aws':
        return values?.aws_access_key_id && values?.aws_secret_access_key;

      case 'azure':
        return (
          values?.projectName &&
          values?.subscriptionId &&
          values?.tenantId &&
          values?.appId &&
          values?.password
        );

      default:
        return false;
    }
  };

  useEffect(() => {
    if (audit) {
      const valid = isFormValid();
      updateData({ ...values, provider });
      setIsComplete(valid);
    }
  }, [values, audit, provider]);

  const handleChange = (e) => {
    const { name, value } = e.target;

    setValues((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const renderProviderSpecificFields = () => {
    switch (provider) {
      case 'gcp':
        return (
          <>
            <div className="col-span-full mb-4">
              <Label htmlFor="name">Name</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="cloud-account-name"
                  onChange={handleChange}
                  name="name"
                  placeholder="Enter name"
                  value={values?.name}
                  helperText={<NameValidation value={values?.name} min={2} max={16} />}
                  inputProps={{ maxLength: 16, minLength: 2 }}
                  // inputProps={
                  //   provider === 'azure' ? { minLength: 6 } : { maxLength: 16, minLength: 6 }
                  // }
                />
              </div>
            </div>

            <div className="col-span-full">
              <div className="flex justify-between items-center flex-wrap">
                <Label htmlFor="credentials">Service Account Credential (JSON)</Label>
                <FullScreenOverlay
                  customCTA={
                    <div className={'cursor-pointer group hover:text-primary-500'}>
                      <div
                        className={
                          'flex items-center justify-center text-gray-600 group-hover:text-primary-600 gap-1'
                        }
                      >
                        <InformationCircleIcon className="w-5 h-5 text-gray-600 group-hover:text-primary-600" />
                        <span>Setup Guide</span>
                      </div>
                    </div>
                  }
                  title="Cloud Account Setup Guide"
                  size={'4xl'}
                  maxHeight={'90vh'}
                  renderContent={CloudAccountCreationGuide}
                />
              </div>

              <div className="mt-2">
                <Textarea
                  rows={6}
                  id="credentials"
                  name="credentials"
                  placeholder="Enter GCP credentials JSON"
                  variant="outlined"
                  className="focus:outline-none focus:ring-1 focus:ring-primary-500"
                  onChange={handleChange}
                  value={values?.credentials || ''}
                  error={
                    values?.credentials &&
                    (() => {
                      try {
                        const parsed = JSON.parse(values.credentials);
                        return !(
                          typeof parsed === 'object' &&
                          !Array.isArray(parsed) &&
                          Object.keys(parsed).length > 0
                        );
                      } catch {
                        return true;
                      }
                    })()
                  }
                  helperText={
                    values?.credentials &&
                    (() => {
                      try {
                        const parsed = JSON.parse(values.credentials);
                        return !(
                          typeof parsed === 'object' &&
                          !Array.isArray(parsed) &&
                          Object.keys(parsed).length > 0
                        )
                          ? 'Not a valid JSON'
                          : ' ';
                      } catch {
                        return 'Not a valid JSON';
                      }
                    })()
                  }
                />
              </div>
            </div>
          </>
        );

      case 'aws':
        return (
          <>
            <div className="col-span-full mb-4">
              <Label htmlFor="name">Name</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="cloud-account-name"
                  name="name"
                  onChange={handleChange}
                  placeholder="Enter name"
                  value={values?.name}
                  helperText={
                    values.name?.length < 255
                      ? ' '
                      : 'Maximum 255 characters are allowed, limit reached'
                  }
                  inputProps={{ maxLength: 255, minLength: 1 }}
                  helperTextClass="text-yellow-500"
                  className="focus:outline-none focus:ring-1 focus:ring-primary-500"
                />
              </div>
            </div>
            <div className="col-span-full mb-4">
              <div className="flex justify-between items-center flex-wrap mb-2">
                <Label htmlFor="credentials">AWS Access Key ID</Label>
                <FullScreenOverlay
                  customCTA={
                    <div className={'cursor-pointer group hover:text-primary-500'}>
                      <div
                        className={
                          'flex items-center justify-center text-gray-600 group-hover:text-primary-600 gap-1'
                        }
                      >
                        <InformationCircleIcon className="w-5 h-5 text-gray-600 group-hover:text-primary-600" />
                        <span>Setup Guide</span>
                      </div>
                    </div>
                  }
                  title="Cloud Account Setup Guide"
                  size={'4xl'}
                  maxHeight={'90vh'}
                  renderContent={AwsAccountCreationGuide}
                />
              </div>
              <div>
                <Input
                  required={!isEdit}
                  variant="outlined"
                  id="aws_access_key_id"
                  onChange={handleChange}
                  value={values?.aws_access_key_id || ''}
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
                  required={!isEdit}
                  variant="outlined"
                  id="aws_secret_access_key"
                  placeholder="Enter secret access key"
                  onChange={handleChange}
                  value={values?.aws_secret_access_key || ''}
                  name="aws_secret_access_key"
                  helperText={' '}
                />
              </div>
            </div>
          </>
        );

      case 'azure':
        return (
          <>
            <div className="col-span-full mb-4">
              <Label htmlFor="name">Name</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="cloud-account-name"
                  name="name"
                  onChange={handleChange}
                  placeholder="Enter name"
                  value={values?.name}
                  helperText={
                    values.name?.length < 255
                      ? ' '
                      : 'Maximum 255 characters are allowed, limit reached'
                  }
                  inputProps={{ maxLength: 255, minLength: 1 }}
                  helperTextClass="text-yellow-500"
                  className="focus:outline-none focus:ring-1 focus:ring-primary-500"
                />
              </div>
            </div>
            <div className="col-span-full mb-4">
              <div className="flex justify-between items-center flex-wrap">
                <Label htmlFor="credentials">Resource Group</Label>
                <FullScreenOverlay
                  isGuide
                  RenderIcon={InformationCircleIcon}
                  // RenderComponent={AzureAccountCreationGuide}
                  popupTitle={'Service Account Setup Guide'}
                  title={'Cloud Account Setup Guide'}
                />
              </div>

              <div>
                <Input
                  required
                  id="projectName"
                  variant="outlined"
                  placeholder="Enter resource group"
                  onChange={handleChange}
                  name="projectName"
                  value={values?.projectName || ''}
                  helperText={' '}
                />
              </div>
            </div>
            <div className="sm:col-span-3 mb-4">
              <Label htmlFor="credentials">Subscription Key</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="subscriptionId"
                  placeholder="Enter subscription id"
                  onChange={handleChange}
                  name="subscriptionId"
                  value={values?.subscriptionId || ''}
                  helperText={' '}
                />
              </div>
            </div>
            <div className="sm:col-span-3 mb-4">
              <Label htmlFor="credentials">Tenant Id</Label>
              <div className="mt-2">
                <Input
                  required
                  variant="outlined"
                  id="tenantId"
                  placeholder="Enter tenant id"
                  onChange={handleChange}
                  name="tenantId"
                  value={values?.tenantId || ''}
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
                  value={values?.appId || ''}
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
                  value={values?.password || ''}
                  helperText={' '}
                />
              </div>
            </div>
          </>
        );

      default:
        return null;
    }
  };

  return (
    <form onSubmit={(e) => e.preventDefault()}>
      <div className="grid grid-cols-1 gap-x-10 sm:max-w-xl sm:grid-cols-6 ml-5">
        {/* Provider Selection Buttons */}
        <div className="col-span-full flex gap-4 mb-8 flex-wrap">
          {providers?.map((providerItem) => (
            <Button
              key={providerItem.key}
              variant={activeProvider.key === providerItem.key ? 'primary' : 'secondary'}
              onClick={() => handleProviderChange(providerItem)}
              className="flex items-center"
            >
              {PROVIDER_ICON_MAPPER[providerItem.key.toLowerCase()]}
              {providerItem.label}
            </Button>
          ))}
        </div>

        {/* Dynamic Provider-Specific Fields */}
        {renderProviderSpecificFields()}
      </div>
    </form>
  );
};

export default CloudForm;
