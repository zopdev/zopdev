import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { PlusCircleIcon } from '@heroicons/react/20/solid';

import BlankCloudAccountSvg from '@/assets/svg/BlankCloudAccount';
import Button from '@/components/atom/Button';
import EmptyComponent from '@/components/atom/EmptyComponent';
import ErrorComponent from '@/components/atom/ErrorComponent';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import LinearProgress from '@/components/atom/Loaders/LinerProgress';
import PageHeading from '@/components/atom/PageHeading';
import ResourceGroupManager from '@/components/container/resources/AddResourceGroup';
import ResourceGroupAccordion from '@/components/container/resources/ResourceGroupAccordion';
import { CloudResourceRow } from '@/components/container/resources/ResourceTableRow';
import BreadCrumb from '@/components/molecules/BreadCrumb';
import Table from '@/components/molecules/Table';
import { Tabs } from '@/components/molecules/Tabs';
import { useGetCloudResources, useGetResourceGroup } from '@/queries/cloud-resources';

const headers = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const pageTabs = [{ label: 'Resources' }, { label: 'Resource Group' }];

const CreateResourceGroupButton = ({ resources }) => (
  <FullScreenOverlay
    customCTA={
      <Button startEndornment={<PlusCircleIcon className="h-5 w-5" />}>
        Create Resource Group
      </Button>
    }
    title="Create Resource Group"
    size="xl"
    variant="drawer"
    renderContent={ResourceGroupManager}
    renderContentProps={{ resources }}
  />
);

const CloudResourcesPage = () => {
  const { cloudId } = useParams();
  const [activeTab, setActiveTab] = useState('Resources');

  const { data: resourceData = [], isLoading, isError, error } = useGetCloudResources(cloudId);

  const { data: resourceGroupData = [] } = useGetResourceGroup(cloudId);

  const breadcrumbList = [
    { name: 'Cloud Accounts', link: '/cloud-accounts' },
    { name: 'Resources', link: '#', disable: true },
  ];

  return (
    <>
      <BreadCrumb breadcrumbList={breadcrumbList} />
      <PageHeading title="Resources" />
      <LinearProgress isLoading={isLoading} />

      {!isLoading && !isError && (
        <>
          <div className="flex justify-between flex-wrap xs:space-y-2 md:space-y-0">
            <Tabs tabs={pageTabs} activeTab={activeTab} onTabChange={setActiveTab} size="md" />
            <div className="flex items-center">
              <CreateResourceGroupButton resources={resourceData} />
            </div>
          </div>

          {activeTab === 'Resources' ? (
            <Table
              headers={headers}
              data={resourceData}
              enableRowClick={false}
              renderRow={CloudResourceRow}
              emptyStateTitle="No Resources Found"
              // emptyStateDescription="Looks like your cloud account has no active resources right now"
            />
          ) : (
            <>
              <ResourceGroupAccordion groups={resourceGroupData} defaultExpandedIds={[]} />
              {resourceGroupData.length === 0 && (
                <EmptyComponent
                  imageComponent={<BlankCloudAccountSvg />}
                  customButton={<CreateResourceGroupButton resources={resourceData} />}
                  title="Please start by setting up your first resource group"
                />
              )}
            </>
          )}
        </>
      )}

      {isError && <ErrorComponent errorText={error?.message || 'Something went wrong'} />}
    </>
  );
};

export default CloudResourcesPage;
