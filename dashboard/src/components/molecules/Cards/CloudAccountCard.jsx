import { formatTime } from '@/utils/helper';
import { PROVIDER_ICON_MAPPER } from '@/utils/componentMapper';
import { Link } from 'react-router-dom';

const CloudAccountCard = ({ item }) => {
  const isGCP = item?.provider === 'gcp';
  const CardWrapper = isGCP ? Link : 'div';
  const wrapperProps = isGCP ? { to: `/cloud-accounts/${item?.id}/resources` } : {};

  return (
    <CardWrapper {...wrapperProps}>
      <div
        className={`break-words rounded-md px-6 py-5 w-[330px] h-[166px] border border-gray-200
          duration-200 ease-in-out
          ${isGCP ? 'hover:bg-gray-50 hover:shadow-lg cursor-pointer transition-shadow' : ''}
        `}
      >
        <div className="flex items-center gap-2 mb-5">
          {PROVIDER_ICON_MAPPER?.[item?.provider]}
          <span className="text-base font-medium text-gray-600">{item?.name}</span>
        </div>

        <div className="flex items-center justify-between mb-5">
          <span className="text-sm text-gray-500">{item?.providerId}</span>
        </div>

        <div className="text-xs text-gray-400">
          Updated At <span className="text-xs text-gray-500">{formatTime(item?.updatedAt)}</span>
        </div>
      </div>
    </CardWrapper>
  );
};

export default CloudAccountCard;
