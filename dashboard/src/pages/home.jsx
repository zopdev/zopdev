import { CloudIcon, RocketLaunchIcon } from '@heroicons/react/24/outline';
import DashBoardCard from '@/components/molecule/Cards/DashBoardCard.jsx';
import DashboardSection from '@/components/organisms/DashBoardSection.jsx';
import Button from '@/components/atom/Button/index.jsx';
import { useNavigate } from 'react-router-dom';
import CloudAccountAuditCards from '@/components/molecule/Cards/CloudAuditCard.jsx';
import SimpleLoader from '@/components/atom/Loaders/SimpleLoader.jsx';

const Home = () => {
  const navigate = useNavigate();

  const handleAuditClick = () => {
    navigate('/cloud-setup');
  };

  // const cloudAccounts = [
  //   {
  //     id: 'account1',
  //     title: 'Zop Cloud',
  //     subtitle: '2 Apps',
  //     status: 'READY',
  //     icon: 'cloud',
  //     providerType: 'gcp',
  //     lastUpdatedBy: 'owner@zop.dev',
  //     lastUpdatedDate: '28th January 2025, 15:38',
  //     auditData: {
  //       stale: {
  //         danger: 3,
  //         warning: 5,
  //         pending: 2,
  //         compliant: 18,
  //         unchecked: 1,
  //         total: 29,
  //       },
  //       overprovisioned: {
  //         danger: 2,
  //         warning: 4,
  //         pending: 1,
  //         compliant: 22,
  //         unchecked: 3,
  //         total: 32,
  //       },
  //       security: {
  //         danger: 5,
  //         warning: 3,
  //         pending: 0,
  //         compliant: 15,
  //         unchecked: 2,
  //         total: 25,
  //       },
  //       network: {
  //         danger: 1,
  //         warning: 2,
  //         pending: 1,
  //         compliant: 10,
  //         unchecked: 0,
  //         total: 14,
  //       },
  //       storage: {
  //         danger: 4,
  //         warning: 2,
  //         pending: 3,
  //         compliant: 20,
  //         unchecked: 1,
  //         total: 30,
  //       },
  //       compute: {
  //         danger: 2,
  //         warning: 3,
  //         pending: 1,
  //         compliant: 25,
  //         unchecked: 2,
  //         total: 33,
  //       },
  //     },
  //     categoryIcons: {
  //       stale: 'server',
  //       overprovisioned: 'exclamation',
  //       security: 'shield',
  //     },
  //   },
  //   {
  //     id: 'account2',
  //     title: 'AWS Production',
  //     subtitle: '5 Services',
  //     status: 'READY',
  //     icon: 'cloud',
  //     providerType: 'aws',
  //     lastUpdatedBy: 'admin@zop.dev',
  //     lastUpdatedDate: '30th January 2025, 09:15',
  //     auditData: {
  //       network: {
  //         danger: 1,
  //         warning: 2,
  //         pending: 1,
  //         compliant: 10,
  //         unchecked: 0,
  //         total: 14,
  //       },
  //       storage: {
  //         danger: 4,
  //         warning: 2,
  //         pending: 3,
  //         compliant: 20,
  //         unchecked: 1,
  //         total: 30,
  //       },
  //       compute: {
  //         danger: 2,
  //         warning: 3,
  //         pending: 1,
  //         compliant: 25,
  //         unchecked: 2,
  //         total: 33,
  //       },
  //     },
  //     categoryIcons: {
  //       network: 'globe',
  //       storage: 'server',
  //       compute: 'server',
  //     },
  //   },
  //   {
  //     id: 'account3',
  //     title: 'Azure Development',
  //     subtitle: '3 Services',
  //     status: 'READY',
  //     providerType: 'azure',
  //     icon: 'cloud',
  //     lastUpdatedBy: 'developer@zop.dev',
  //     lastUpdatedDate: '25th January 2025, 11:22',
  //     auditData: {
  //       security: {
  //         danger: 7,
  //         warning: 2,
  //         pending: 1,
  //         compliant: 10,
  //         unchecked: 0,
  //         total: 20,
  //       },
  //       compute: {
  //         danger: 1,
  //         warning: 4,
  //         pending: 2,
  //         compliant: 15,
  //         unchecked: 3,
  //         total: 25,
  //       },
  //     },
  //     categoryIcons: {
  //       security: 'shield',
  //       compute: 'server',
  //     },
  //   },
  // ];

  const cloudAccounts = [];
  const auditCardData = {
    title: 'Audit Your Cloud',
    description:
      'Analyze your cloud infrastructure to identify stale resources, overprovisioned assets, and critical security vulnerabilities.',
    icon: <CloudIcon className="h-6 w-6 text-white" />,
    features: ['Ensure Compliance', 'Performance optimization', 'Cost efficiency recommendations'],
    buttonText: 'Audit Your Cloud',
    buttonIcon: <CloudIcon className="h-5 w-5 text-white" />,
    onClick: handleAuditClick,
  };

  const deployCardData = {
    title: 'Deploy Application',
    description: 'Deploy your apps quickly with our streamlined and reliable deployment process.',
    icon: <RocketLaunchIcon className="h-6 w-6 text-white" />,
    features: ['Automate Deployments', 'Improve Reliability', 'Scale Seamlessly'],
    buttonText: 'Deploy Application',
    buttonIcon: <RocketLaunchIcon className="h-5 w-5 text-white" />,
    buttonVariant: 'primary',
    onClick: () => {},
  };
  console.log('sadada');
  return (
    <main className="flex-1 p-6">
      <div
        className={`mx-auto ${cloudAccounts?.length === 0 ? 'max-w-5xl' : 'max-w-6xl'}  space-y-8`}
      >
        <div className="text-center">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
            Welcome to ZOP Dashboard
          </h1>
          <p className="mt-4 text-lg text-muted-foreground">
            Get started by choosing one of the options below
          </p>
        </div>

        <div className="flex flex-col md:flex-row gap-6">
          {cloudAccounts?.length === 0 ? (
            <>
              <DashboardSection>
                <DashBoardCard {...auditCardData} />
              </DashboardSection>
              <DashboardSection>
                <DashBoardCard {...deployCardData} />
              </DashboardSection>
            </>
          ) : (
            <>
              <div className="flex-1 flex flex-col gap-2">
                <div className="flex justify-between">
                  <h2 className="text-xl font-semibold">Cloud Accounts</h2>
                  <Button className="">Audit Cloud Accounts</Button>
                </div>
                <div className="border rounded-xl p-6 space-y-4 shadow-sm bg-white flex flex-col">
                  <div className="space-y-4 flex justify-center items-center flex-col">
                    <CloudAccountAuditCards cloudAccounts={cloudAccounts} />
                  </div>
                </div>
              </div>

              <div className="flex-1 flex flex-col gap-2 mt-12">
                <div className="border rounded-xl p-6 space-y-4 shadow-sm bg-white flex flex-col">
                  <DashBoardCard {...deployCardData} />
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </main>
  );
};

export default Home;
