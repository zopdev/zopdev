'use client';

import React, { useContext } from 'react';
import HeadingComponent from '../../../../../components/HeaderComponents';
import Link from 'next/link';
import { useParams, usePathname, useRouter } from 'next/navigation';
import { PlusCircleIcon, PlusIcon } from '@heroicons/react/20/solid';
import Table from '../../../../../components/Table/table';
import DraggableList from '../../../../../partials/draggableList';
import useGetDeploymentSpace from '../../../../../hooks/deploymentSpace/getDeploymentSpace';
import BreadCrumbComp from '../../../../../components/BreadCrumb';
import CustomLinearProgress from '../../../../../components/Loaders/LinearLoader';
import ErrorComponent from '../../../../../components/ErrorComponent';
import { AppContext } from '../../../../../libs/context';

const headers = [
  { key: 'name', label: 'Environment', align: 'left', colClassName: 'sm:w-1/2' },
  { key: 'deployment_space', label: 'Deployment Space', align: 'left', colClassName: 'sm:w-1/2' },
];

const DeploymentSpace = () => {
  const pathname = usePathname();
  const router = useRouter();
  const params = useParams();
  const { appData } = useContext(AppContext);

  const { value, loading, error } = useGetDeploymentSpace();

  const handleAddEnvConfig = (e, id) => {
    e.stopPropagation();
    router.push(`${pathname}/${id}/configureDeploymentSpace`);
  };

  const handleDeploymentList = (data) => {
    const chips = [];
    while (data?.name) {
      chips.push({ name: data.name, selected: true });
      data = data.next;
    }
    return <DraggableList chips={chips} disableDrag={true} />;
  };

  const body = value?.data?.environments?.map((item) => {
    return {
      id: item.id,
      name: item.name,
      deployment_space: item?.deploymentSpace ? (
        handleDeploymentList(item.deploymentSpace)
      ) : (
        <button
          className='className="inline-flex mt-4 items-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"'
          onClick={(e) => handleAddEnvConfig(e, item.id)}
        >
          <div className={`flex gap-1 items-center`}>
            <PlusIcon aria-hidden="true" className="size-5" />
            <div className={`w-[200px]`}>Configure deployment space</div>
          </div>{' '}
        </button>
      ),
    };
  });

  // const handleEdit = (row) => alert(`Edit ${row.name}`);
  const handleDelete = (row) => alert(`Delete ${row.name}`);

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const handleRowClick = (row) => {
    router.push(`${pathname}/${row.id}/services`);
  };

  const breadcrumbList = [
    { name: 'Applications', link: '/applications' },
    {
      name: appData?.APPLICATION_DATA?.isSuccess ? data?.name : 'loading...',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <div>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent
        title={'Environments'}
        actions={
          //   <Tooltip title={currentApplication?.editTooltip}>
          <Link
            href={`${pathname}/addEnvironment`}
            className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'}
            
          `}
          >
            {<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
            Add Environment
          </Link>
          //   </Tooltip>
        }
      />
      <CustomLinearProgress isLoading={loading} />
      <div className="mt-6">
        <Table
          headers={headers}
          data={body}
          // onEdit={handleEdit}
          onDelete={handleDelete}
          action={false}
          handleRowClick={handleRowClick}
          enableRowClick={true}
        />
      </div>

      {error && <ErrorComponent fullPageError errorText={error || 'Something went wrong !'} />}
    </div>
  );
};

export default DeploymentSpace;
