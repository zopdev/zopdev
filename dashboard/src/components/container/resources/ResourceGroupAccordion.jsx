import React, { useState } from 'react';
import {
  ChevronRightIcon,
  PencilSquareIcon,
  //  TrashIcon
} from '@heroicons/react/20/solid';
import IconButton from '@/components/atom/Button/IconButton';
import FullScreenOverlay from '@/components/atom/FullScreenOverlay';
import ResourceGroupManager from '@/components/container/resources/AddResourceGroup';
import { CloudResourceRow } from '@/components/container/resources/ResourceTableRow';
import Table from '@/components/molecules/Table';
import { TrashIcon } from '@heroicons/react/24/outline/index.js';
import DeleteModal from '@/components/organisms/DeleteModal.jsx';
import { toast } from '@/components/molecules/Toast/index.jsx';

const tableHeaders = [
  { key: 'name', label: 'Name', align: 'left', width: '200px' },
  { key: 'state', label: 'State', align: 'left', width: '150px' },
  { key: 'instance_type', label: 'Instance Type', align: 'left', width: '120px' },
  { key: 'region', label: 'Region', align: 'left', width: '120px' },
];

const ResourceGroupAccordion = ({
  groups = [],
  defaultExpandedIds = [],
  resources,
  resourceDelete,
}) => {
  const [expandedGroups, setExpandedGroups] = useState(new Set(defaultExpandedIds));

  const toggleGroup = (groupId) => {
    setExpandedGroups((prev) => {
      const updated = new Set(prev);
      updated.has(groupId) ? updated.delete(groupId) : updated.add(groupId);
      return updated;
    });
  };
  const [deleteItem, setDeleteItem] = useState({});

  const deleteConfirmation = () => {
    resourceDelete.mutate(
      {
        cloudAccId: deleteItem?.cloud_account_id,
        resourceGroupId: deleteItem?.id,
      },
      {
        onSuccess: () => {
          toast.success('Resource Group deleted successfully!');
        },
        onError: () => {
          toast.failed('Failed to delete Resource Group.');
        },
      },
    );
  };

  return (
    <div className="space-y-4 py-4">
      {groups?.map((group) => {
        const hasResources = group?.resources?.length > 0;
        const isExpanded = expandedGroups.has(group.id);
        const runningCount = group.resources?.filter((r) => r.status === 'RUNNING')?.length || 0;
        return (
          <div key={group.id} className="border border-gray-200 rounded-lg overflow-hidden">
            <div
              className={`bg-gray-50 px-4 py-3 sm:px-6 transition-colors ${
                hasResources ? 'cursor-pointer hover:bg-gray-100' : ''
              }`}
              onClick={() => {
                if (hasResources) {
                  toggleGroup(group.id);
                }
              }}
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
                  </div>
                  <DeleteModal
                    onDelete={deleteConfirmation}
                    deleteTitle={'Delete Resource Group'}
                    deleteKey={deleteItem?.name}
                    isConfirmation
                    isLoading={resourceDelete?.isPending}
                    customCTA={
                      <IconButton
                        onClick={(e) => {
                          e.stopPropagation();
                          setDeleteItem(group);
                        }}
                      >
                        <TrashIcon className="text-gray-500 h-4 w-4 hover:text-red-600 rounded" />
                      </IconButton>
                    }
                  />
                  {/*<ExampleUsage group={group} resourceDelete={resourceDelete} />*/}
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
