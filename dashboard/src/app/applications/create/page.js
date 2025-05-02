import Head from 'next/head';
import React from 'react';
import CreateAppForm from '../../../../partials/application/createForm';

const CreateApplication = () => {
  return (
    <>
      <Head>
        <title>Create Application</title>
      </Head>
      <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
        <div className="flex items-center justify-center flex-col">
          <div className="divide-y divide-white/5  ">
            <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 pt-7 sm:px-6 md:grid-cols-12 lg:px-8">
              <div className="md:col-span-5">
                <h2 className="text-base font-semibold leading-7 text-gray-900">{`Add Application`}</h2>
                <p className="mt-1 text-sm leading-6 text-gray-600">
                  Add application details by giving it a unique name and configuring environments as
                  needed.
                </p>
              </div>
              <div className="md:col-span-7">
                <CreateAppForm />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default CreateApplication;
