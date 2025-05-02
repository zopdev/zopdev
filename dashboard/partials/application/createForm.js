'use client';

import React, { useState } from 'react';
import Label from '../../components/Label';
import Input from '../../components/Input';
import InputWithButton from '../../components/Input/inputWithButton';
import DraggableList from '../draggableList';
import Button from '../../components/Button';
import useAddApplication from '../../hooks/application/addApplication';
import ErrorComponent from '../../components/ErrorComponent';

const CreateAppForm = () => {
  const { values, setValues, handleSubmit, isLoading, error } = useAddApplication();

  const [chips, setChips] = useState([
    { name: 'prod', selected: false },
    { name: 'stage', selected: false },
    { name: 'test', selected: false },
  ]);
  const [inputValue, setInputValue] = useState('');

  // Handle input change
  const handleInputChange = (e) => setInputValue(e.target.value);

  // Add new chip
  const handleAddChip = () => {
    const chipName = inputValue.trim();
    if (chipName && !chips.some((chip) => chip.name === chipName)) {
      setChips([...chips, { name: chipName, selected: true }]);
      setInputValue(''); // Clear input
    }
  };

  // Toggle chip selection
  const toggleChipSelection = (name) => {
    setChips((prevChips) =>
      prevChips.map((chip) => (chip.name === name ? { ...chip, selected: !chip.selected } : chip)),
    );
  };

  const handleChange = (e) => {
    setValues({ ...values, [e.target.name]: e.target.value });
  };

  const handleListUpdate = (updatedList) => {
    setValues({ ...values, environments: updatedList });
  };

  const handleOnSubmit = (e) => {
    e.preventDefault();
    handleSubmit(values);
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
              id="cloud-account-name"
              onChange={handleChange}
              name="name"
              placeholder="Enter name"
              helperText={' '}
              value={values?.name}
              //   disabled={isEditableFields ? isEditableFields.includes('name') : isLoading}
              //   helperText={<NameValidation value={values.name} type={provider} />}
              //   inputProps={provider === 'azure' ? { minLength: 6 } : { maxLength: 16, minLength: 6 }}
            />
          </div>
        </div>
        <div className="col-span-full">
          <Label htmlFor="name">Environment (optional)</Label>
          <div className="mt-2">
            <InputWithButton
              //   required
              variant="outlined"
              id="cloud-account-name"
              onChange={handleInputChange}
              name="name"
              placeholder="Enter name"
              value={inputValue}
              onClick={handleAddChip}
              //   disabled={isEditableFields ? isEditableFields.includes('name') : isLoading}
              //   helperText={<NameValidation value={values.name} type={provider} />}
              //   inputProps={provider === 'azure' ? { minLength: 6 } : { maxLength: 16, minLength: 6 }}
            />
          </div>
        </div>
        <div className="col-span-full mt-6">
          <Label htmlFor="name">Choose Environment</Label>
          <div className="flex flex-wrap gap-3 mt-2">
            {chips.map((chip) => (
              <span
                key={chip.name}
                onClick={() => toggleChipSelection(chip.name)}
                className={`inline-flex items-center rounded-full px-3 py-2 text-xs font-medium cursor-pointer hover:bg-gray-200 ${
                  chip.selected
                    ? 'bg-primary-100 text-primary-600 hover:!bg-primary-200 '
                    : 'bg-gray-100 text-gray-600'
                }`}
              >
                {chip.name}
              </span>
            ))}
          </div>
        </div>

        <div className="col-span-full mt-10">
          {values?.environments?.length > 0 && (
            <Label htmlFor="name">Choose the order of environment</Label>
          )}
          <DraggableList chips={chips} handleListUpdate={handleListUpdate} />
          {values?.environments?.length > 0 && (
            <p className="text-xs italic text-gray-400 mt-2">
              *Rearrange the order by drag and drop
            </p>
          )}
        </div>
      </div>

      <div className="text-sm text-red-500 mt-5 grid grid-cols-1 gap-x-6 gap-y-0 sm:max-w-xl sm:grid-cols-6">
        <div className="col-span-full">
          {error && (
            <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
          )}
        </div>
      </div>

      <div className="mt-8 flex gap-4">
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

export default CreateAppForm;
