import React from 'react';
import Label from '../../components/Label';
import Input from '../../components/Input';
import useAddEnvironment from '../../hooks/environment/addEnvironment';
import Button from '../../components/Button';
import ErrorComponent from '../../components/ErrorComponent';
import { useParams, usePathname, useRouter } from 'next/navigation';

const AddEnvironment = () => {
  const router = useParams();
  const { values, setValues, handleSubmit, isLoading, error } = useAddEnvironment();

  const handleChange = (e) => {
    setValues({ ...values, [e.target.name]: e.target.value });
  };

  const handleOnSubmit = (e) => {
    e.preventDefault();
    handleSubmit(router?.['application-id'], values);
  };

  return (
    <form onSubmit={handleOnSubmit}>
      <div className="grid grid-cols-1 gap-x-6 gap-y-0 sm:max-w-xl sm:grid-cols-6">
        <div className="col-span-full">
          <Label htmlFor="name">Name</Label>
          <div className="mt-2">
            <Input
              required
              variant="outlined"
              id="environment-name"
              onChange={handleChange}
              name="name"
              placeholder="Enter environment"
              helperText={' '}
              value={values?.name}
              //   disabled={isEditableFields ? isEditableFields.includes('name') : isLoading}
              //   helperText={<NameValidation value={values.name} type={provider} />}
              //   inputProps={provider === 'azure' ? { minLength: 6 } : { maxLength: 16, minLength: 6 }}
            />
          </div>
        </div>

        <div className="col-span-full mb-2">
          {error && (
            <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
          )}
        </div>
      </div>
      <div className="mt-2 flex gap-4">
        <Button type="submit" id="connect-provider" disabled={isLoading}>
          Add
          {isLoading && (
            <div
              style={{ borderTopColor: 'transparent' }}
              className="w-4 h-4 border-4 border-blue-200 rounded-full animate-spin"
            />
          )}
        </Button>
      </div>
    </form>
  );
};

export default AddEnvironment;
