import { ChevronRightIcon } from '@heroicons/react/20/solid';
import { Link } from 'react-router-dom';

export default function BreadCrumb({ breadcrumbList, style }) {
  return (
    <>
      {breadcrumbList[1]?.name !== 'loading...' && (
        <nav
          className={`flex overflow-x-auto max-w-[90vw] -mt-4 mb-3 ${style}`}
          aria-label="Breadcrumb"
        >
          <ol role="list" className="flex items-center space-x-2 whitespace-nowrap">
            {breadcrumbList.map((page, idx) => (
              <li key={idx}>
                <div className="flex items-center">
                  <Link
                    to={page.link}
                    className={`text-xs font-medium ${
                      page.disable
                        ? 'border-secondary-500 border-solid cursor-default text-secondary-400 '
                        : 'text-secondary-500 hover:text-secondary-600'
                    }`}
                  >
                    {page.name?.length > 16
                      ? `${page.name.slice(0, 10)}...${page.name.slice(-3)}`
                      : page.name}
                  </Link>
                  {breadcrumbList.length !== idx + 1 && (
                    <ChevronRightIcon
                      className=" ml-2 h-3 w-3 flex-shrink-0 text-secondary-400"
                      aria-hidden="true"
                    />
                  )}
                </div>
              </li>
            ))}
          </ol>
        </nav>
      )}
    </>
  );
}
