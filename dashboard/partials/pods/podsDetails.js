import React from 'react';
import usePodDetails from '../../hooks/pods/getPodDetails';
import CustomLinearProgress from '../../components/Loaders/LinearLoader';
import ErrorComponent from '../../components/ErrorComponent';

const PodsDetails = ({ formData }) => {
  const { pod, loading, error } = usePodDetails(formData);
  return (
    <>
      <CustomLinearProgress isLoading={loading} classNames={{ root: 'my-2' }} />
      {!loading && !error && (
        <div className="px-6">
          <>
            <div className="sm:px-0 ">
              <h3 className="text-base/7 font-semibold text-gray-900">Properties</h3>
              <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
                The Properties section offers a concise overview of the Pod's metadata,
                configuration, and status, including creation details, labels, annotations, and
                operational insights for easy management.
              </p>
            </div>
            <div className="mt-6 border-t border-b border-gray-200">
              <dl className="divide-y divide-gray-200">
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Created</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.metadata?.creationTimestamp}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Name</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.metadata?.name}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Namespace</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.metadata?.namespace}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Labels</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {pod?.data?.metadata?.labels &&
                        Object?.entries(pod?.data?.metadata?.labels).map(([key, value]) => (
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
                      {pod?.data?.metadata?.annotations &&
                        Object?.entries(pod?.data?.metadata?.annotations).map(([key, value]) => (
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
                  <dt className="text-sm/6 font-medium text-gray-900">Controlled By</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.metadata?.ownerReferences?.[0]?.kind}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Node</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.spec?.nodeName}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Pod IP</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.status?.podIP}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Service Account</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    {pod?.data?.spec?.serviceAccountName}
                  </dd>
                </div>
                <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                  <dt className="text-sm/6 font-medium text-gray-900">Conditions</dt>
                  <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                    <div className="flex flex-wrap gap-2">
                      {pod?.data?.status?.conditions.map((container, index) => (
                        <div
                          key={index}
                          className={`inline-block bg-primary-50 text-primary-700 text-xs font-medium px-2 py-1 rounded-full border border-primary-300`}
                        >
                          {container.type}
                        </div>
                      ))}
                    </div>
                  </dd>
                </div>
              </dl>
            </div>
          </>
          <div className="mt-10 ">
            <div className="sm:px-0 pb-4">
              <h3 className="text-base/7 font-semibold  text-gray-900">Containers</h3>
              <p className="mt-1 max-w-2xl text-sm/6 text-gray-500">
                This section provides detailed insights into each container's configuration,
                including runtime settings, resource allocations, environment variables, ports, and
                volume mounts.
              </p>
            </div>

            {pod?.data?.spec?.containers?.map((container, index) => {
              return (
                <div key={index} className="mt-4 border-t py-4">
                  <div className="flex gap-2 items-center">
                    <div
                      className={`w-4 h-4 ${container.ready ? 'bg-green-500' : 'bg-orange-500'} rounded-full group-hover`}
                    ></div>
                    <h3 className="text-base/7 font-semibold text-gray-900 ">{container.name}</h3>
                  </div>
                  <div className="mt-6 border-t border-gray-200">
                    <dl className="divide-y divide-gray-200">
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Status</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          {pod?.data?.status?.containerStatuses &&
                            pod?.data?.status?.containerStatuses?.map((item, containerIndex) => {
                              const stateKeys = Object?.keys(item.state).join(', ');
                              return <span key={containerIndex}>{stateKeys}</span>;
                            })}
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Image</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          {container?.image}
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">ImagePullPolicy</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          {container?.imagePullPolicy}
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Ports</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            {container?.ports?.map((item, index) => (
                              <span
                                key={index}
                                className="inline-block bg-green-50 text-green-700 text-xs font-medium px-2 py-1 rounded-full border border-green-300"
                              >
                                {`${item.name}:${item.containerPort}/${item.protocol}`}
                              </span>
                            ))}
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Environments</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            <div className="flex gap-2 flex-wrap">
                              {container?.env.map((env, index) => (
                                <div
                                  key={index}
                                  className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300"
                                >
                                  <strong>{env.name}:</strong> {env.value || 'N/A'}
                                </div>
                              ))}
                            </div>
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Mounts</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            <div className="flex gap-2 flex-wrap">
                              {container?.volumeMounts.map((mount, index) => (
                                <div
                                  key={index}
                                  className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300"
                                >
                                  <strong>{mount.name}:</strong> {mount.mountPath || 'N/A'}
                                </div>
                              ))}
                            </div>
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Liveness</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            <div className="flex gap-2 flex-wrap">
                              {container?.livenessProbe?.httpGet?.path && (
                                <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                  {container?.livenessProbe?.httpGet?.path}
                                </div>
                              )}
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                delay={container?.livenessProbe?.initialDelaySeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                timeout={container?.livenessProbe?.timeoutSeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                period={container?.livenessProbe?.periodSeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                #success={container?.livenessProbe?.successThreshold}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                #failure={container?.livenessProbe?.failureThreshold}
                              </div>
                            </div>
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Readiness</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            <div className="flex gap-2 flex-wrap">
                              {container?.readinessProbe?.httpGet?.path && (
                                <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                  {container?.readinessProbe?.httpGet?.path}
                                </div>
                              )}
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                delay={container?.readinessProbe?.initialDelaySeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                timeout={container?.readinessProbe?.timeoutSeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                period={container?.readinessProbe?.periodSeconds}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                #success={container?.readinessProbe?.successThreshold}
                              </div>
                              <div className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300">
                                #failure={container?.readinessProbe?.failureThreshold}
                              </div>
                            </div>
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Arguments</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          <div className="flex flex-wrap gap-2">
                            <div className="flex gap-2 flex-wrap">
                              {container?.args &&
                                container?.args.map((arg, index) => (
                                  <div
                                    key={index}
                                    className="px-3 py-1 bg-gray-100 text-gray-700 rounded-md border border-gray-300"
                                  >
                                    {arg}
                                  </div>
                                ))}
                            </div>
                          </div>
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Requests</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          CPU - {container?.resources?.requests?.cpu} , Memory -{' '}
                          {container?.resources?.requests?.memory}
                        </dd>
                      </div>
                      <div className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt className="text-sm/6 font-medium text-gray-900">Limits</dt>
                        <dd className="mt-1 text-sm/6 text-gray-700 sm:col-span-2 sm:mt-0">
                          CPU - {container?.resources?.limits?.cpu} , Memory -{' '}
                          {container?.resources?.limits?.memory}
                        </dd>
                      </div>
                    </dl>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
      {error && (
        <ErrorComponent errorText={error || 'Something went wrong !'} className={' !p-2'} />
      )}
    </>
  );
};

export default PodsDetails;
