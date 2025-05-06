import React from 'react';
import IconButton from '../Button/IconButton/index.jsx';
import { colourCode, REFRESH_STATUS } from '@/constant.jsx';
import DotWithProgress from '@/components/atom/Loaders/DotWithProgress.jsx';
import Tooltip from '@/components/atom/Tooltip/index.jsx';
import RefreshIcon from '@/assets/svg/RefreshIcon.jsx';

const ResourceStatus = ({
  status,
  retry = false,
  retryTooltip = 'Retry',
  id,
  handleRetry,
  isRetrying,
}) => {
  const isRefreshing = REFRESH_STATUS.includes(status);
  const color = colourCode[status] || '#ccc';

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
      <Tooltip title={`Status: ${status}`} className={{ root: 'flex' }} arrow>
        {isRefreshing ? (
          <DotWithProgress color={color} />
        ) : (
          <span
            className={`status-rounded`}
            style={{
              backgroundColor: color,
            }}
          />
        )}
      </Tooltip>
      {retry && (
        <Tooltip title={retryTooltip} arrow>
          <div>
            <IconButton disabled={isRetrying} onClick={(e) => handleRetry(e, id)} size="small">
              <RefreshIcon fontSize="small" />
            </IconButton>
          </div>
        </Tooltip>
      )}
    </div>
  );
};

export default ResourceStatus;
