import React, { useEffect, useState } from 'react';
import Textarea from '../../components/TextArea';
import Button from '@/components/atom/Button/index.jsx';
import Label from '@/components/atom/Loaders/index.jsx';
import Input from '@/components/atom/Input/index.jsx';
import { PROVIDER_ICON_MAPPER } from '@/componentMapper.jsx';
import { enforceCharLimit } from '@/utils/index.js';

export const isValidJSON = (str) => {
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

const useAddCloudAccount = (initialValues = {}) => {
  const [values, setValues] = useState({
    name: '',
    credentials: '',
    ...initialValues,
  });

  return { values, setValues };
};

const CloudForm = ({
  provider,
  options,
  tabValue,
  handleTabChange,
  data,
  updateData,
  audit,
  setIsComplete,
}) => {
  const { values, setValues } = useAddCloudAccount(data);

  const isFormValid = () => {
    return values?.name && values?.credentials && isValidJSON(values?.credentials);
  };

  useEffect(() => {
    if (audit) {
      const valid = isFormValid();
      updateData({ ...data, ...values });
      if (valid) setIsComplete(true);
      else setIsComplete(false);
    }
  }, [values, audit]);

  const handleChange = (e) => {
    const { name, value } = e.target;

    setValues((prev) => {
      if (name === 'credentials') {
        if (isValidJSON(value)) {
          return {
            ...prev,
            [name]: JSON.stringify(JSON.parse(value), null, 4),
          };
        }
        return { ...prev, [name]: value };
      }

      if (name === 'name' && provider === 'azure') {
        return { ...prev, [name]: enforceCharLimit(value, 11) };
      }

      return { ...prev, [name]: value };
    });
  };

  return (
    <form onSubmit={(e) => e.preventDefault()}>
      <div className="grid grid-cols-1 gap-x-10 sm:max-w-xl sm:grid-cols-6 ml-5">
        <div className="col-span-full flex gap-4 mb-8 flex-wrap">
          {options?.map((item, idx) => (
            <Button
              key={idx}
              variant={tabValue === idx ? 'primary' : 'secondary'}
              onClick={() => handleTabChange(idx)}
              className="pointer-events-none"
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

        {provider === 'gcp' && (
          <div className="col-span-full">
            <Label htmlFor="credentials">Service Account Credential (JSON)</Label>
            <div className="mt-2">
              <Textarea
                rows={6}
                placeholder="Enter GCP credentials"
                multiline="true"
                id="credentials"
                variant="outlined"
                name="credentials"
                className="focus:outline-none focus:ring-1 focus:ring-primary-500"
                onChange={handleChange}
                value={values?.credentials}
                error={values.credentials ? !isValidJSON(values.credentials) : false}
                helperText={
                  values.credentials && !isValidJSON(values.credentials) ? 'Not a valid JSON' : ' '
                }
              />
            </div>
          </div>
        )}
      </div>
    </form>
  );
};

export default CloudForm;
