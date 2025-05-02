'use client';

import React from 'react';
import HeadingComponent from '../../../components/HeaderComponents';
import { PlusCircleIcon } from '@heroicons/react/20/solid';
import CloudAccountCard from '../../../components/Cards/cloudAccountCard';
import EmptyComponent from '../../../components/EmptyPageComponent';
import BlankCloudAccountSvg from '../../../svg/emptyCloudAccount';
import Link from 'next/link';
import useCloudAccountList from '../../../hooks/cloudAccount/getCloudAccountList';
import CompleteLoader from '../../../components/Loaders/CompletePageLoader';
import ErrorComponent from '../../../components/ErrorComponent';

const CloudAccounts = () => {
  const { cloudAccounts, loading, error } = useCloudAccountList();

  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
      {cloudAccounts?.data?.length > 0 && (
        <HeadingComponent
          title={'Cloud Accounts'}
          actions={
            <Link href={'/cloud-accounts/create'}>
              <button
                type="button"
                className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600`}
              >
                <PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />
                Add Cloud Account
              </button>
            </Link>
          }
        />
      )}

      {loading && <CompleteLoader />}
      <div className="flex gap-4 w-full justify-start mt-4 flex-wrap">
        {cloudAccounts?.data?.map((item, idx) => {
          return <CloudAccountCard key={idx} item={item} />;
        })}
      </div>

      {cloudAccounts?.data?.length === 0 && (
        <EmptyComponent
          imageComponent={<BlankCloudAccountSvg />}
          redirectLink={'/cloud-accounts/create'}
          buttonTitle={'Add Cloud Account'}
          title={'Please start by setting up your first cloud account'}
        />
      )}

      {error && <ErrorComponent errorText={error || 'Something went wrong'} />}
    </div>
  );
};

export default CloudAccounts;
