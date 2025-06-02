import Table from '@/components/molecules/Table';
import React, { useState } from 'react';

const ResourceGroupAccordion = ({ groups = [], defaultExpandedIds = [], onAction = () => {} }) => {
  const [expandedGroups, setExpandedGroups] = useState(new Set(defaultExpandedIds));

  const toggleGroup = (groupId) => {
    const newExpanded = new Set(expandedGroups);
    newExpanded.has(groupId) ? newExpanded.delete(groupId) : newExpanded.add(groupId);
    setExpandedGroups(newExpanded);
  };

  // Define table headers
  const tableHeaders = [
    { key: 'name', label: 'Resource', align: 'left', width: '200px' },
    { key: 'type', label: 'Type', align: 'left', width: '150px' },
    { key: 'status', label: 'Status', align: 'left', width: '120px' },
    { key: 'actions', label: 'Actions', align: 'left', width: '100px' },
  ];

  return (
    <div className="space-y-4 py-4">
      {groups.map((group) => (
        <div key={group.id} className="border border-gray-200 rounded-lg overflow-hidden">
          {/* Header */}
          <div
            className="bg-gray-50 px-4 py-3 sm:px-6 cursor-pointer hover:bg-gray-100 transition-colors"
            onClick={() => toggleGroup(group.id)}
          >
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="flex items-start gap-3">
                <svg
                  className={`w-4 h-4 text-primary-600 mt-1 transition-transform ${expandedGroups.has(group.id) ? 'rotate-90' : ''}`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 5l7 7-7 7"
                  />
                </svg>
                <div>
                  <h3 className="font-semibold text-gray-900">{group.name}</h3>
                  <p className="text-sm text-gray-500">{group.description}</p>
                </div>
              </div>

              <div className="flex items-center gap-3 flex-wrap sm:flex-nowrap">
                <span className="text-sm text-gray-600 whitespace-nowrap">
                  {group.runningResources}/{group.totalResources} running
                </span>
                <div className="flex space-x-2">
                  {/* Actions */}
                  {['start', 'stop', 'edit', 'delete'].map((action) => (
                    <button
                      key={action}
                      className={`p-2 rounded ${
                        action === 'stop' || action === 'delete'
                          ? 'text-red-600 hover:bg-red-50'
                          : action === 'start'
                            ? 'text-primary-600 hover:bg-primary-50'
                            : 'text-gray-600 hover:bg-gray-50'
                      }`}
                      onClick={(e) => {
                        e.stopPropagation();
                        onAction(action, group);
                      }}
                    >
                      <ActionIcon type={action} />
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Your Table Component */}
          {expandedGroups.has(group.id) && (
            <Table
              headers={tableHeaders}
              data={group.resources || []}
              enableRowClick={false}
              // renderRow={CloudResourceRow}
              emptyStateTitle="No Resources Found"
              // emptyStateDescription="Looks like your cloud account has no active resources right now"
            />
          )}
        </div>
      ))}
    </div>
  );
};

// Sub-component for rendering action icons
const ActionIcon = ({ type }) => {
  switch (type) {
    case 'start':
      return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1M12 2l3 7H9l3-7zm0 0v18"
          />
        </svg>
      );
    case 'stop':
      return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 10h6v4H9z" />
        </svg>
      );
    case 'edit':
      return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
          />
        </svg>
      );
    case 'delete':
      return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1v3M4 7h16"
          />
        </svg>
      );
    default:
      return null;
  }
};

// Demo with sample data
const App = () => {
  const sampleGroups = [
    {
      id: 1,
      name: 'Production Resources',
      description: 'Critical production infrastructure',
      runningResources: 3,
      totalResources: 5,
      resources: [
        { id: 1, name: 'Web Server 1', type: 'EC2', status: 'running' },
        { id: 2, name: 'Database Server', type: 'RDS', status: 'running' },
        { id: 3, name: 'Load Balancer', type: 'ALB', status: 'stopped' },
        { id: 4, name: 'Cache Server', type: 'ElastiCache', status: 'pending' },
        { id: 5, name: 'File Storage', type: 'S3', status: 'running' },
      ],
    },
    {
      id: 2,
      name: 'Development Resources',
      description: 'Development and testing environment',
      runningResources: 1,
      totalResources: 3,
      resources: [
        { id: 6, name: 'Dev Server', type: 'EC2', status: 'running' },
        { id: 7, name: 'Test Database', type: 'RDS', status: 'pending' },
        { id: 8, name: 'Dev Cache', type: 'ElastiCache', status: 'stopped' },
      ],
    },
    {
      id: 3,
      name: 'Staging Resources',
      description: 'Pre-production staging environment',
      runningResources: 0,
      totalResources: 0,
      resources: [],
    },
  ];

  const handleAction = (action, item) => {
    console.log(`Action: ${action}`, item);
  };

  return (
    <div>
      <ResourceGroupAccordion
        groups={sampleGroups}
        defaultExpandedIds={[1]}
        onAction={handleAction}
      />
    </div>
  );
};

export default App;
