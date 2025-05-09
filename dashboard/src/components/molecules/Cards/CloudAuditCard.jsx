import React, { useEffect, useState } from 'react';
import {
  ExclamationCircleIcon,
  CheckCircleIcon,
  ClockIcon,
  CloudIcon,
  ServerIcon,
  ShieldCheckIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline';
import ResourceStatus from '@/components/atom/ResourceStatus/index.jsx';
import { PROVIDER_ICON_MAPPER } from '@/utils/componentMapper.jsx';
import { toast } from '@/components/molecules/Toast/index.jsx';
import ErrorComponent from '@/components/atom/ErrorComponent/index.jsx';

const ICONS = {
  cloud: CloudIcon,
  server: ServerIcon,
  shield: ShieldCheckIcon,
  globe: GlobeAltIcon,
  exclamation: ExclamationCircleIcon,
  check: CheckCircleIcon,
  clock: ClockIcon,
};

export default function CloudAccountAuditCards({ cloudAccounts = [], reRunAudit }) {
  return (
    <div className="space-y-4 flex justify-center items-center flex-col w-full ">
      {cloudAccounts?.map((account, index) => (
        <CloudAccountAuditCard
          key={account.id || index}
          {...account}
          account={account}
          reRunAudit={reRunAudit}
        />
      ))}
    </div>
  );
}

function CloudAccountAuditCard({
  // id,
  name,
  subtitle,
  status,
  provider,
  auditData = {},
  lastUpdatedBy,
  updatedAt,
  initialActiveTab,
  reRunAudit,
  account,
  categoryIcons = {},
  statusBarColors = {
    danger: 'bg-red-500',
    warning: 'bg-yellow-500',
    pending: 'bg-primary-500',
    compliant: 'bg-green-500',
    unchecked: 'bg-secondary-300',
  },
  statusIconColors = {
    danger: { bg: 'bg-red-100', icon: 'text-red-500' },
    warning: { bg: 'bg-yellow-100', icon: 'text-yellow-500' },
    pending: { bg: 'bg-primary-100', icon: 'text-primary-500' },
    compliant: { bg: 'bg-green-100', icon: 'text-green-500' },
    unchecked: { bg: 'bg-secondary-100', icon: 'text-secondary-500' },
  },
}) {
  const [activeTab, setActiveTab] = useState(
    initialActiveTab || (Object.keys(auditData).length > 0 ? Object.keys(auditData)[0] : ''),
  );
  const getCategoryIcon = (category) => {
    if (categoryIcons[category]) {
      const IconComponent = ICONS[categoryIcons[category]] || ICONS.exclamation;
      return <IconComponent className="h-4 w-4" />;
    }

    switch (category) {
      case 'stale':
        return <ServerIcon className="h-4 w-4" />;
      case 'overprovision':
        return <ExclamationCircleIcon className="h-4 w-4" />;
      case 'security':
        return <ShieldCheckIcon className="h-4 w-4" />;
      case 'network':
        return <GlobeAltIcon className="h-4 w-4" />;
      case 'storage':
        return <GlobeAltIcon className="h-4 w-4" />;
      case 'compute':
        return <ServerIcon className="h-4 w-4" />;
      // default:
      //   return <ExclamationCircleIcon className="h-4 w-4" />;
    }
  };

  const getStatusPercentage = (status, category) => {
    if (!auditData[category] || auditData[category].total === 0) return 0;
    return (auditData[category][status] / auditData[category].total) * 100;
  };

  const renderStatusBar = (category) => {
    if (!auditData[category]) return null;

    return (
      <div className="w-full h-2 flex rounded-full overflow-hidden mt-2">
        {Object.keys(statusBarColors).map((statusKey) => (
          <div
            key={statusKey}
            className={statusBarColors[statusKey]}
            style={{
              width: `${getStatusPercentage(statusKey, category)}%`,
              minWidth: auditData[category][statusKey] > 0 ? '4px' : '0',
            }}
          />
        ))}
      </div>
    );
  };

  // const handleRerun = () => {
  //   reRunAudit.mutate({
  //     id,
  //     selectedOption: activeTab,
  //   });
  // };
  useEffect(() => {
    if (reRunAudit.isError) toast.failed(reRunAudit.error?.message);
  }, [reRunAudit]);

  const renderStatusDetails = (category) => {
    if (!auditData[category]) return null;

    const statuses = Object.keys(statusIconColors);

    return (
      <div className="flex flex-wrap justify-around gap-y-3 mt-4">
        {statuses.map((statusKey) => {
          if (typeof auditData[category][statusKey] === 'undefined') return null;

          let StatusIcon = statusKey === 'pending' ? ClockIcon : ExclamationCircleIcon;
          if (statusKey === 'compliant') StatusIcon = CheckCircleIcon;

          const bgColor = statusIconColors[statusKey].bg;
          const iconColor = statusIconColors[statusKey].icon;

          return (
            <div key={statusKey} className="flex flex-col items-center w-16 sm:w-auto">
              <div className={`flex items-center justify-center w-8 h-8 rounded-full ${bgColor}`}>
                <StatusIcon className={`h-5 w-5 ${iconColor}`} />
              </div>
              <span className="text-xs font-medium mt-1 text-center">
                {statusKey.charAt(0).toUpperCase() + statusKey.slice(1)}
              </span>
              <span className="text-lg">{auditData[category][statusKey]}</span>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className="w-full lg:max-w-lg bg-white rounded-lg overflow-hidden shadow-sm min-h-78">
      <div className="p-3 sm:p-4 sm:pb-2 mt-2">
        <div className="flex items-center space-x-2">
          <div className="flex items-center justify-center w-8 h-8 rounded-full shrink-0">
            {PROVIDER_ICON_MAPPER[provider]}
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="text-left font-medium text-gray-600 text-xl">{name}</h3>
              <ResourceStatus status={status} />
            </div>
            <p className="text-xs sm:text-sm text-secondary-500 truncate">{subtitle}</p>
          </div>
        </div>
      </div>
      {!account?.auditDetails?.status && (
        <div className={'mx-2'}>
          <ErrorComponent
            complete={true}
            errorText={'Something went wrong while fetching audit details.'}
            className={'!w-full !h-56'}
          />
        </div>
      )}
      {account?.auditDetails?.data && (
        <div className="p-3 sm:p-4">
          <div className="w-full">
            {Object.keys(auditData).length > 0 && (
              <div className="w-full bg-secondary-100 flex justify-center items-center rounded-lg">
                <div className="w-full overflow-x-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
                  <div className="flex min-w-full p-1 whitespace-nowrap">
                    {Object.keys(auditData).map((category) => (
                      <button
                        key={category}
                        onClick={() => setActiveTab(category)}
                        className={`flex items-center gap-1 px-2 cursor-pointer sm:px-3 py-1.5 text-xs sm:text-sm font-medium rounded-md transition-colors flex-1 justify-center ${
                          activeTab === category
                            ? 'bg-white shadow-sm'
                            : 'text-secondary-600 hover:bg-secondary-200'
                        }`}
                      >
                        <span className="block">{getCategoryIcon(category)}</span>
                        <span className="block truncate">
                          {category.charAt(0).toUpperCase() + category.slice(1)}
                        </span>
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            )}

            <div className="mt-4">
              {Object.keys(auditData).map((category) => (
                <div key={category} className={`${activeTab === category ? 'block' : 'hidden'}`}>
                  {renderStatusBar(category)}
                  {renderStatusDetails(category)}
                </div>
              ))}
            </div>
          </div>

          {(lastUpdatedBy || updatedAt) && (
            <div className="mt-1 pt-3 text-xs flex text-secondary-900 justify-between items-center">
              <div className={'md:flex'}>
                {' '}
                <span className={'text-secondary-400'}>Last Run on&nbsp;</span>
                {updatedAt && <p>{updatedAt}</p>}
              </div>

              {/*<Button*/}
              {/*  loading={reRunAudit?.isPending}*/}
              {/*  onClick={() => handleRerun()}*/}
              {/*  variant={'primary-outline'}*/}
              {/*  size={'sm'}*/}
              {/*>*/}
              {/*  Re-Run {activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}*/}
              {/*</Button>*/}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
