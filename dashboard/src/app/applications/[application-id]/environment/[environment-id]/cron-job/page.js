'use client';

import React, { useContext } from 'react';
import HeadingComponent from '../../../../../../../components/HeaderComponents';
import CustomLinearProgress from '../../../../../../../components/Loaders/LinearLoader';
import Table from '../../../../../../../components/Table/table';
import ErrorComponent from '../../../../../../../components/ErrorComponent';
import { calculateAge } from '../../../../../../../utlis/helperFunc';
import BreadCrumbComp from '../../../../../../../components/BreadCrumb';
import { useParams } from 'next/navigation';
import { AppContext } from '../../../../../../../libs/context';
import useCronJobList from '../../../../../../../hooks/cronJob/getCronJobList';
import CronJobDetails from '../../../../../../../partials/cronJob/cronJobDetails';

const headers = [
  { key: 'name', label: 'Name', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'schedule', label: 'Schedule', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'suspend', label: 'Suspend', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'active', label: 'Active', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  {
    key: 'last_schedule',
    label: 'Last Schedule',
    align: 'left',
    colClassName: 'sm:!min-w-[175px]',
  },
  { key: 'age', label: 'Age', align: 'left', colClassName: 'sm:!min-w-[175px]' },
];

const CronJobs = () => {
  const { cronJob, loading, error } = useCronJobList();
  const { appData } = useContext(AppContext);
  const params = useParams();

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const bodyData = cronJob?.data?.cronjobs?.map((item) => {
    return {
      name: item.metadata.name,
      schedule: item.spec.schedule,
      suspend: item.spec.suspend.toString(),
      active: 0,
      last_schedule: calculateAge(item.status.active.lastScheduleTime),
      age: calculateAge(item.metadata.creationTimestamp),
    };
  });

  const breadcrumbList = [
    { name: 'Applications', link: '/applications' },
    {
      name: appData?.APPLICATION_DATA?.isSuccess ? data?.name : 'loading...',
      link: `/applications/${params?.['application-id']}/environment`,
    },
    {
      name: cronJob?.metadata?.environmentName ? cronJob?.metadata?.environmentName : 'loading...',
      link: `/applications/${params?.['application-id']}/environment/${params?.['environment-id']}/services`,
    },
    {
      name: 'Cron Jobs',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent
        title={'Cron Jobs'}
        // actions={
        //   <Link
        //     // href={`${pathname}/addEnvironment`}
        //     href={'#'}
        //     className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'}

        //   `}
        //   >
        //     {<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
        //     Add Cron Job
        //   </Link>
        // }
      />
      <CustomLinearProgress isLoading={loading} />
      <div>
        {!error && (
          <Table headers={headers} data={bodyData} action={true} renderComponent={CronJobDetails} />
        )}
      </div>

      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default CronJobs;
