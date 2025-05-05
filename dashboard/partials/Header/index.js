'use client';

import { Fragment, useContext, useEffect, useState } from 'react';
import { Disclosure, Transition, Dialog } from '@headlessui/react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline';
import AppLogo from '../../svg/appLogo';
import Image from 'next/image';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import CloudAccountSVG from '../../svg/cloudAccount';
import APP_LOGO_IMAGE from '../../public/applogoWithText.svg';
import ApplicationSvg from '../../svg/application';
import { useInitializeHeader } from '../../hooks/Header/addHeader';
import { AppContext } from '../../libs/context';

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export function TopBar() {
  const router = useRouter();
  const pathname = usePathname();
  useInitializeHeader();
  const [tab, setTab] = useState();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { appData } = useContext(AppContext);

  useEffect(() => {
    let DashboardView = navigation
      .map(function (ele) {
        return `/${ele.routeName.split('/')[1]}`;
      })
      .indexOf(`/${pathname.split('/')[1]}`);

    DashboardView = DashboardView !== -1 ? DashboardView || 0 : -1;

    setTab(DashboardView);
  }, [pathname]);

  //   const handleAppData = (entity, values) =>
  //     setAppData((prevValues) => ({ ...prevValues, [entity]: values }));

  //   useEffect(() => {
  //     if (providerList.isSuccess) {
  //       handleAppData("PROVIDERS_DATA", {
  //         data: providerList.data,
  //         isSuccess: true,
  //       });
  //     }
  //     if (namespaceList.isSuccess) {
  //       handleAppData("NAMESPACE_DATA", {
  //         data: namespaceList.data,
  //         isSuccess: true,
  //       });
  //     }
  //     if (observabilityClusterList) {
  //       handleAppData("OBS_CLUSTERS_DATA", {
  //         ...observabilityClusterList,
  //       });
  //     }
  //   }, [
  //     providerList.isFetching,
  //     namespaceList.isFetching,
  //     observabilityClusterList?.isFetching,
  //   ]);

  //   const { logout, refreshApi } = useAuth(
  //     LOGGED_IN_USER_INFO,
  //     BASE_URL,
  //     LOGIN_ENDPOINT,
  //     ACCESS_TOKEN_ENDPOINT,
  //     [],
  //     "",
  //     decodeJwt(userInfo?.accessToken)?.["tenant-id"]
  //   );

  //   const handleOrgReload = () => {
  //     window?.location.reload();
  //   };

  //   const handleOrgSwitch = (id) => {
  //     refreshApi("", id, handleOrgReload);
  //     if (router?.pathname.includes("/table")) {
  //       router.push("/table");
  //     } else if (router?.pathname !== "/cloud-accounts") {
  //       router.push("/cloud-accounts");
  //     }
  //   };

  //   const handleSignOut = (e) => {
  //     e.preventDefault();
  //     e.stopPropagation();
  //     deleteCookie("token");
  //     setUserInfo({});
  //     logout();
  //     location.href = "/app";
  //   };

  //   const handleRedirect = (path) => {
  //     router.push(path);
  //     setSidebarOpen(false);
  //   };

  const navigation = [
    {
      name: 'Cloud Accounts',
      routeName: '/cloud-accounts',
      hover: true,
      selected: true,
      isLoading: false,
      icon: CloudAccountSVG,
    },
    {
      name: 'Applications',
      routeName: '/applications',
      hover: true,
      selected: true,
      isLoading: false,
      icon: ApplicationSvg,
    },
    // {
    //   name: "Observability",
    //   routeName:
    //     observabilityClusterList?.data?.destinations?.length > 0 &&
    //     observabilityClusterList?.data?.sources?.length > 0
    //       ? "/observability"
    //       : "/observability/create",
    //   hover: false,
    //   selected: !observabilityClusterList?.isLoading,
    //   isLoading: observabilityClusterList?.isLoading,
    //   icon: ObservabilitySvg,
    // },
    // { name: 'Billing', routeName: '#', current: false },
  ];

  const handleRedirect = (path) => {
    router.push(path);
    setSidebarOpen(false);
  };

  return (
    <>
      <Disclosure as="nav" className="bg-white shadow">
        {() => (
          <>
            <div className="mx-auto max-w-8xl px-4 sm:px-6 lg:px-8">
              <div className="flex h-16 justify-between">
                <div className="flex">
                  <div className="md:flex hidden flex-shrink-0 items-center">
                    <Link href={'/'}>
                      <AppLogo />
                    </Link>
                  </div>
                  <button
                    type="button"
                    onClick={() => setSidebarOpen(true)}
                    className="-m-2.5 p-2.5 text-gray-400 md:hidden"
                  >
                    <span className="sr-only">Open sidebar</span>
                    <Bars3Icon aria-hidden="true" className="h-6 w-6" />
                  </button>
                  <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                    {navigation.map((item, idx) => (
                      <Link key={idx} href={`${item?.routeName}`}>
                        <div
                          className={classNames(
                            tab === idx && item?.selected
                              ? ' border-primary-500  text-gray-900'
                              : ' border-transparent  text-gray-500 hover:border-gray-300 hover:text-gray-700',
                            ' inline-flex h-full items-center border-b-2 px-1 pt-1 text-sm font-medium',
                          )}
                          //   onMouseEnter={() => handleMouseEnter(data?.name)}
                          //   onMouseLeave={handleMouseLeave}
                          //   onClick={handleRedirect}
                        >
                          <span>{item?.name}</span>
                        </div>
                      </Link>
                      //   <HoverMenu
                      //     data={item}
                      //     key={idx}
                      //     style={classNames(
                      //       tab === idx && item?.selected
                      //         ? " border-primary-500  text-gray-900"
                      //         : " border-transparent  text-gray-500 hover:border-gray-300 hover:text-gray-700",
                      //       " inline-flex h-full items-center border-b-2 px-1 pt-1 text-sm font-medium"
                      //     )}
                      //     providerList={providerList}
                      //     namespaceList={namespaceList}
                      //     isLoading={item.isLoading}
                      //   />
                    ))}
                  </div>
                </div>
                {/* <div className="ml-6 flex items-center">
                  <Menu as="div" className="relative ml-3">
                    <div>
                      <MenuButton className="relative flex rounded-full bg-white text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2">
                        <span className="absolute -inset-1.5" />
                        <span className="sr-only">Open user menu</span>
                        {false ? (
                          <Image
                            width={44}
                            height={44}
                            className="h-8 w-8 rounded-full"
                            // src={userInfo.profile_picture}
                            src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                            alt="Profile Image"
                            referrerPolicy="no-referrer"
                          />
                        ) : (
                          <div className="relative w-10 h-10 overflow-hidden bg-gray-100 rounded-full dark:bg-gray-600">
                            <svg
                              className="absolute w-12 h-12 text-gray-400 -left-1"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                              xmlns="http://www.w3.org/2000/svg"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                                clipRule="evenodd"
                              ></path>
                            </svg>
                          </div>
                        )}
                      </MenuButton>
                    </div>
                  </Menu>
                </div> */}
              </div>
            </div>
          </>
        )}
      </Disclosure>
      <Transition.Root show={sidebarOpen} as={Fragment}>
        <Dialog as="div" className="relative z-50 md:hidden" onClose={setSidebarOpen}>
          <Transition.Child
            as={Fragment}
            enter="transition-opacity duration-300 ease-linear"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="transition-opacity duration-300 ease-linear"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-gray-900/80" />
          </Transition.Child>

          <div className="fixed inset-0 flex">
            <Transition.Child
              as={Fragment}
              enter="transform transition ease-in-out duration-300 sm:duration-300"
              enterFrom="-translate-x-full"
              enterTo="translate-x-0"
              leave="transform transition ease-in-out duration-300 sm:duration-300"
              leaveFrom="translate-x-0"
              leaveTo="-translate-x-full"
            >
              <Dialog.Panel className="relative mr-16 flex w-full max-w-xs flex-1 transform transition duration-300 ease-in-out">
                <div className="absolute left-full top-0 flex w-16 justify-center pt-5 duration-300 ease-in-out">
                  <button
                    type="button"
                    onClick={() => setSidebarOpen(false)}
                    className="-m-2.5 p-2.5"
                  >
                    <span className="sr-only">Close sidebar</span>
                    <XMarkIcon aria-hidden="true" className="h-6 w-6 text-white" />
                  </button>
                </div>
                <div className="flex grow flex-col gap-y-8 overflow-y-auto bg-white px-6 pb-2  ring-1 ring-white/10">
                  <div className="flex h-16 shrink-0 items-center ">
                    <Image width={134} height={42} src={APP_LOGO_IMAGE} alt="Profile Image" />
                  </div>
                  <nav className="flex flex-1 flex-col">
                    <ul role="list" className="flex flex-1 flex-col gap-y-7">
                      <li>
                        <ul role="list" className="-mx-2 space-y-4">
                          {navigation.map((item, idx) => (
                            <li key={idx}>
                              <div
                                onClick={() => handleRedirect(item.routeName)}
                                className={classNames(
                                  tab === idx && item?.selected
                                    ? 'bg-gray-50 text-primary-500'
                                    : 'text-gray-600 hover:bg-gray-50',
                                  'group flex gap-x-3 rounded-md p-2 text-sm font-medium leading-6 cursor-pointer',
                                )}
                              >
                                <item.icon
                                  aria-hidden="true"
                                  className={`h-6 w-6 shrink-0`}
                                  color={tab === idx ? 'rgb(6 182 212)' : '#9197B3'}
                                />
                                {item.name}
                              </div>
                            </li>
                          ))}
                        </ul>
                      </li>
                    </ul>
                  </nav>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
    </>
  );
}

export default function TopBarWrapper() {
  return (
    // <TopBarContainer {...props}>
    <TopBar />
    // </TopBarContainer>
  );
}
