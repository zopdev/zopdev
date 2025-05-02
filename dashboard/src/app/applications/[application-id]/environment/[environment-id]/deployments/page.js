'use client';

import React, { useContext } from 'react';
import HeadingComponent from '../../../../../../../components/HeaderComponents';
import CustomLinearProgress from '../../../../../../../components/Loaders/LinearLoader';
import Table from '../../../../../../../components/Table/table';
import useDeploymentList from '../../../../../../../hooks/deployment/getDeployment';
import ErrorComponent from '../../../../../../../components/ErrorComponent';
import { calculateAge } from '../../../../../../../utlis/helperFunc';
import DeploymentDetails from '../../../../../../../partials/deployments/deploymentDetails';
import BreadCrumbComp from '../../../../../../../components/BreadCrumb';
import { useParams } from 'next/navigation';
import { AppContext } from '../../../../../../../libs/context';

const headers = [
  { key: 'name', label: 'Name', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'pods', label: 'Pods', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'replica', label: 'Replicas', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'age', label: 'Age', align: 'left', colClassName: 'sm:!min-w-[175px]' },
  { key: 'condition', label: 'Conditions', align: 'left', colClassName: 'sm:!min-w-[175px]' },
];

const Deployments = () => {
  const { deployment, loading, error } = useDeploymentList();
  const { appData } = useContext(AppContext);
  const params = useParams();

  const data = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const bodyData = deployment?.data?.deployments?.map((item) => {
    return {
      name: item.metadata.name,
      pods: `${item.status.readyReplicas}/${item.spec.replicas}`,
      replica: item.spec.replicas,
      age: calculateAge(item.metadata.creationTimestamp),
      condition: (
        <p className=" whitespace-nowrap">
          {item.status.conditions
            .map((condition) => condition.type)
            .reverse()
            .join(' ')}
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
      name: deployment?.metadata?.environmentName
        ? deployment?.metadata?.environmentName
        : 'loading...',
      link: `/applications/${params?.['application-id']}/environment/${params?.['environment-id']}/services`,
    },
    {
      name: 'Deployments',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent
        title={'Deployments'}
        // actions={
        //   <Link
        //     // href={`${pathname}/addEnvironment`}
        //     href={'#'}
        //     className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'}

        //   `}
        //   >
        //     {<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
        //     Add Deployment
        //   </Link>
        // }
      />
      <CustomLinearProgress isLoading={loading} />
      <div>
        {!error && (
          <Table
            headers={headers}
            data={bodyData}
            action={true}
            renderComponent={DeploymentDetails}
          />
        )}
      </div>

      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default Deployments;
