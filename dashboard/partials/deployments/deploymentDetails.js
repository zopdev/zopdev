import React from 'react';
import useDeploymentDetails from '../../hooks/deployment/getDeploymentDetails';
import { calculateAge } from '../../utlis/helperFunc';
import CustomLinearProgress from '../../components/Loaders/LinearLoader';
import ErrorComponent from '../../components/ErrorComponent';

const StatusColorCode = {
  AVAILABLE: 'text-green-700 bg-green-50 ring-green-300',
  PROGRESSING: 'text-primary-700 bg-primary-50 ring-primary-300',
};

const DeploymentDetails = ({ formData }) => {
  const { deployment, loading, error } = useDeploymentDetails(formData);
  return (
    <>
      <CustomLinearProgress isLoading={loading} classNames={{ root: 'my-2' }} />
      {!loading && !error && (
        <div className="px-6">
          <>
            <div className="sm:px-0 ">
              <h3 className="text-base/7 font-semibold text-gray-900">Properties</h3>
              <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
                The Properties section summarizes the deployment's configuration and state,
                including creation time, name, namespace, labels, annotations, replicas, selectors,
                strategy type, and status conditions.
              </p>
            </div>
            <div className="mt-6 border-t border-b border-gray-200">
              <dl className="divide-y divide-gray-200">
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Created</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {calculateAge(deployment?.data?.metadata?.creationTimestamp)} ago
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Name</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {deployment?.data?.metadata?.name}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Namespace</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {deployment?.data?.metadata?.namespace}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Labels</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {deployment?.data?.metadata?.labels &&
                        Object?.entries(deployment?.data?.metadata?.labels).map(([key, value]) => (
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
                      {deployment?.data?.metadata?.annotations &&
                        Object?.entries(deployment?.data?.metadata?.annotations).map(
                          ([key, value]) => (
                            <span
                              key={key}
                              className="inline-block bg-primary-50 text-primary-700 text-xs font-medium px-2 py-1 rounded-full border border-primary-300"
                            >
                              {`${key}=${value}`}
                            </span>
                          ),
                        )}
                    </div>
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Replicas</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {deployment?.data?.status?.replicas} desired,{' '}
                      {deployment?.data?.status?.availableReplicas} available,{' '}
                      {deployment?.data?.status?.replicas -
                        deployment?.data?.status?.availableReplicas}{' '}
                      unavailable
                    </div>
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Selectors</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {deployment?.data?.spec?.selector?.matchLabels &&
                        Object?.entries(deployment?.data?.spec?.selector?.matchLabels).map(
                          ([key, value]) => (
                            <span
                              key={key}
                              className="inline-block bg-primary-50 text-primary-700 text-xs font-medium px-2 py-1 rounded-full border border-primary-300"
                            >
                              {`${key}=${value}`}
                            </span>
                          ),
                        )}
                    </div>
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Strategy Type</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {deployment?.data?.spec?.strategy?.type}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Conditions</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {deployment?.data?.status?.conditions &&
                        Object?.entries(deployment?.data?.status?.conditions).map(
                          ([key, value]) => (
                            <span
                              key={key}
                              className={`${
                                StatusColorCode[value?.type.toUpperCase()]
                              } inline-block rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset`}
                            >
                              {`${value?.type}`}
                            </span>
                          ),
                        )}
                    </div>
                  </dd>
                </div>
              </dl>
            </div>
          </>
        </div>
      )}
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default DeploymentDetails;
