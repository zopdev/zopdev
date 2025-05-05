'use client';

import React from 'react';
import HeadingComponent from '../../../components/HeaderComponents';
import { PlusCircleIcon } from '@heroicons/react/20/solid';
import EmptyComponent from '../../../components/EmptyPageComponent';
import BlankCloudAccountSvg from '../../../svg/emptyCloudAccount';
import Link from 'next/link';
import CompleteLoader from '../../../components/Loaders/CompletePageLoader';
import ErrorComponent from '../../../components/ErrorComponent';
import useApplicationList from '../../../hooks/application/getApplicationList';
import ApplicationCard from '../../../components/Cards/applicationCard';

const CloudAccounts = () => {
  const { applications, loading, error } = useApplicationList();

  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
      {applications?.data?.length > 0 && (
        <HeadingComponent
          title={'Applications'}
          actions={
            <Link href={'/applications/create'}>
              <button
                type="button"
                className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600`}
              >
                <PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />
                Add Application
              </button>
            </Link>
          }
        />
      )}

      {loading && <CompleteLoader />}
      <div className="flex gap-4 w-full justify-start mt-4 flex-wrap">
        {applications?.data?.map((item, idx) => {
          return (
            <ApplicationCard
              key={idx}
              data={item}
              //   view={'cloudAccount'}
              //   handleLogsOpen={handleLogsOpen}
              //   handleRetryStatus={handleRetryStatus}
              //   putToProvider={putToProvider}
            />
          );
        })}
      </div>

      {applications?.data?.length === 0 && (
        <EmptyComponent
          imageComponent={<BlankCloudAccountSvg />}
          redirectLink={'/applications/create'}
          buttonTitle={'Add Application'}
          title={'Please start by setting up your first application'}
        />
      )}

      {error && <ErrorComponent errorText={error || 'Something went wrong'} />}
    </div>
  );
};

export default CloudAccounts;
