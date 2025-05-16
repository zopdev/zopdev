import { CloudIcon, RocketLaunchIcon } from '@heroicons/react/24/outline';
import DashBoardCard from '@/components/molecules/Cards/DashBoardCard.jsx';
import DashboardSection from '@/components/organisms/DashBoardSection.jsx';
import { useNavigate } from 'react-router-dom';
import CloudAccountAuditCards from '@/components/container/CloudAuditCard.jsx';
import { useGetCloudAccounts } from '@/queries/CloudAccount/index.js';
import ErrorComponent from '@/components/atom/ErrorComponent/index.jsx';
import React from 'react';
import CompleteLoader from '@/components/atom/Loaders/CompleteLoader.jsx';
import Button from '@/components/atom/Button/index.jsx';
import { toast } from '@/components/molecules/Toast/index.jsx';

const Dashboard = () => {
  const getData = useGetCloudAccounts();

  const navigate = useNavigate();

  const handleAuditClick = () => {
    navigate('/cloud-setup');
  };

  const handleApplicationClick = () => {
    toast.info('This feature is currently not available. Coming Soon!');
  };

  const auditCardData = {
    title: 'Audit Your Cloud',
    description:
      'Analyze your cloud infrastructure to identify stale resources, overprovision assets, and critical security vulnerabilities.',
    icon: <CloudIcon className="h-6 w-6 text-white" />,
    features: ['Ensure compliance', 'Performance optimization', 'Cost efficiency recommendations'],
    buttonText: 'Audit Your Cloud',
    buttonIcon: <CloudIcon className="h-5 w-5 text-white" />,
    onClick: handleAuditClick,
  };

  const deployCardData = {
    title: 'Deploy Application',
    description: 'Deploy your apps quickly with our streamlined and reliable deployment process.',
    icon: <RocketLaunchIcon className="h-6 w-6 text-white" />,
    features: ['Automate deployments', 'Improve reliability', 'Scale seamlessly'],
    buttonText: 'Deploy Application',
    buttonIcon: <RocketLaunchIcon className="h-5 w-5 text-white" />,
    buttonVariant: 'primary',
    onClick: handleApplicationClick,
  };

  return (
    <main className="flex-1 p-6">
      {getData?.isLoading && <CompleteLoader />}
      <div
        className={`mx-auto ${getData?.data?.data?.length === 0 ? 'max-w-5xl' : 'max-w-6xl'}  space-y-8`}
      >
        {getData?.isSuccess && !getData?.data?.data?.length > 0 && (
          <div className="text-center">
            <h1 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
              Welcome to Zopdev
            </h1>
            <p className="mt-4 text-lg text-muted-foreground">
              Get started by choosing one of the options below.
            </p>
          </div>
        )}

        <div className="flex flex-col md:justify-center lg:flex-row gap-6 my-5">
          {getData?.isSuccess && !getData?.data?.data?.length > 0 ? (
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
              <div className="flex-1 flex flex-col w-full">
                <div className="flex justify-between items-center">
                  <h2 className="text-left font-medium text-gray-600 text-xl">Cloud Accounts</h2>
                  <Button href={'cloud-setup'} variant="text-link" size="md">
                    Start New Cloud Audit
                  </Button>
                </div>
                <div className="border border-borderDefault rounded-xl p-6 space-y-4 shadow-sm bg-white flex flex-col">
                  <div className="space-y-4 flex justify-center items-center flex-col !min-h-78">
                    {getData?.isError && (
                      <div className={'!min-h-78 w-full'}>
                        <ErrorComponent
                          complete={true}
                          errorText={getData?.error?.message}
                          className={'!w-full !min-h-78'}
                        />
                      </div>
                    )}
                    {getData?.isSuccess && getData?.data?.data?.length > 0 && (
                      <CloudAccountAuditCards cloudAccounts={getData?.data?.data} />
                    )}
                  </div>
                </div>
              </div>

              <div className="flex-1 flex flex-col gap-2 lg:mt-12 w-full">
                <div className="border border-borderDefault rounded-xl p-6 space-y-4 shadow-sm bg-white flex flex-col">
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

export default Dashboard;
