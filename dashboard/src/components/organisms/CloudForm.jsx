import { useEffect, useState } from 'react';
import Button from '@/components/atom/Button';
import Label from '@/components/atom/Loaders';
import Input from '@/components/atom/Input';
import { PROVIDER_ICON_MAPPER } from '@/utils/componentMapper';
import Textarea from '@/components/atom/Textarea';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import { InformationCircleIcon } from '@heroicons/react/24/outline';
import CloudAccountCreationGuide from '@/components/molecules/SetupGuides/CloudAccountSetupGuide';
import AwsAccountCreationGuide from '@/components/molecules/SetupGuides/AwsAccountCreationGuid';
import NameValidation from '@/components/molecules/Validations/NameValidation';

const isValidJSON = (str) => {
  try {
    const parsed = JSON.parse(str);
    return typeof parsed === 'object' && !Array.isArray(parsed) && Object.keys(parsed).length > 0;
  } catch {
    return false;
  }
};

const useAddCloudAccount = (initialValues = {}) => {
  const [values, setValues] = useState({
    name: '',
    credentials: '',
    ...initialValues,
  });
  return { values, setValues };
};

const providerFieldsMap = {
  gcp: (values, handleChange) => (
    <div className="col-span-full">
      <div className="flex justify-between items-center flex-wrap">
        <Label htmlFor="credentials">Service Account Credential (JSON)</Label>
        <FullScreenOverlay
          customCTA={
            <div className="cursor-pointer group hover:text-primary-500 flex items-center gap-1 text-gray-600 group-hover:text-primary-600">
              <InformationCircleIcon className="w-5 h-5" />
              <span>Setup Guide</span>
            </div>
          }
          title="Cloud Account Setup Guide"
          size="4xl"
          maxHeight="90vh"
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
          error={values?.credentials && !isValidJSON(values?.credentials)}
          helperText={
            values?.credentials && !isValidJSON(values?.credentials) ? 'Not a valid JSON' : ' '
          }
        />
      </div>
    </div>
  ),

  aws: (values, handleChange, isEdit) => (
    <>
      <div className="col-span-full mb-4">
        <div className="flex justify-between items-center flex-wrap mb-2">
          <Label htmlFor="aws_access_key_id">AWS Access Key ID</Label>
          <FullScreenOverlay
            customCTA={
              <div className="cursor-pointer group hover:text-primary-500 flex items-center gap-1 text-gray-600 group-hover:text-primary-600">
                <InformationCircleIcon className="w-5 h-5" />
                <span>Setup Guide</span>
              </div>
            }
            title="Cloud Account Setup Guide"
            size="4xl"
            maxHeight="90vh"
            renderContent={AwsAccountCreationGuide}
          />
        </div>
        <Input
          required={!isEdit}
          variant="outlined"
          id="aws_access_key_id"
          name="aws_access_key_id"
          placeholder="Enter access key id"
          onChange={handleChange}
          value={values?.aws_access_key_id || ''}
          helperText=" "
        />
      </div>
      <div className="col-span-full">
        <Label htmlFor="aws_secret_access_key">AWS Secret Access Key</Label>
        <div className="mt-2">
          <Input
            required={!isEdit}
            variant="outlined"
            id="aws_secret_access_key"
            name="aws_secret_access_key"
            placeholder="Enter secret access key"
            onChange={handleChange}
            value={values?.aws_secret_access_key || ''}
            helperText=" "
          />
        </div>
      </div>
    </>
  ),

  azure: (values, handleChange) => (
    <>
      <div className="col-span-full mb-4">
        <div className="flex justify-between items-center flex-wrap">
          <Label htmlFor="projectName">Resource Group</Label>
          <FullScreenOverlay
            isGuide
            RenderIcon={InformationCircleIcon}
            popupTitle="Service Account Setup Guide"
            title="Cloud Account Setup Guide"
          />
        </div>
        <Input
          required
          id="projectName"
          variant="outlined"
          name="projectName"
          placeholder="Enter resource group"
          onChange={handleChange}
          value={values?.projectName || ''}
          helperText=" "
        />
      </div>
      {[
        { id: 'subscriptionId', label: 'Subscription Key' },
        { id: 'tenantId', label: 'Tenant Id' },
        { id: 'appId', label: 'Client Id' },
        { id: 'password', label: 'Client Password' },
      ].map(({ id, label }) => (
        <div key={id} className="sm:col-span-3 mb-4">
          <Label htmlFor={id}>{label}</Label>
          <div className="mt-2">
            <Input
              required
              variant="outlined"
              id={id}
              name={id}
              placeholder={`Enter ${label.toLowerCase()}`}
              onChange={handleChange}
              value={values?.[id] || ''}
              helperText=" "
            />
          </div>
        </div>
      ))}
    </>
  ),
};

const isFormValidByProvider = (provider, values) => {
  const checks = {
    gcp: () => values.name && values.credentials && isValidJSON(values.credentials),
    aws: () => values.aws_access_key_id && values.aws_secret_access_key,
    azure: () =>
      values.projectName &&
      values.subscriptionId &&
      values.tenantId &&
      values.appId &&
      values.password,
  };
  return checks[provider]?.() ?? false;
};

const CloudForm = ({
  providers,
  activeProvider,
  handleProviderChange,
  data,
  updateData,
  audit,
  setIsComplete,
  provider,
  isEdit = false,
}) => {
  const { values, setValues } = useAddCloudAccount(data);

  useEffect(() => {
    if (audit) {
      const valid = isFormValidByProvider(provider, values);
      updateData({ ...values, provider });
      setIsComplete(valid);
    }
  }, [values, audit, provider]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setValues((prev) => ({ ...prev, [name]: value }));
  };

  return (
    <form onSubmit={(e) => e.preventDefault()}>
      <div className="grid grid-cols-1 gap-x-10 sm:max-w-xl sm:grid-cols-6 ml-5">
        {/* Provider Selection */}
        <div className="col-span-full flex gap-4 mb-8 flex-wrap">
          {providers?.map((providerItem) => (
            <Button
              key={providerItem.key}
              variant={activeProvider.key === providerItem.key ? 'primary' : 'secondary'}
              onClick={() => handleProviderChange(providerItem)}
              className="flex items-center"
              startEndornment={PROVIDER_ICON_MAPPER[providerItem.key.toLowerCase()]}
            >
              {providerItem.label}
            </Button>
          ))}
        </div>

        {/* Common Name Field */}
        <div className="col-span-full mb-4">
          <Label htmlFor="name">Name</Label>
          <div className="mt-2">
            <Input
              required
              variant="outlined"
              id="cloud-account-name"
              name="name"
              placeholder="Enter name"
              value={values?.name}
              onChange={handleChange}
              helperText={<NameValidation value={values?.name} min={2} max={255} />}
              inputProps={{ maxLength: 255, minLength: 2 }}
              className="focus:outline-none focus:ring-1 focus:ring-primary-500"
            />
          </div>
        </div>

        {/* Provider Specific Fields */}
        {providerFieldsMap[provider]?.(values, handleChange, isEdit)}
      </div>
    </form>
  );
};

export default CloudForm;
