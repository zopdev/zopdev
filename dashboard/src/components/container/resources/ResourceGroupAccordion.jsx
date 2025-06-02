import React, { useState } from 'react';

const ResourceGroupAccordion = () => {
  const [resourceGroups] = useState([
    {
      id: 1,
      name: 'Production Services',
      description: 'Core production services that need to run 24/7',
      totalResources: 2,
      runningResources: 2,
      resources: [
        {
          id: 1,
          name: 'web-server-prod',
          type: 'EC2',
          status: 'Running',
        },
        {
          id: 2,
          name: 'order-processing-db',
          type: 'RDS',
          status: 'Running',
        },
      ],
    },
    {
      id: 2,
      name: 'Development Environment',
      description: 'Development and testing resources that can be shut down outside of work hours',
      totalResources: 3,
      runningResources: 0,
      resources: [
        {
          id: 3,
          name: 'dev-web-server',
          type: 'EC2',
          status: 'Stopped',
        },
        {
          id: 4,
          name: 'test-database',
          type: 'RDS',
          status: 'Stopped',
        },
        {
          id: 5,
          name: 'dev-cache',
          type: 'ElastiCache',
          status: 'Stopped',
        },
      ],
    },
    {
      id: 3,
      name: 'Analytics Platform',
      description: 'Data processing and analytics services',
      totalResources: 1,
      runningResources: 1,
      resources: [
        {
          id: 6,
          name: 'analytics-cluster',
          type: 'EMR',
          status: 'Running',
        },
      ],
    },
  ]);

  const [expandedGroups, setExpandedGroups] = useState(new Set([1]));

  const toggleGroup = (groupId) => {
    const newExpanded = new Set(expandedGroups);
    newExpanded.has(groupId) ? newExpanded.delete(groupId) : newExpanded.add(groupId);
    setExpandedGroups(newExpanded);
  };

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case 'running':
        return 'bg-green-100 text-green-800';
      case 'stopped':
        return 'bg-red-100 text-red-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const handleAction = (action, resourceName) => {
    console.log(`${action} action triggered for ${resourceName}`);
  };

  return (
    <div className="space-y-4 py-4">
      {resourceGroups.map((group) => (
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
                  {/* Start */}
                  <button
                    className="p-2 text-primary-600 hover:bg-primary-50 rounded"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleAction('start', group.name);
                    }}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1M12 2l3 7H9l3-7zm0 0v18"
                      />
                    </svg>
                  </button>
                  {/* Stop */}
                  <button
                    className="p-2 text-red-600 hover:bg-red-50 rounded"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleAction('stop', group.name);
                    }}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 10h6v4H9z"
                      />
                    </svg>
                  </button>
                  {/* Edit */}
                  <button
                    className="p-2 text-gray-600 hover:bg-gray-50 rounded"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleAction('edit', group.name);
                    }}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                  </button>
                  {/* Delete */}
                  <button
                    className="p-2 text-red-600 hover:bg-red-50 rounded"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleAction('delete', group.name);
                    }}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Table Section */}
          {expandedGroups.has(group.id) && (
            <div className="bg-white overflow-x-auto">
              <table className="min-w-full">
                <thead className="bg-gray-50 border-t border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Resource
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {group.resources.map((resource) => (
                    <tr key={resource.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {resource.name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {resource.type}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(resource.status)}`}
                        >
                          {resource.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <button
                          className="text-red-600 hover:bg-red-50 p-1 rounded"
                          onClick={() => handleAction('delete', resource.name)}
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                          </svg>
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default ResourceGroupAccordion;
