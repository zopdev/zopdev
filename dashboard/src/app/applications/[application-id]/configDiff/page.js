'use client';

import React, { useContext, useEffect, useState } from 'react';
import HeadingComponent from '../../../../../components/HeaderComponents';
import JsonComparisonTableContainer from '../../../../../components/JsonTable';
import useGetDeploymentSpace from '../../../../../hooks/deploymentSpace/getDeploymentSpace';
import { formatTime } from '../../../../../utlis/helperFunc';
import DraggableList from '../../../../../partials/draggableList';
import CustomLinearProgress from '../../../../../components/Loaders/LinearLoader';
import ErrorComponent from '../../../../../components/ErrorComponent';
import { useParams } from 'next/navigation';
import BreadCrumbComp from '../../../../../components/BreadCrumb';
import { AppContext } from '../../../../../libs/context';

const handleDeploymentList = (data) => {
  const chips = [];
  while (data?.name) {
    chips.push({ name: data.name, selected: true });
    data = data.next;
  }
  return <DraggableList chips={chips} disableDrag={true} />;
};

function transformEnvironments(inputArray) {
  return inputArray?.map((item) => ({
    id: item.id,
    label: item.name,
    data: {
      Name: item.name,
      'Deployment Space': () => handleDeploymentList(item.deploymentSpace),
      'Created at': formatTime(item.createdAt),
      'Updated at': formatTime(item.updatedAt),
    },
  }));
}

const ConfigDiff = () => {
  const params = useParams();
  const { appData } = useContext(AppContext);
  const { value, loading, error } = useGetDeploymentSpace();
  const [data, setData] = useState([]);

  useEffect(() => {
    if (value?.data?.environments) {
      const data = transformEnvironments(value?.data?.environments);
      setData(data);
    }
  }, [value]);

  const appName = appData?.APPLICATION_DATA?.data?.filter(
    (item) => item.id === Number(params?.['application-id']),
  )?.[0];

  const breadcrumbList = [
    { name: 'Applications', link: '/applications' },
    {
      name: appData?.APPLICATION_DATA?.isSuccess ? appName?.name : 'loading...',
      link: `/applications/${params?.['application-id']}/environment`,
    },
    {
      name: 'Config Diff',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <div>
      <BreadCrumbComp breadcrumbList={breadcrumbList} />
      <HeadingComponent title={'Configuration Diff'} />
      <CustomLinearProgress isLoading={loading} />
      <JsonComparisonTableContainer data={data} isKeysVisible={true} difference={false} />
      {error && <ErrorComponent fullPageError errorText={error || 'Something went wrong !'} />}
    </div>
  );
};

export default ConfigDiff;
