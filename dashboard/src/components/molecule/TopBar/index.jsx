'use client';

import { Fragment, useEffect, useState } from 'react';
import { Disclosure, Transition, Dialog } from '@headlessui/react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline/index.js';
import AppLogo from '@/assets/svg/AppLogo.jsx';
import CloudAccountSVG from '@/assets/svg/CloudAccountSvg.jsx';
import AppLogoWithText from '../../../assets/images/applogoWithText.svg'; // adjust path based on your component location

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export function TopBar() {
  const location = useLocation();
  const [tab, setTab] = useState();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  useEffect(() => {
    let DashboardView = navigation
      .map(function (ele) {
        return `/${ele.routeName.split('/')[1]}`;
      })
      .indexOf(`/${location.pathname.split('/')[1]}`); // Use location.pathname here

    DashboardView = DashboardView !== -1 ? DashboardView || 0 : -1;

    setTab(DashboardView);
  }, [location.pathname]); // Add location.pathname to dependency array

  const navigation = [
    {
      name: 'Dashboard',
      routeName: '/',
      hover: true,
      selected: true,
      isLoading: false,
      icon: CloudAccountSVG,
    },
  ];

  const navigate = useNavigate();
  const handleRedirect = (path) => {
    navigate(path);
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
                    <Link to={'/'}>
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
                      <Link key={idx} to={`${item?.routeName}`}>
                        <div
                          className={classNames(
                            tab === idx && item?.selected
                              ? ' border-primary-500  text-gray-900'
                              : ' border-transparent  text-gray-500 hover:border-gray-300 hover:text-gray-700',
                            ' inline-flex h-full items-center border-b-2 px-1 pt-1 text-sm font-medium',
                          )}
                        >
                          <span>{item?.name}</span>
                        </div>
                      </Link>
                    ))}
                  </div>
                </div>
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
                    <img width={134} height={42} src={AppLogoWithText} alt="App Logo" />
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

export default TopBar;
