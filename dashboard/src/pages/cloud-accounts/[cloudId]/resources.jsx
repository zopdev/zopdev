import ErrorComponent from '@/components/atom/ErrorComponent';
import LinearProgress from '@/components/atom/Loaders/LinerProgress';
import PageHeading from '@/components/atom/PageHeading';
import SwitchButton from '@/components/atom/Switch';
import BreadCrumb from '@/components/molecules/BreadCrumb';
import Table from '@/components/molecules/Table';
import { useGetCloudResources } from '@/queries/cloud-resources';
import { useParams } from 'react-router-dom';

const headers = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'schedule', label: 'Schedule', align: 'left', width: '250px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const CloudResourceRow = (resource) => {
  const handleToggle = (state) => {
    console.log(state);
  };

  return {
    id: resource?.instance_name,
    name: resource?.instance_name,
    schedule: 'schedule',
    state: (
      <div className="min-w-36">
        <SwitchButton
          labelPosition="right"
          isEnabled={false}
          disabled={false}
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
  const cloudResources = useGetCloudResources(cloudId);

  const data = cloudResources?.data?.data || [];

  const handleRowClick = (row) => {
    console.log('Row clicked:', row);
  };

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
        <Table
          headers={headers}
          data={data}
          handleRowClick={handleRowClick}
          enableRowClick={false}
          renderRow={CloudResourceRow}
          emptyStateTitle="No Resources Found"
          // emptyStateDescription="Looks like your cloud account has no active resources right now"
        />
      )}
      {cloudResources?.isError && (
        <ErrorComponent errorText={cloudResources?.error?.message || 'Something went wrong'} />
      )}
    </>
  );
};

export default CloudResourcesPage;
