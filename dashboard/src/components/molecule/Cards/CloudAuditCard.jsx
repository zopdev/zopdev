import { useState } from 'react';
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
import { PROVIDER_ICON_MAPPER } from '@/componentMapper.jsx';

const ICONS = {
  cloud: CloudIcon,
  server: ServerIcon,
  shield: ShieldCheckIcon,
  globe: GlobeAltIcon,
  exclamation: ExclamationCircleIcon,
  check: CheckCircleIcon,
  clock: ClockIcon,
};

export default function CloudAccountAuditCards({ cloudAccounts = [] }) {
  return (
    <div className="space-y-4 flex justify-center items-center flex-col">
      {cloudAccounts?.map((account, index) => (
        <CloudAccountAuditCard key={account.id || index} {...account} />
      ))}
    </div>
  );
}

function CloudAccountAuditCard({
  title,
  subtitle,
  status,
  providerType,
  auditData = {},
  lastUpdatedBy,
  lastUpdatedDate,
  initialActiveTab,
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
      case 'overprovisioned':
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
              minWidth: auditData[category][statusKey] > 0 ? '4px' : '0', // Ensure visibility for small values
            }}
          />
        ))}
      </div>
    );
  };

  const renderStatusDetails = (category) => {
    if (!auditData[category]) return null;

    const statuses = Object.keys(statusIconColors);

    return (
      <div className="flex justify-between mt-4">
        {statuses.map((statusKey) => {
          if (typeof auditData[category][statusKey] === 'undefined') return null;

          let StatusIcon = statusKey === 'pending' ? ClockIcon : ExclamationCircleIcon;
          if (statusKey === 'compliant') StatusIcon = CheckCircleIcon;

          const bgColor = statusIconColors[statusKey].bg;
          const iconColor = statusIconColors[statusKey].icon;

          return (
            <div key={statusKey} className="flex flex-col items-center">
              <div className={`flex items-center justify-center w-8 h-8 rounded-full ${bgColor}`}>
                <StatusIcon className={`h-5 w-5 ${iconColor}`} />
              </div>
              <span className="text-xs font-medium mt-1">
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
    <>
      <div className="w-full max-w-lg bg-white rounded-lg border border-borderDefault overflow-hidden">
        {/* Card Header */}
        <div className="p-4 pb-2">
          <div className="flex items-center space-x-2">
            <div className="flex items-center justify-center w-8 h-8 rounded-full ">
              {PROVIDER_ICON_MAPPER[providerType]}
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-3">
                <h3 className="text-xl font-semibold">{title}</h3>
                <ResourceStatus status={status} />
              </div>
              <p className="text-sm text-secondary-500">{subtitle}</p>
            </div>
          </div>
        </div>

        <div className="p-4">
          <div className="w-full">
            {Object.keys(auditData).length > 0 && (
              <div className="w-full overflow-x-auto pb-2 [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden">
                <div className="inline-flex min-w-full p-1 bg-secondary-100 rounded-lg">
                  {Object.keys(auditData).map((category) => (
                    <button
                      key={category}
                      onClick={() => setActiveTab(category)}
                      className={`flex items-center gap-1 px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                        activeTab === category
                          ? 'bg-white shadow-sm'
                          : 'text-secondary-600 hover:bg-secondary-200'
                      }`}
                    >
                      {getCategoryIcon(category)}
                      <span>{category.charAt(0).toUpperCase() + category.slice(1)}</span>
                    </button>
                  ))}
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

          {(lastUpdatedBy || lastUpdatedDate) && (
            <div className="mt-4 pt-4 border-t text-xs text-secondary-500">
              {/* {lastUpdatedBy && <p>Updated By {lastUpdatedBy}</p>} */}
              {lastUpdatedDate && <p>{lastUpdatedDate}</p>}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
