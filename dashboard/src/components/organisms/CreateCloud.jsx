import React, { useState } from 'react';
import CloudForm from '@/components/organisms/CloudForm.jsx';

const CreateCloud = ({ audit, setIsComplete, updateData, data }) => {
  const PROVIDERS = [
    { key: 'gcp', label: 'GCP' },
    { key: 'aws', label: 'AWS' },
    { key: 'azure', label: 'Azure' },
  ];

  const [activeProvider, setActiveProvider] = useState(PROVIDERS[0]);

  const handleProviderChange = (provider) => {
    setActiveProvider(provider);
  };

  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
      <div className="flex items-center justify-center flex-col">
        <div className="divide-y divide-white/5">
          <div className="grid max-w-7xl grid-cols-1 gap-x-3 gap-y-10 md:px-4 md:pt-7 sm:px-6 md:grid-cols-12 lg:px-8">
            <div className="md:col-span-5">
              <h2 className="text-base font-semibold leading-7 text-secondary-900">
                Add Cloud Account
              </h2>
              <p className="mt-1 text-md md:pr-10 leading-6 text-secondary-600">
                Provide the necessary credentials for your cloud account.
              </p>
            </div>
            <div className="md:col-span-7">
              <CloudForm
                key={activeProvider.key}
                audit={audit}
                provider={activeProvider.key}
                providers={PROVIDERS}
                activeProvider={activeProvider}
                handleProviderChange={handleProviderChange}
                setIsComplete={setIsComplete}
                updateData={updateData}
                data={data}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateCloud;
