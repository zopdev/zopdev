import { useEffect, useRef, useState } from 'react';
import { useLocation, Link } from 'react-router-dom';

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export default function Sidebar({ menu }) {
  const location = useLocation();
  const selectedTabRef = useRef(null);
  const [activeIndex, setActiveIndex] = useState(0);

  const currentMenu = menu || [];

  useEffect(() => {
    if (selectedTabRef.current) {
      selectedTabRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
        inline: 'start',
      });
    }
  }, [activeIndex]);

  useEffect(() => {
    const currentPath = location.pathname.split('/').pop();
    const index = currentMenu.findIndex((item) => item.link.endsWith(currentPath || ''));
    setActiveIndex(index >= 0 ? index : 0);
  }, [location.pathname, currentMenu]);

  return (
    <>
      {/* Mobile Sidebar */}
      <div className="border-b border-gray-200 border-opacity-5 overflow-x-auto md:hidden">
        <nav className="-mb flex items-center px-2 max-w-full" aria-label="Tabs">
          {currentMenu.map((item, index) =>
            !item.heading ? (
              <Link
                key={index}
                to={`/${item.link}`}
                className={classNames(
                  index === activeIndex
                    ? 'text-primary-500 bg-primary-selected rounded-md'
                    : 'text-gray-500 hover:text-gray-700',
                  'group flex items-center py-2 px-2 text-sm font-medium !p-[1rem]',
                )}
                ref={index === activeIndex ? selectedTabRef : null}
              >
                <item.icon className="mr-2 h-5 w-5" color="currentColor" />
                <span className="whitespace-nowrap">{item.text}</span>
              </Link>
            ) : null,
          )}
        </nav>
      </div>

      {/* Desktop Sidebar */}
      <div className="md:flex grow flex-col gap-y-5 overflow-y-auto p-8 border-t border-r border-gray-200 bg-white xs:hidden min-h-full">
        <nav className="flex flex-1 flex-col">
          <ul className="flex flex-1 flex-col gap-y-7">
            <li>
              <ul className="-mx-2 space-y-1">
                {currentMenu.map((item, index) => {
                  if (item.heading) {
                    if (item.divider) {
                      return (
                        <div key={index} className="space-y-0 !mb-6 !mt-6">
                          <div className="border-b border-b-gray-200 border-opacity-30" />
                          <div className="border-t border-t-gray-200 border-opacity-95" />
                        </div>
                      );
                    }
                    return (
                      <li key={index} className="text-left">
                        <div className="text-xs font-semibold leading-6 text-gray-400 mt-6 mb-0 p-2">
                          {item.headingTitle}
                        </div>
                      </li>
                    );
                  }

                  return (
                    <Link key={item.text} to={`/${item.link}`} className="text-left">
                      <div
                        className={classNames(
                          index === activeIndex
                            ? 'bg-gray-50 text-primary-600'
                            : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600',
                          'group flex gap-x-3 rounded-md p-2 text-sm font-medium leading-6 cursor-pointer w-48 mb-4',
                        )}
                      >
                        <item.icon
                          className={classNames(
                            index === activeIndex
                              ? 'text-primary-600'
                              : 'text-gray-400 group-hover:text-primary-600',
                            'h-6 w-6 shrink-0',
                          )}
                          color="currentColor"
                        />
                        {item.text}
                      </div>
                    </Link>
                  );
                })}
              </ul>
            </li>
          </ul>
        </nav>
      </div>
    </>
  );
}
