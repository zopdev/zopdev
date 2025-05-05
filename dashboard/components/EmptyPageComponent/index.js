import React from 'react';
import Link from 'next/link';
import { PlusCircleIcon } from '@heroicons/react/20/solid';

const EmptyComponent = ({ imageComponent, redirectLink, buttonTitle, title, buttonIcon }) => {
  return (
    <div className="h-[70vh] w-full flex justify-center items-center">
      <div className="flex flex-col gap-4 items-center">
        <div className="flex flex-col items-center">
          {imageComponent}
          <p className="text-gray-400 text-base font-medium text-wrap mt-2">{title}</p>
        </div>
        {redirectLink && (
          <Link href={redirectLink}>
            <button
              type="button"
              className={`inline-flex items-center gap-x-1.5 rounded-md bg-primary-600    px-3 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none  hover:text-white`}
            >
              {buttonIcon || <PlusCircleIcon className="h-5 w-5" aria-hidden="true" />}
              {buttonTitle}
            </button>
          </Link>
        )}
      </div>
    </div>
  );
};

export default EmptyComponent;
