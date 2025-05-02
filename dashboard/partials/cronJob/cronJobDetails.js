import React from 'react';
import { calculateAge } from '../../utlis/helperFunc';
import useCronJobDetails from '../../hooks/cronJob/getCronJobDetails';
import ErrorComponent from '../../components/ErrorComponent';
import CustomLinearProgress from '../../components/Loaders/LinearLoader';

const CronJobDetails = ({ formData }) => {
  const { cronJob, loading, error } = useCronJobDetails(formData);

  return (
    <>
      <CustomLinearProgress isLoading={loading} classNames={{ root: 'my-2' }} />
      {!loading && !error && (
        <div className="px-6">
          <div className="sm:px-0 ">
            <h3 className="text-base/7 font-semibold text-gray-900">Properties</h3>
            <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
              This section provides key CronJob details, including creation time, namespace, name,
              schedule, labels, annotations, active instances, suspension status, and last scheduled
              time.
            </p>
          </div>
          <div className="mt-6 border-t border-b border-gray-200">
            <dl className="divide-y divide-gray-200">
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Created</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  {calculateAge(cronJob?.data?.metadata?.creationTimestamp)} ago
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Name</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  {cronJob?.data?.metadata?.name}
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Namespace</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  {cronJob?.data?.metadata?.namespace}
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Labels</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">
                    {cronJob?.data?.metadata?.labels &&
                      Object?.entries(cronJob?.data?.metadata?.labels).map(([key, value]) => (
                        <span
                          key={key}
                          className="inline-block bg-primary-50 text-primary-700 text-xs font-medium px-2 py-1 rounded-full border border-primary-300"
                        >
                          {`${key}=${value}`}
                        </span>
                      ))}
                  </div>
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Annotations</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">
                    {cronJob?.data?.metadata?.annotations &&
                      Object?.entries(cronJob?.data?.metadata?.annotations).map(([key, value]) => (
                        <span
                          key={key}
                          className="inline-block bg-primary-50 text-primary-700 text-xs font-medium px-2 py-1 rounded-full border border-primary-300"
                        >
                          {`${key}=${value}`}
                        </span>
                      ))}
                  </div>
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Schedule</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">{cronJob?.data?.spec?.schedule}</div>
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Active</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">
                    {cronJob?.data?.status?.active?.length}
                  </div>
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Suspend</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">
                    {cronJob?.data?.spec?.suspend.toString()}
                  </div>
                </dd>
              </div>
              <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                <dt className="text-sm/6 font-medium text-gray-900">Last schedule</dt>
                <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                  <div className="flex flex-wrap gap-2">
                    {calculateAge(cronJob?.data?.status?.lastScheduleTime)}
                  </div>
                </dd>
              </div>
            </dl>
          </div>
        </div>
      )}
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default CronJobDetails;
