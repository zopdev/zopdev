'use client';

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import Label from '../../components/Label';
import SimpleLoader from '../../components/Loaders/SimpleLoader';
import AnimateHeight from '../../components/Animation/animateHeight';
import useAddDeploymentConfig from '../../hooks/deploymentSpace/addDeploymentSpace';
import { useParams } from 'next/navigation';
import Button from '../../components/Button';
import ErrorComponent from '../../components/ErrorComponent';
import { transformDeploymentConfigData } from './deploymentConfig';

const StepperComponent = () => {
  const params = useParams();
  const { handleSubmit, isLoading, error } = useAddDeploymentConfig();

  const [step, setStep] = useState(1);
  const [cloudAccounts, setCloudAccounts] = useState([]);
  const [selectedCloudAccount, setSelectedCloudAccount] = useState(null);
  const [deploymentOptions, setDeploymentOptions] = useState([]);
  const [selectedOption, setSelectedOption] = useState(null);
  const [dropdowns, setDropdowns] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showSubmit, setShowSubmit] = useState(false);

  // Fetch cloud accounts initially
  useEffect(() => {
    const fetchCloudAccounts = async () => {
      setLoading(true);
      try {
        const { data } = await axios.get('http://localhost:8000/cloud-accounts');
        setCloudAccounts(data?.data || []);
      } catch (error) {
        console.error('Error fetching cloud accounts:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchCloudAccounts();
  }, []);

  // Handle cloud account selection
  const handleCloudAccountChange = async (selectedData) => {
    setSelectedCloudAccount(selectedData);
    resetSteps(2);
    setLoading(true);

    try {
      const { data } = await axios.get(
        `http://localhost:8000/cloud-accounts/${selectedData?.id}/deployment-space/options`,
      );
      setDeploymentOptions(data?.data || []);
    } catch (error) {
      console.error('Error fetching deployment options:', error);
    } finally {
      setLoading(false);
    }
  };

  // Handle deployment option selection
  const handleDeploymentOptionChange = async (option) => {
    setSelectedOption(option);
    resetSteps(3);

    setLoading(true);

    try {
      const { data } = await axios.get(`http://localhost:8000${option.path}`);
      if (data?.data?.next) {
        addDropdown(data.data.options || [], data.data.next, data.data.metadata);
      } else {
        alert('All selections complete.');
      }
    } catch (error) {
      console.error('Error fetching additional options:', error);
    } finally {
      setLoading(false);
    }
  };

  // Handle dynamic dropdown changes
  const handleDropdownChange = async (index, selectedOption) => {
    const updatedDropdowns = [...dropdowns];
    updatedDropdowns[index].selected = selectedOption;
    setDropdowns(updatedDropdowns.slice(0, index + 1));

    const next = updatedDropdowns[index].next;
    if (!next) {
      setShowSubmit(true);
      return;
    }

    const queryParams = buildQueryParams(next.params, selectedOption);
    const apiUrl = `http://localhost:8000${next.path}?${queryParams}`;

    setLoading(true);

    try {
      const { data } = await axios.get(apiUrl);
      if (data?.data?.options) {
        addDropdown(data.data.options || [], data.data.next, data.data.metadata);
      } else {
        alert('All selections complete.');
      }
    } catch (error) {
      console.error('Error fetching next options:', error);
    } finally {
      setLoading(false);
    }
  };

  // Utility: Reset steps and dropdowns
  const resetSteps = (newStep) => {
    setStep(newStep);
    // setDeploymentOptions([]);
    setDropdowns([]);
  };

  // Utility: Add a new dropdown
  const addDropdown = (options, next, metadata) => {
    setDropdowns((prev) => [...prev, { options, selected: null, next, metadata }]);
  };

  // Utility: Build query parameters for API requests
  const buildQueryParams = (params, option) => {
    const queryParams = new URLSearchParams();
    for (const key in params) {
      if (option[params[key]] !== undefined) {
        queryParams.append(key, option[params[key]]);
      }
    }
    return queryParams.toString();
  };

  // Handle submit button click
  const handleFormSubmit = () => {
    const selectedData = {
      cloudAccount: selectedCloudAccount,
      deploymentOption: selectedOption,
      dropdownSelections: dropdowns.map((dropdown) => dropdown.selected),
    };

    const updatedData = transformDeploymentConfigData(selectedData);

    handleSubmit(params?.['environment-id'], updatedData);
  };

  return (
    <AnimateHeight>
      <div className="p-6 space-y-8 bg-gray-50">
        {/* Step 1: Cloud Accounts */}
        {step >= 1 && (
          <div>
            <Label htmlFor="name">Select Cloud Account</Label>
            <select
              className="w-full p-2 border rounded-md"
              value={selectedCloudAccount?.id || ''}
              onChange={(e) => {
                const selectedId = e.target.value;
                const selectedAccount = cloudAccounts.find(
                  (account) => account.id.toString() === selectedId,
                );
                handleCloudAccountChange(selectedAccount); // Pass the full selected object
              }}
            >
              <option value="" disabled>
                Select a cloud account
              </option>
              {cloudAccounts.map((account) => (
                <option key={account.id} value={account.id}>
                  {account.name}
                </option>
              ))}
            </select>
          </div>
        )}
        {/* Step 2: Deployment Options */}
        {step >= 2 && (
          <div>
            <Label htmlFor="name">Select Deployment Space </Label>
            <select
              className="w-full p-2 border rounded-md"
              value={selectedOption?.name || ''}
              onChange={(e) => {
                handleDeploymentOptionChange(deploymentOptions[e.target.selectedIndex - 1]);
              }}
            >
              <option value="" disabled>
                Select an option
              </option>
              {deploymentOptions.map((option, index) => (
                <option key={index} value={option.name}>
                  {option.name}
                </option>
              ))}
            </select>
          </div>
        )}
        {/* Step 3+: Dynamic Dropdowns */}
        {dropdowns.map((dropdown, index) => (
          <div key={index}>
            <Label htmlFor="name">Select {dropdown?.metadata?.name}</Label>
            <select
              className="w-full p-2 border rounded-md"
              value={dropdown?.selected?.name || ''}
              onChange={(e) =>
                handleDropdownChange(index, dropdown.options[e.target.selectedIndex - 1])
              }
            >
              <option value="" disabled>
                Select an option
              </option>
              {dropdown.options.map((option, i) => (
                <option key={i} value={option.name}>
                  {option.name}
                </option>
              ))}
            </select>
          </div>
        ))}
        {/* Submit Button */}
        {loading && <SimpleLoader />}
        <div className="text-sm text-red-500 mt-5 grid grid-cols-1 gap-x-6 gap-y-0 sm:max-w-xl sm:grid-cols-6">
          <div className="col-span-full">
            {error && (
              <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
            )}
          </div>
        </div>
        {showSubmit && (
          <Button onClick={handleFormSubmit} disabled={isLoading}>
            Add
            {isLoading && (
              <div
                style={{ borderTopColor: 'transparent' }}
                className="w-4 h-4 border-4 border-blue-200 rounded-full animate-spin"
              />
            )}
          </Button>
        )}
      </div>
    </AnimateHeight>
  );
};

export default StepperComponent;
