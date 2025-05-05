import React from 'react';
import SimpleLoader from '../SimpleLoader';

/**
 * this function is to display loader when query is fetching
 * @returns loader
 */
export default function CompleteLoader() {
  return (
    <div className="center-complete-loader">
      <SimpleLoader size={50} thickness={3} />
    </div>
  );
}
