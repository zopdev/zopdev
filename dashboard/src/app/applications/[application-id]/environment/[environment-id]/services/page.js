'use client';

import React, { useContext } from 'react';
import HeadingComponent from '../../../../../../../components/HeaderComponents';
import CustomLinearProgress from '../../../../../../../components/Loaders/LinearLoader';
import Table from '../../../../../../../components/Table/table';
import { formatTime } from '../../../../../../../utlis/helperFunc';
import useServiceList from '../../../../../../../hooks/service/getServiceList';
import ErrorComponent from '../../../../../../../components/ErrorComponent';
import BreadCrumbComp from '../../../../../../../components/BreadCrumb';
import { useParams } from 'next/navigation';
import { AppContext } from '../../../../../../../libs/context';
import ServiceDetails from '../../../../../../../partials/service/serviceDetail';

const headers = [
  { key: 'name', label: 'Name', align: 'left', colClassName: 'sm:w-1/4' },
  { key: 'loadBalancerIP', label: 'Load Balance IP', align: 'left', colClassName: 'sm:w-1/4' },
  { key: 'timestamps', label: 'Timestamps', align: 'left', colClassName: 'sm:w-1/3' },
  { key: 'status', label: 'Status', align: 'end', colClassName: 'sm:w-1/4' },
];

const Services = () => {
  const { services, loading, error } = useServiceList();
  const { appData } = useContext(AppContext);
  const params = useParams();

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const bodyData = services?.data?.services?.map((item) => {
    return {
      name: item.metadata.name,
      loadBalancerIP: item.spec.loadBalancerIP,
      timestamps: formatTime(item.metadata.creationTimestamp),
      status: item.status.loadBalancer.ingress ? (
        <p
          className={` inline-block text-green-700 bg-green-50 ring-green-600/20 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset`}
        >
          Active
        </p>
      ) : (
        <p
          className={` inline-block text-red-700 bg-red-50 ring-red-600/10 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset`}
        >
          Inactive
        </p>
      ),
    };
  });

  const breadcrumbList = [
    { name: 'Applications', link: '/applications' },
    {
      name: appData?.APPLICATION_DATA?.isSuccess ? data?.name : 'loading...',
      link: `/applications/${params?.['application-id']}/environment`,
    },
    {
      name: 'Services',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent
        title={'Services'}
        // actions={
        //   <Link
        //     // href={`${pathname}/addEnvironment`}
        //     href={'#'}
        //     className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'}

        //   `}
        //   >
        //     {<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
        //     Add Service
        //   </Link>
        // }
      />
      <CustomLinearProgress isLoading={loading} />
      <div>
        {!error && (
          <Table headers={headers} data={bodyData} action={true} renderComponent={ServiceDetails} />
        )}
      </div>
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default Services;
