import { ExclamationCircleIcon } from '@heroicons/react/20/solid';

export default function ErrorComponent({ errorText, className }) {
  return (
    <div className={`rounded-md bg-red-600 p-4 bg-opacity-20 border border-red-500 ${className} `}>
      <div className="flex">
        <div className="flex-shrink-0">
          <ExclamationCircleIcon className="h-6 w-6 text-red-500" aria-hidden="true" />
        </div>
        <div className="ml-3">
          <h3 className="text-sm font-medium text-red-600 text-opacity-70 break-all">
            {errorText}
          </h3>
        </div>
      </div>
    </div>
  );
}
