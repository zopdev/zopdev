import React from 'react';
import { formatTime } from '../../utlis/helperFunc';
import { PROVIDER_ICON_MAPPER } from '../../constant';

const CloudAccountCard = ({ item }) => {
  return (
    <div
      className={`break-words rounded-md transition duration-200 ease-in-out px-6 py-5 w-[330px] h-[166px] border border-gray-200`}
    >
      <div className="flex items-center gap-2 mb-5">
        {PROVIDER_ICON_MAPPER?.[item?.provider]}
        <span className="text-base font-medium text-gray-600">{item?.name}</span>
      </div>

      <div className="flex items-center justify-between mb-5">
        <span className="text-sm text-gray-500">{item?.providerId}</span>
      </div>

      <div className="text-xs text-gray-400 ">
        Updated At <span className="text-xs text-gray-500">{formatTime(item?.updatedAt)}</span>
      </div>
    </div>
  );
};

export default CloudAccountCard;
