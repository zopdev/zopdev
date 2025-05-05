'use client';

import React, { useContext } from 'react';
import HeadingComponent from '../../../../../../../components/HeaderComponents';
import CustomLinearProgress from '../../../../../../../components/Loaders/LinearLoader';
import Table from '../../../../../../../components/Table/table';
import { calculateAge } from '../../../../../../../utlis/helperFunc';
import ErrorComponent from '../../../../../../../components/ErrorComponent';
import BreadCrumbComp from '../../../../../../../components/BreadCrumb';
import { useParams } from 'next/navigation';
import { AppContext } from '../../../../../../../libs/context';
import usePodList from '../../../../../../../hooks/pods/getPodsList';
import PodsDetails from '../../../../../../../partials/pods/podsDetails';

const headers = [
  { key: 'name', label: 'Name', align: 'left', colClassName: 'sm:!min-w-[300px]' },
  { key: 'container', label: 'Container', align: 'left', colClassName: 'sm:!min-w-[200px]' },
  { key: 'restart', label: 'Restart', align: 'left', colClassName: 'sm:!min-w-[200px]' },
  {
    key: 'controlled',
    label: 'Controlled By',
    align: 'left',
    colClassName: 'sm:!min-w-[200px] whitespace-nowrap',
  },
  //   { key: 'Qos', label: 'Qos', align: 'left', colClassName: 'sm:!min-w-[200px]' },
  { key: 'age', label: 'Age', align: 'left', colClassName: ' sm:!min-w-[200px] whitespace-nowrap' },
  { key: 'state', label: 'Status', align: 'end', colClassName: 'sm:!min-w-[200px]' },
];

const StatusBoxes = ({ data }) => {
  return (
    <div className="flex space-x-2">
      {data.map((container, index) => (
        <div
          key={index}
          className={`w-4 h-4 ${container.ready ? 'bg-green-500' : 'bg-orange-500'} rounded-full group-hover`}
          title={container.name} // Optional: shows container name on hover
        ></div>
      ))}
    </div>
  );
};

const Services = () => {
  const { pods, loading, error } = usePodList();
  const { appData } = useContext(AppContext);
  const params = useParams();

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const bodyData = pods?.data?.pods?.map((item) => {
    return {
      name: item.metadata.name,
      container: <StatusBoxes data={item?.status?.containerStatuses} />,
      restart: item?.status?.containerStatuses?.reduce((acc, curr) => acc + curr.restartCount, 0),
      controlled: item?.metadata?.ownerReferences?.[0]?.kind,
      age: calculateAge(item.metadata.creationTimestamp),
      state:
        item.status.phase === 'Running' || item.status.phase === 'Succeeded' ? (
          <p
            className={` inline-block text-green-700 bg-green-50 ring-green-600/20 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset`}
          >
            {item.status.phase}
          </p>
        ) : (
          <p
            className={` inline-block text-red-700 bg-red-50 ring-red-600/10 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset`}
          >
            {item.status.phase}
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
      name: pods?.metadata?.environmentName ? pods?.metadata?.environmentName : 'loading...',
      link: `/applications/${params?.['application-id']}/environment/${params?.['environment-id']}/services`,
    },
    {
      name: 'Pods',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent
        title={'Pods'}
        // actions={
        //   <Link
        //     href={'#'}
        //     className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'}

        //   `}
        //   >
        //     {<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
        //     Add Pod
        //   </Link>
        // }
      />
      <CustomLinearProgress isLoading={loading} />
      <div>
        {!error && (
          <Table headers={headers} data={bodyData} action={true} renderComponent={PodsDetails} />
        )}
      </div>
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default Services;
