import React, { useState } from 'react';
import { ChevronRightIcon, PencilSquareIcon, TrashIcon } from '@heroicons/react/20/solid';
import IconButton from '@/components/atom/Button/IconButton';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import ResourceGroupManager from '@/components/container/resources/AddResourceGroup';
import { CloudResourceRow } from '@/components/container/resources/ResourceTableRow';
import Table from '@/components/molecules/Table';

const tableHeaders = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const ResourceGroupAccordion = ({ groups = [], defaultExpandedIds = [], resources }) => {
  const [expandedGroups, setExpandedGroups] = useState(new Set(defaultExpandedIds));

  const toggleGroup = (groupId) => {
    setExpandedGroups((prev) => {
      const updated = new Set(prev);
      updated.has(groupId) ? updated.delete(groupId) : updated.add(groupId);
      return updated;
    });
  };

  return (
    <div className="space-y-4 py-4">
      {groups.map((group) => {
        const hasResources = group?.resources?.length > 0;
        const isExpanded = expandedGroups.has(group.id);
        const runningCount = group.resources?.filter((r) => r.status === 'RUNNING')?.length || 0;

        return (
          <div key={group.id} className="border border-gray-200 rounded-lg overflow-hidden">
            <div
              className={`bg-gray-50 px-4 py-3 sm:px-6 transition-colors ${
                hasResources ? 'cursor-pointer hover:bg-gray-100' : ''
              }`}
              onClick={() => hasResources && toggleGroup(group.id)}
            >
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex items-start gap-3">
                  <ChevronRightIcon
                    className={`size-5 text-primary-600 mt-1 transition-transform ${
                      isExpanded ? 'rotate-90' : ''
                    }`}
                  />
                  <div>
                    <h3 className="font-semibold text-gray-900">{group.name}</h3>
                    <p className="text-sm text-gray-500">{group.description}</p>
                  </div>
                </div>

                <div className="flex items-center gap-3 flex-wrap sm:flex-nowrap">
                  {hasResources ? (
                    <span className="text-sm text-gray-600 whitespace-nowrap">
                      {runningCount}/{group.resources.length} running
                    </span>
                  ) : (
                    <span className="text-sm text-gray-400 whitespace-nowrap">
                      No resources have been added to the group
                    </span>
                  )}

                  <div className="flex space-x-2" onClick={(e) => e.stopPropagation()}>
                    <FullScreenOverlay
                      customCTA={
                        <IconButton>
                          <PencilSquareIcon className="size-4 text-gray-500" />
                        </IconButton>
                      }
                      title="Create Resource Group"
                      size="xl"
                      variant="drawer"
                      renderContent={ResourceGroupManager}
                      renderContentProps={{ resources, initialData: group }}
                    />
                    <IconButton
                      onClick={() => {
                        /* TODO: Add delete logic */
                      }}
                    >
                      <TrashIcon className="size-4 text-red-400" />
                    </IconButton>
                  </div>
                </div>
              </div>
            </div>

            {isExpanded && (
              <Table
                headers={tableHeaders}
                data={group.resources || []}
                enableRowClick={false}
                renderRow={CloudResourceRow}
                emptyStateTitle="No Resources added"
              />
            )}
          </div>
        );
      })}
    </div>
  );
};

export default ResourceGroupAccordion;
