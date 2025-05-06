'use client';

import React, { useState } from 'react';
import CloudForm from '@/components/organisms/CloudForm.jsx';

const CreateCloud = ({ audit, setIsComplete, updateData, data }) => {
  const PROVIDER_OPTIONS = {
    0: 'gcp',
  };

  const PROVIDER_TO_TAB_OPTIONS = {
    gcp: 0,
  };
  const FORM_PROVIDER_OPTIONS = ['GCP'];

  const [value, setValue] = useState(PROVIDER_TO_TAB_OPTIONS.gcp);

  const options = FORM_PROVIDER_OPTIONS.map((item) => ({
    label: item,
  }));

  const handleChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
      <div className="flex items-center justify-center flex-col">
        <div className="divide-y divide-white/5  ">
          <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 pt-7 sm:px-6 md:grid-cols-12 lg:px-8">
            <div className="md:col-span-5">
              <h2 className="text-base font-semibold leading-7 text-secondary-900">
                Add Cloud Account
              </h2>
              <p className="mt-1 text-md leading-6 text-secondary-600">
                To add a cloud account, provide your Service Account Key JSON and assign a
                meaningful name to identify your cloud account easily.
              </p>
            </div>
            <div className="md:col-span-7">
              <CloudForm
                audit={audit}
                provider={PROVIDER_OPTIONS[value]}
                handleTabChange={handleChange}
                options={options}
                tabValue={value}
                value={value}
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
