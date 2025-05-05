import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { formatTime } from '../../utlis/helperFunc';

const ApplicationCard = ({ data, view }) => {
  const router = useRouter();

  useEffect(() => {
    if (data?.id) {
      router.prefetch(`/applications/${data?.id}`);
    }
  }, [data?.id, router]);

  const handleRedirect = () => {
    router.push(`/applications/${data?.id}/environment`);
  };

  return (
    <div
      className={`relative w-[330px] ${view === 'navBar' ? 'h-[146px]' : 'h-auto'} p-6 ${
        view === 'navBar' ? 'px-6 py-4' : 'px-5 py-5'
      } cursor-pointer rounded-md border border-gray-200 hover:bg-gray-50 hover:shadow-lg transition-shadow`}
      onClick={handleRedirect}
    >
      <div className="flex items-center mb-5 gap-2">
        <h3 className="text-lg font-medium text-gray-600">{data?.name}</h3>
      </div>

      <div className="flex justify-between mb-5">
        <div className="flex w-full  items-center">
          {data?.services !== 0 && (
            <p className="text-xs text-gray-500 font-medium">{data?.environments?.length}</p>
          )}
          <p className="text-xs text-gray-500 font-light ">
            &nbsp;
            {data?.environments?.length === 1
              ? 'Environment'
              : data?.environments?.length > 1
                ? 'Environments'
                : ' '}
          </p>
        </div>
      </div>

      {data?.environments?.length > 1 ? (
        <div className={`flex gap-2 overflow-auto scroll-hidden`}>
          {data.environments.map((single) => (
            <div
              className="h-14 bg-gray-900/5 rounded-lg flex justify-center items-center p-2 mb-3 min-w-[100] w-full"
              key={single.id}
            >
              <div>
                {/* <p className="text-xs text-gray-500">Environment</p> */}
                <p className="text-sm text-gray-900 text-start">{single.name}</p>
              </div>

              {/* <div>
                <div className={`relative flex gap-1`}>
                  <p className="text-xs text-gray-500">Order</p>
                </div>
                <p className="text-sm text-gray-900 text-right">{single.level ?? '1'}</p>
              </div> */}
            </div>
          ))}
        </div>
      ) : (
        <div className="h-14 bg-gray-900/5 rounded-lg flex justify-between p-2 mb-3">
          <div>
            <p className="text-xs text-gray-500">Environment</p>
            <p className="text-sm text-gray-900 text-start">{data?.environments?.[0]?.name}</p>
          </div>

          <div>
            <p className="text-xs text-gray-500">Order</p>
            <p className="text-sm text-gray-900 text-center">{data?.environments?.[0]?.level}</p>
          </div>
        </div>
      )}

      <div className="text-xs text-gray-400 pt-1 ">
        Updated At{' '}
        {/* <span className="text-gray-900">
          {item?.updatedByIdentifier || item?.createdByIdentifier}
        </span> */}
        <span className="text-xs text-gray-500">{formatTime(data?.updatedAt)}</span>
      </div>
      {/* <p className="text-xs text-gray-500">
        <span>{formatTime(data?.updatedAt)}</span>
      </p> */}
    </div>
  );
};

export default ApplicationCard;
