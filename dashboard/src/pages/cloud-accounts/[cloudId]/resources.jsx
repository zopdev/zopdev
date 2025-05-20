import PageHeading from '@/components/atom/PageHeading';
import BreadCrumb from '@/components/molecules/BreadCrumb';
import Table from '@/components/molecules/Table';

const headers = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'schedule', label: 'Schedule', align: 'left', width: '250px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const data = [
  {
    id: 1,
    name: 'name',
    schedule: 'schedule',
    state: 'Running',
    instance_type: 'instance_type',
    region: 'region',
  },

  {
    id: 1,
    name: 'name',
    schedule: 'schedule',
    state: 'Running',
    instance_type: 'instance_type',
    region: 'region',
  },

  {
    id: 1,
    name: 'name',
    schedule: 'schedule',
    state: 'Running',
    instance_type: 'instance_type',
    region: 'region',
  },
];

const CloudResourcesPage = () => {
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
    <div>
      <BreadCrumb breadcrumbList={breadcrumbList} />
      <PageHeading title={'Resources'} />
      <Table
        headers={headers}
        //data={data}
        data={[]}
        handleRowClick={handleRowClick}
        enableRowClick={false}
        stickyHeader={true}
        emptyStateTitle="No Resources Found"
        // emptyStateDescription="Looks like your cloud account has no active resources right now"
      />
    </div>
  );
};

export default CloudResourcesPage;
