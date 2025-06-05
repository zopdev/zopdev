import SwitchButton from '@/components/atom/Switch';
import { toast } from '@/components/molecules/Toast';
import { usePostResourceState } from '@/queries/cloud-resources';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

export const CloudResourceRow = (resource) => {
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
