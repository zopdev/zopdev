'use client';

import React, { useState } from 'react';
import CloudForm from '../../../../partials/cloudAccount/createForm';

const PROVIDER_OPTIONS = {
  0: 'gcp',
  // 1: "aws",
  // 2: "azure",
};

const PROVIDER_TO_TAB_OPTIONS = {
  gcp: 0,
  // aws: 1,
  // azure: 2,
};

// const FORM_PROVIDER_OPTIONS = ["GCP", "AWS", "AZURE"];
const FORM_PROVIDER_OPTIONS = ['GCP'];

const CreateCloud = () => {
  const [value, setValue] = useState(PROVIDER_TO_TAB_OPTIONS.gcp);

  const options = FORM_PROVIDER_OPTIONS.map((item) => ({
    label: item,
    // disabled: isEdit && item.toLowerCase() !== initialValues?.cloudPlatform,
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
              <h2 className="text-base font-semibold leading-7 text-gray-900">Add Cloud Account</h2>
              <p className="mt-1 text-sm leading-6 text-gray-600">
                To add a cloud account, provide your Service Account Key JSON and assign a
                meaningful name to identify your cloud account easily.
              </p>
            </div>
            <div className="md:col-span-7">
              <CloudForm
                provider={PROVIDER_OPTIONS[value]}
                handleTabChange={handleChange}
                options={options}
                tabValue={value}
                value={value}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateCloud;
