import React, { useEffect, useMemo, useState } from 'react';
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
import { useGetAuditDetails, usePostAuditData } from '@/queries/CloudAccount/index.js';
import Button from '@/components/atom/Button/index.jsx';
import { refreshInterval } from '@/utils/constant.js';

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
          account={account}
          reRunAudit={reRunAudit}
        />
      ))}
    </div>
  );
}

function CloudAccountAuditCard({ account }) {
  function transformAuditApiResponse(apiResponse) {
    const getAuditSummary = (items) => {
      const summary = {
        danger: 0,
        warning: 0,
        pending: 0,
        compliant: 0,
        unchecked: 0,
        total: 0,
      };

      items.forEach((item) => {
        const status = item.status || 'unchecked';
        if (Object.prototype.hasOwnProperty.call(summary, status)) {
          summary[status]++;
        } else {
          summary.unchecked++;
        }
      });

      summary.total = Object.values(summary).reduce((acc, val) => acc + val, 0);
      return summary;
    };

    const rawAuditData = {};
    const categoryIcons = {};
    const allData = [];

    const getDefaultIcon = (category) => {
      const iconMap = {
        overprovision: 'exclamation',
        stale: 'clock',
        security: 'shield',
      };
      return iconMap[category] || 'circle';
    };

    for (const category in apiResponse) {
      const entries = apiResponse[category] || [];
      const allItems = entries.flatMap((entry) => entry.result?.items || []);
      rawAuditData[category] = getAuditSummary(allItems);
      categoryIcons[category] = getDefaultIcon(category);

      entries.forEach(({ result, ...rest }) => {
        allData.push({
          category,
          ...rest,
        });
      });
    }

    const allSummary = {
      danger: 0,
      warning: 0,
      pending: 0,
      compliant: 0,
      unchecked: 0,
      total: 0,
    };

    for (const category in rawAuditData) {
      const summary = rawAuditData[category];
      for (const key in allSummary) {
        allSummary[key] += summary[key] || 0;
      }
    }
    const auditData = {
      all: allSummary,
      ...rawAuditData,
    };

    return {
      auditData,
      categoryIcons,
      allData,
    };
  }
  let categoryIcons = {};
  let statusBarColors = {
    danger: 'bg-red-500',
    warning: 'bg-yellow-500',
    pending: 'bg-primary-500',
    compliant: 'bg-green-500',
    unchecked: 'bg-secondary-300',
  };
  let statusIconColors = {
    danger: { bg: 'bg-red-100', icon: 'text-red-500' },
    warning: { bg: 'bg-yellow-100', icon: 'text-yellow-500' },
    pending: { bg: 'bg-primary-100', icon: 'text-primary-500' },
    compliant: { bg: 'bg-green-100', icon: 'text-green-500' },
    unchecked: { bg: 'bg-secondary-100', icon: 'text-secondary-500' },
  };

  const [shouldPoll, setShouldPoll] = useState(false);
  const { data: auditResponseData } = useGetAuditDetails(
    { id: account?.id },
    {
      enabled: !!account?.id,
      refetchInterval: shouldPoll ? refreshInterval : false,
    },
  );
  const individualAuditData = auditResponseData?.data || {};
  const transformedData = useMemo(() => {
    return transformAuditApiResponse(individualAuditData);
  }, [individualAuditData]);

  const auditData = transformedData?.auditData;
  useEffect(() => {
    if (auditData?.all?.pending > 0) {
      setShouldPoll(true);
    } else {
      setShouldPoll(false);
    }
  }, [auditData]);

  const reRunAudit = usePostAuditData();

  const [activeTab, setActiveTab] = useState(
    Object.keys(auditData).length > 0 ? Object.keys(auditData)[0] : '',
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
    const data = auditData[category];
    if (!data) return null;

    const total = data.total;

    const allZero = Object.keys(statusBarColors).every((statusKey) => data[statusKey] === 0);
    if (total === 0 || allZero) {
      return <div className="w-full h-2 mt-2 rounded-full overflow-hidden bg-gray-300" />;
    }

    return (
      <div className="w-full h-2 flex rounded-full overflow-hidden mt-2">
        {Object.keys(statusBarColors).map((statusKey) => (
          <div
            key={statusKey}
            className={statusBarColors[statusKey]}
            style={{
              width: `${getStatusPercentage(statusKey, category)}%`,
              minWidth: data[statusKey] > 0 ? '4px' : '',
            }}
          />
        ))}
      </div>
    );
  };

  const handleRerun = () => {
    reRunAudit.mutate({
      id: account?.id,
      selectedOption: activeTab,
    });
  };
  useEffect(() => {
    if (reRunAudit?.isError) toast.failed(reRunAudit.error?.message);
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
  const getLatestEvaluatedAt = (category, allData) => {
    if (category === 'all') {
      const timestamps = allData
        .filter((item) => item.evaluatedAt)
        .map((item) => new Date(item.evaluatedAt));
      if (timestamps.length === 0) return null;
      const latestTimestamp = new Date(Math.max(...timestamps));

      return latestTimestamp.toISOString();
    }
    const timestamps = allData
      .filter((item) => item.category === category && item.evaluatedAt)
      .map((item) => new Date(item.evaluatedAt));
    if (timestamps.length === 0) return null;

    const latestTimestamp = new Date(Math.max(...timestamps));

    return latestTimestamp.toISOString();
  };

  const lastEvaluatedAt = getLatestEvaluatedAt(activeTab, transformedData?.allData);

  const shimmerClass =
    'relative overflow-hidden before:absolute before:inset-0 before:-translate-x-full before:animate-shimmer  animate-pulse before:bg-gradient-to-r before:from-transparent before:via-white/60 before:to-transparent';

  return (
    <div className="w-full lg:max-w-lg bg-white rounded-lg overflow-hidden shadow-sm min-h-78">
      <div className="p-3 sm:p-4 sm:pb-2 mt-2">
        <div className="flex items-center space-x-2">
          <div className="flex items-center justify-center w-8 h-8 rounded-full shrink-0">
            {PROVIDER_ICON_MAPPER[account?.provider]}
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="text-left font-medium text-gray-600 text-xl">{account?.name}</h3>
              <ResourceStatus status={account?.status} />
            </div>
            {/*<p className="text-xs sm:text-sm text-secondary-500 truncate">{subtitle}</p>*/}
          </div>
        </div>
      </div>

      {auditResponseData?.isError && (
        <div className="mx-2">
          <ErrorComponent
            complete
            errorText="Something went wrong while fetching audit details."
            className="!w-full !h-56"
          />
        </div>
      )}

      {auditResponseData?.isLoading && <Skeleton shimmerClass={shimmerClass} />}

      {!auditResponseData?.isError && !auditResponseData?.isLoading && (
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
          {lastEvaluatedAt && (
            <div className="mt-1 pt-3 text-xs flex text-secondary-900 justify-between items-center">
              <div className="md:flex">
                <span className="text-secondary-400">Last Run on&nbsp;</span>
                <p>{new Date(lastEvaluatedAt).toLocaleString()}</p>
              </div>
              <Button
                loading={reRunAudit?.isPending}
                onClick={() => handleRerun()}
                variant={'primary-outline'}
                size={'sm'}
              >
                Re-Run {activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}
              </Button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

const Skeleton = ({ shimmerClass }) => (
  <div className="p-3 sm:p-4">
    <div className="w-full bg-secondary-100 flex justify-center items-center rounded-lg">
      <div className="w-full overflow-x-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
        <div className="flex min-w-full p-1 whitespace-nowrap gap-2">
          {Array.from({ length: 4 }).map((_, index) => (
            <div
              key={index}
              className="flex items-center gap-2 px-2 sm:px-3 py-1.5 rounded-md w-28 sm:w-32 h-9"
            >
              <div className={`h-3 bg-secondary-200 rounded w-3/4 ${shimmerClass}`} />
            </div>
          ))}
        </div>
      </div>
    </div>

    <div className={`w-full h-2 mt-6 rounded-full bg-secondary-200 ${shimmerClass}`}></div>

    <div className="flex flex-wrap justify-around gap-y-3 mt-4">
      {Array.from({ length: 5 }).map((_, index) => (
        <div key={index} className="flex flex-col items-center w-16 sm:w-auto">
          <div className={`w-8 h-8 bg-secondary-200 rounded-full ${shimmerClass}`}></div>
        </div>
      ))}
    </div>
    <div className="mt-4 pt-3 flex justify-between items-center">
      <div className={`h-4 w-full bg-secondary-200 ${shimmerClass} rounded`}></div>
    </div>
  </div>
);
