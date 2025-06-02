import React, { useEffect, useState } from 'react';
import Button from '@/components/atom/Button';
import Checkbox from '@/components/atom/Checkbox';
import ErrorComponent from '@/components/atom/ErrorComponent';
import Input from '@/components/atom/Input';
import Label from '@/components/atom/Label';
import Textarea from '@/components/atom/Textarea';
import { usePostResourceGroup, usePutResourceGroup } from '@/queries/cloud-resources';

const ResourceGroupManager = ({ resources, onClose, initialData = null }) => {
  const [form, setForm] = useState({ name: '', description: '' });
  const [selectedResourceIds, setSelectedResourceIds] = useState([]);

  const {
    mutate: createGroup,
    isPending: isCreating,
    isSuccess: createSuccess,
    isError: isCreateError,
    error: createError,
  } = usePostResourceGroup();

  const {
    mutate: updateGroup,
    isPending: isUpdating,
    isSuccess: updateSuccess,
    isError: isUpdateError,
    error: updateError,
  } = usePutResourceGroup();

  const isEditMode = Boolean(initialData?.id);

  useEffect(() => {
    if (initialData) {
      setForm({ name: initialData.name || '', description: initialData.description || '' });
      setSelectedResourceIds(initialData?.resources?.map((item) => item?.id) || []);
    }
  }, [initialData]);

  useEffect(() => {
    if (createSuccess || updateSuccess) {
      onClose();
      setForm({ name: '', description: '' });
      setSelectedResourceIds([]);
    }
  }, [createSuccess, updateSuccess]);

  const handleInputChange = (e) => {
    const { id, value } = e.target;
    setForm((prev) => ({ ...prev, [id]: value }));
  };

  const handleResourceToggle = (id) => {
    setSelectedResourceIds((prev) =>
      prev.includes(id) ? prev.filter((r) => r !== id) : [...prev, id],
    );
  };

  const handleSelectAll = () => {
    if (selectedResourceIds.length === resources.length) {
      setSelectedResourceIds([]);
    } else {
      setSelectedResourceIds(resources.map((r) => r.id));
    }
  };

  const handleSave = () => {
    const payload = {
      name: form.name,
      description: form.description,
      cloudAccId: 1,
      resource_ids: selectedResourceIds,
    };

    if (isEditMode) {
      updateGroup({ ...payload, id: initialData.id });
    } else {
      createGroup(payload);
    }
  };

  const isFormValid = form.name.trim() !== '' && selectedResourceIds.length > 0;
  const isAllSelected = selectedResourceIds.length === resources.length;
  const isIndeterminate = selectedResourceIds.length > 0 && !isAllSelected;
  const isPending = isCreating || isUpdating;
  const isError = isCreateError || isUpdateError;
  const error = createError || updateError;

  return (
    <div>
      <div className="mb-4">
        <Label htmlFor="name">Name</Label>
        <Input id="name" value={form.name} onChange={handleInputChange} placeholder="Enter name" />
      </div>

      <div className="mb-6">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={form.description}
          onChange={handleInputChange}
          placeholder="Enter description (optional)"
          rows={3}
          variant="outlined"
          className="focus:outline-none focus:ring-1 focus:ring-primary-500"
        />
      </div>

      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-800 mb-4">Select Resources</h2>

        <div className="mb-3 pb-3 border-b border-gray-200">
          <Label className="flex items-center cursor-pointer">
            <Checkbox
              checked={isAllSelected}
              ref={(el) => {
                if (el) el.indeterminate = isIndeterminate;
              }}
              onChange={handleSelectAll}
            />
            <span className="ml-2 text-sm font-medium text-gray-700">
              Select All ({selectedResourceIds.length} of {resources.length} selected)
            </span>
          </Label>
        </div>

        <div className="space-y-2 max-h-64 overflow-y-auto">
          {resources.map(({ id, name, type }) => (
            <div key={id} className="flex items-center p-2 hover:bg-gray-50 rounded">
              <Label className="flex items-center cursor-pointer w-full">
                <Checkbox
                  checked={selectedResourceIds.includes(id)}
                  onChange={() => handleResourceToggle(id)}
                />
                <div className="ml-3 flex-1">
                  <div className="text-sm font-medium text-gray-900">{name}</div>
                  <div className="text-xs text-gray-500">{type}</div>
                </div>
              </Label>
            </div>
          ))}
        </div>
      </div>

      {!isFormValid && (
        <div className="mt-3 text-sm text-red-600">
          {selectedResourceIds.length === 0 && 'Please select at least one resource.'}
        </div>
      )}

      {isError && <ErrorComponent errorText={error?.message} className="mb-4" />}

      <div className="flex justify-end">
        <Button onClick={handleSave} disabled={!isFormValid} loading={isPending}>
          {isEditMode ? 'Update Resource Group' : 'Save Resource Group'}
        </Button>
      </div>
    </div>
  );
};

export default ResourceGroupManager;
