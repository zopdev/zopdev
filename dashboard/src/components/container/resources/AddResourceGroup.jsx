import Button from '@/components/atom/Button';
import Checkbox from '@/components/atom/Checkbox';
import ErrorComponent from '@/components/atom/ErrorComponent';
import Input from '@/components/atom/Input';
import Label from '@/components/atom/Label';
import Textarea from '@/components/atom/Textarea';
import { usePostResourceGroup } from '@/queries/cloud-resources';
import React, { useEffect, useState } from 'react';

const ResourceGroupManager = ({ resources, onClose }) => {
  const [resourceGroupName, setResourceGroupName] = useState('');
  const [description, setDescription] = useState('');

  const postToResourceGroup = usePostResourceGroup();
  const [selectedResources, setSelectedResources] = useState([]);

  const handleResourceSelect = (resourceId) => {
    setSelectedResources((prev) => {
      if (prev.includes(resourceId)) {
        return prev.filter((id) => id !== resourceId);
      } else {
        return [...prev, resourceId];
      }
    });
  };

  const handleSelectAll = () => {
    if (selectedResources.length === resources.length) {
      setSelectedResources([]);
    } else {
      setSelectedResources(resources.map((resource) => resource.id));
    }
  };

  const isAllSelected = selectedResources.length === resources.length;
  const isIndeterminate =
    selectedResources.length > 0 && selectedResources.length < resources.length;

  const handleSave = () => {
    const selectedResourceDetails = resources.filter((resource) =>
      selectedResources.includes(resource.id),
    );

    postToResourceGroup.mutate({
      name: resourceGroupName,
      description,
      cloudAccId: 1,
      resource_ids: selectedResourceDetails.map((item) => item?.id),
    });
  };

  const isFormValid = resourceGroupName.trim() !== '' && selectedResources.length > 0;

  useEffect(() => {
    if (postToResourceGroup?.isSuccess) {
      onClose();
      setSelectedResources([]);
      setResourceGroupName('');
      setDescription('');
    }
  }, [postToResourceGroup?.isPending]);

  return (
    <div>
      <div className="mb-4">
        <Label htmlFor="resourceGroupName">Name</Label>
        <Input
          id="resourceGroupName"
          value={resourceGroupName}
          onChange={(e) => setResourceGroupName(e.target.value)}
          placeholder="Enter name"
        />
      </div>

      <div className="mb-6">
        <Label htmlFor="description"> Description</Label>
        <Textarea
          rows={3}
          id="description"
          name="description"
          placeholder="Enter description (optional)"
          variant="outlined"
          className="focus:outline-none focus:ring-1 focus:ring-primary-500"
          onChange={(e) => setDescription(e.target.value)}
          value={description}
        />
      </div>

      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-800 mb-4">Select Resources</h2>

        <div className="mb-3 pb-3 border-b border-gray-200">
          <Label className="flex items-center cursor-pointer">
            <Checkbox
              checked={isAllSelected}
              ref={(input) => {
                if (input) input.indeterminate = isIndeterminate;
              }}
              onChange={handleSelectAll}
            />
            <span className="ml-2 text-sm font-medium text-gray-700">
              Select All ({selectedResources.length} of {resources.length} selected)
            </span>
          </Label>
        </div>

        <div className="space-y-2 max-h-64 overflow-y-auto">
          {resources.map((resource) => (
            <div key={resource.id} className="flex items-center p-2 hover:bg-gray-50 rounded">
              <Label className="flex items-center cursor-pointer w-full">
                <Checkbox
                  checked={selectedResources.includes(resource.id)}
                  onChange={() => handleResourceSelect(resource.id)}
                />
                <div className="ml-3 flex-1">
                  <div className="text-sm font-medium text-gray-900">{resource.name}</div>
                  <div className="text-xs text-gray-500">{resource.type}</div>
                </div>
              </Label>
            </div>
          ))}
        </div>
      </div>

      {!isFormValid && (resourceGroupName.trim() === '' || selectedResources.length === 0) && (
        <div className="mt-3 text-sm text-red-600">
          {selectedResources.length === 0 && 'Please select at least one resource.'}
        </div>
      )}
      {postToResourceGroup?.isError && (
        <ErrorComponent errorText={postToResourceGroup?.error?.message} className="mb-4" />
      )}

      <div className="flex justify-end">
        <Button
          onClick={handleSave}
          disabled={!isFormValid}
          loading={postToResourceGroup?.isPending}
        >
          Save Resource Group
        </Button>
      </div>
    </div>
  );
};

export default ResourceGroupManager;
