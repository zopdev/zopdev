import SwitchButton from '@/components/atom/Switch';
import Table from '@/components/molecules/Table';
import { toast } from '@/components/molecules/Toast';
import { usePostResourceState } from '@/queries/cloud-resources';
import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

const CloudResourceRow = (resource) => {
  const [currentState, setCurrentState] = useState(resource?.status === 'RUNNING' ? true : false);
  const { cloudId } = useParams();
  const resourceStateChanger = usePostResourceState();

  const handleToggle = (state) => {
    setCurrentState(state);
    resourceStateChanger.mutate({
      cloudAccId: parseInt(cloudId),
      id: resource?.id,
      name: resource?.name,
      type: resource?.type,
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
    id: resource?.id,
    name: resource?.name,
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
    instance_type: resource?.type,
    region: resource?.region,
  };
};

const ResourceGroupAccordion = ({ groups = [], defaultExpandedIds = [] }) => {
  const [expandedGroups, setExpandedGroups] = useState(new Set(defaultExpandedIds));

  const toggleGroup = (groupId) => {
    const newExpanded = new Set(expandedGroups);
    newExpanded.has(groupId) ? newExpanded.delete(groupId) : newExpanded.add(groupId);
    setExpandedGroups(newExpanded);
  };

  const tableHeaders = [
    { key: 'name', label: 'Name', align: 'left', width: '200px' },
    { key: 'state', label: 'State', align: 'left', width: '150px' },
    { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
    { key: 'region', label: 'Region', align: 'left', width: '120px' },
  ];

  return (
    <div className="space-y-4 py-4">
      {groups.map((group) => (
        <div key={group.id} className="border border-gray-200 rounded-lg overflow-hidden">
          <div
            className={`bg-gray-50 px-4 py-3 sm:px-6  transition-colors ${
              group?.resources?.length > 0 && 'cursor-pointer hover:bg-gray-100'
            } `}
            onClick={() => group?.resources?.length > 0 && toggleGroup(group.id)}
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
                {/* <span className="text-sm text-gray-600 whitespace-nowrap">
                  {group.runningResources}/{group.totalResources} running
                </span> */}
                {/* <div className="flex space-x-2">
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
                </div> */}
              </div>
            </div>
          </div>

          {expandedGroups.has(group.id) && (
            <Table
              headers={tableHeaders}
              data={group.resources || []}
              enableRowClick={false}
              renderRow={CloudResourceRow}
              emptyStateTitle="No Resources added"
            />
          )}
        </div>
      ))}
    </div>
  );
};

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

export default ResourceGroupAccordion;
