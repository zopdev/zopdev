import IconButton from '@/components/atom/Button/IconButton';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import ResourceGroupManager from '@/components/container/resources/AddResourceGroup';
import { CloudResourceRow } from '@/components/container/resources/ResourceTableRow';
import Table from '@/components/molecules/Table';
import { ChevronRightIcon, PencilSquareIcon, TrashIcon } from '@heroicons/react/20/solid';
import React, { useState } from 'react';

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
                <ChevronRightIcon
                  className={`size-5 text-primary-600 mt-1 transition-transform ${expandedGroups.has(group.id) ? 'rotate-90' : ''}`}
                />
                <div>
                  <h3 className="font-semibold text-gray-900">{group.name}</h3>
                  <p className="text-sm text-gray-500">{group.description}</p>
                </div>
              </div>

              <div className="flex items-center gap-3 flex-wrap sm:flex-nowrap">
                {group?.resources?.length > 0 ? (
                  <span className="text-sm text-gray-600 whitespace-nowrap">
                    {group.resources?.filter((item) => item?.status === 'RUNNING')?.length}/
                    {group?.resources?.length} running
                  </span>
                ) : (
                  <span className="text-sm text-gray-400 whitespace-nowrap">
                    No resources have been added to the group
                  </span>
                )}
                <div className="flex space-x-2">
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
                    renderContentProps={{
                      resources: [],
                    }}
                  />

                  <IconButton
                    onClick={(e) => {
                      e.stopPropagation();
                      // onAction('delete', group);
                    }}
                  >
                    <TrashIcon className="size-4 text-red-400 " />
                  </IconButton>
                </div>
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

export default ResourceGroupAccordion;
