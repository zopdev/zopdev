'use client';

import Head from 'next/head';
import React, { useContext } from 'react';
import AddEnvironment from '../../../../../../partials/Environment/addForm';
import { AppContext } from '../../../../../../libs/context';
import BreadCrumbComp from '../../../../../../components/BreadCrumb';
import { useParams } from 'next/navigation';

const CreateEnvironment = () => {
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
      name: 'Add Environment',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <Head>
        <title>Add Environment</title>
      </Head>

      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <div className="flex items-center justify-center flex-col">
        <div className="divide-y divide-white/5  ">
          <div className="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 pt-7 sm:px-6 md:grid-cols-12 lg:px-8">
            <div className="md:col-span-5">
              <h2 className="text-base font-semibold leading-7 text-gray-900">{`Add Environment`}</h2>
              <p className="mt-1 text-sm leading-6 text-gray-600">
                Create a new environment by providing a unique name. This allows you to define a
                workspace tailored to your needs. You can customize it further after setup.
              </p>
            </div>
            <div className="md:col-span-7">
              <AddEnvironment />
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default CreateEnvironment;
