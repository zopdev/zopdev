import React from 'react';
import useServiceDetails from '../../hooks/service/getServiceDetails';
import CustomLinearProgress from '../../components/Loaders/LinearLoader';
import ErrorComponent from '../../components/ErrorComponent';
import { calculateAge } from '../../utlis/helperFunc';

const ServiceDetails = ({ formData }) => {
  const { service, loading, error } = useServiceDetails(formData);
  return (
    <>
      <CustomLinearProgress isLoading={loading} classNames={{ root: 'my-2' }} />
      {!loading && !error && (
        <div className="px-6">
          <>
            <div className="sm:px-0 ">
              <h3 className="text-base/7 font-semibold text-gray-900">Properties</h3>
              <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
                The Properties section provides key details about the service, including its
                creation time, name, namespace, labels, selectors, type, and session affinity. It
                offers a clear overview of the service's configuration, role, and behavior, aiding
                in monitoring and management.
              </p>
            </div>
            <div className="mt-6 border-t border-b border-gray-200">
              <dl className="divide-y divide-gray-200">
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Created</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {calculateAge(service?.data?.metadata?.creationTimestamp)}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Name</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {service?.data?.metadata?.name}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Namespace</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {service?.data?.metadata?.namespace}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Labels</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {service?.data?.metadata?.labels &&
                        Object?.entries(service?.data?.metadata?.labels).map(([key, value]) => (
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
                  <dt className="text-sm/6 font-medium text-gray-900">Selectors</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {service?.data?.spec?.selector &&
                        Object?.entries(service?.data?.spec?.selector).map(([key, value]) => (
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
                  <dt className="text-sm/6 font-medium text-gray-900">Type</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {service?.data?.spec?.type}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Session Affinity</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {service?.data?.spec?.sessionAffinity}
                  </dd>
                </div>
              </dl>
            </div>
          </>
          <div className="mt-10 ">
            <div className="sm:px-0">
              <h3 className="text-base/7 font-semibold text-gray-900">Connection</h3>
              <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
                The Connection section details the service's network setup, including its Cluster IP
                for internal routing and port mappings with target ports and protocols. It provides
                a clear view of the service's accessibility and network configuration.
              </p>
            </div>
            <div className="mt-6 border-t border-gray-200">
              <dl className="divide-y divide-gray-200">
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Cluster IP</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {service?.data?.spec?.clusterIP}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Ports</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {service?.data?.spec?.ports?.map((item, index) => (
                        <span
                          key={index}
                          className="inline-block bg-green-50 text-green-700 text-xs font-medium px-2 py-1 rounded-full border border-green-300"
                        >
                          {`${item.port}:${item.targetPort}/${item.protocol}`}
                        </span>
                      ))}
                    </div>
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        </div>
      )}
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default ServiceDetails;
