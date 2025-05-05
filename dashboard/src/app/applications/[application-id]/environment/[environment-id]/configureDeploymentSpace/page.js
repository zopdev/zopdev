'use client';

import Head from 'next/head';
import React, { useContext } from 'react';
import StepperComponent from '../../../../../../../partials/nestedList';
import { AppContext } from '../../../../../../../libs/context';
import { useParams } from 'next/navigation';
import BreadCrumbComp from '../../../../../../../components/BreadCrumb';
// import AddEnvironment from '../../../../../../partials/Environment/addForm';

const AddDeploymentSpace = () => {
  const { appData } = useContext(AppContext);
  const params = useParams();

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const breadcrumbList = [
    { name: 'Applications', link: '/applications' },
    {
      name: appData?.APPLICATION_DATA?.isSuccess ? data?.name : 'loading...',
      link: `/applications/${params?.['application-id']}/environment`,
    },
    {
      name: 'Configure Deployment Space',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <Head>
        <title>Configure Deployment Space</title>
      </Head>

      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <div className="flex items-center justify-center flex-col">
        <div className="divide-y divide-white/5  ">
          <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 pt-7 sm:px-6 md:grid-cols-12 lg:px-8">
            <div className="md:col-span-5">
              <h2 className="text-base font-semibold leading-7 text-gray-900">{`Configure Deployment Space`}</h2>
              <p className="mt-1 text-sm leading-6 text-gray-600">
                Set up a deployment space for an environment to be used for deployments. Customize
                it further after setup.
              </p>
            </div>
            <div className="md:col-span-7">
              <StepperComponent />
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default AddDeploymentSpace;
