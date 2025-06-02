import Button from '@/components/atom/Button';
import ErrorComponent from '@/components/atom/ErrorComponent';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import LinearProgress from '@/components/atom/Loaders/LinerProgress';
import PageHeading from '@/components/atom/PageHeading';
import SwitchButton from '@/components/atom/Switch';
import ResourceGroupManager from '@/components/container/resources/AddResourceGroup';
import ResourceGroupAccordion from '@/components/container/resources/ResourceGroupAccordion';
import BreadCrumb from '@/components/molecules/BreadCrumb';
import Table from '@/components/molecules/Table';
import { Tabs } from '@/components/molecules/Tabs';
import { toast } from '@/components/molecules/Toast';
import {
  useGetCloudResources,
  useGetResourceGroup,
  usePostResourceState,
} from '@/queries/cloud-resources';
import { PlusCircleIcon } from '@heroicons/react/20/solid';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

const header = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const pageTabs = [{ label: 'Resources' }, { label: 'Resource Group' }];

const CloudResourceRow = (resource) => {
  const [currentState, setCurrentState] = useState(resource?.status === 'RUNNING' ? true : false);
  const { cloudId } = useParams();
  const resourceStateChanger = usePostResourceState();

  const handleToggle = (state) => {
    setCurrentState(state);
    resourceStateChanger.mutate({
      cloudAccId: parseInt(cloudId),
      id: resource?.id,
      name: resource?.instance_name,
      type: resource?.instance_type,
      state: state ? 'START' : 'SUSPEND',
    });
  };

  useEffect(() => {
    if (resourceStateChanger.isError) {
      toast.failed(resourceStateChanger?.error?.message);
      setCurrentState((prev) => !prev);
    }
  }, [resourceStateChanger.isError, resourceStateChanger.error]);

  return {
    id: resource?.instance_name,
    name: resource?.instance_name,
    state: (
      <div className="min-w-36">
        <SwitchButton
          labelPosition="right"
          value={currentState}
          disabled={resourceStateChanger?.isPending}
          onChange={handleToggle}
          titleList={{ true: 'Running', false: 'Suspended' }}
          name="status"
        />
      </div>
    ),
    instance_type: resource?.instance_type,
    region: resource?.region,
  };
};

const CloudResourcesPage = () => {
  const { cloudId } = useParams();
  const [activeTab, setActiveTab] = useState('Resources');
  const cloudResources = useGetCloudResources(cloudId);
  const resourceGroup = useGetResourceGroup(cloudId);

  const data = cloudResources?.data || [];

  const breadcrumbList = [
    { name: 'Cloud Accounts', link: '/cloud-accounts' },
    {
      name: 'Resources',
      link: `#`,
      disable: true,
    },
  ];

  return (
    <>
      <BreadCrumb breadcrumbList={breadcrumbList} />
      <PageHeading title={'Resources'} />
      <LinearProgress isLoading={cloudResources?.isLoading} />

      {!cloudResources?.isLoading && !cloudResources?.isError && (
        <>
          <div className="flex justify-between flex-wrap xs:space-y-2 md:space-y-0">
            <Tabs tabs={pageTabs} activeTab={activeTab} onTabChange={setActiveTab} size="md" />
            <div className="flex items-center">
              <FullScreenOverlay
                customCTA={
                  <Button
                    startEndornment={<PlusCircleIcon className="h-5 w-5" aria-hidden="true" />}
                  >
                    Create Resource Group
                  </Button>
                }
                title="Create Resource Group"
                size="xl"
                variant="drawer"
                renderContent={ResourceGroupManager}
                renderContentProps={{
                  resources: data,
                }}
              />
            </div>
          </div>

          <Table
            headers={header}
            data={data}
            enableRowClick={false}
            renderRow={CloudResourceRow}
            emptyStateTitle="No Resources Found"
            // emptyStateDescription="Looks like your cloud account has no active resources right now"
          />

          <ResourceGroupAccordion />
        </>
      )}
      {cloudResources?.isError && (
        <ErrorComponent errorText={cloudResources?.error?.message || 'Something went wrong'} />
      )}
    </>
  );
};

export default CloudResourcesPage;
