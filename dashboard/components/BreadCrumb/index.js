import { ChevronRightIcon } from '@heroicons/react/20/solid';
import Link from 'next/link';

export default function BreadCrumbComp({ breadcrumbList, style }) {
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
                    href={page.link}
                    className={`text-xs font-medium ${
                      page.disable
                        ? 'border-gray-500 border-solid cursor-default text-gray-500 '
                        : 'text-gray-500 hover:text-gray-700'
                    }`}
                  >
                    {page.name?.length > 16
                      ? `${page.name.slice(0, 10)}...${page.name.slice(-3)}`
                      : page.name}
                  </Link>
                  {breadcrumbList.length !== idx + 1 && (
                    <ChevronRightIcon
                      className=" ml-2 h-3 w-3 flex-shrink-0 text-gray-400"
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
