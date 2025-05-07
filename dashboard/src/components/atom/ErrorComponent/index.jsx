import { ExclamationCircleIcon } from '@heroicons/react/20/solid';

export default function ErrorComponent({ errorText, className, complete = false }) {
  return (
    <>
      {!complete && (
        <div className={`rounded-md bg-red-600/20 p-4 border border-red-500 ${className} `}>
          <div className="flex">
            <div className="flex-shrink-0">
              <ExclamationCircleIcon className="h-6 w-6 text-red-500" aria-hidden="true" />
            </div>
            <div className="ml-3 min-w-0">
              <h3 className="text-sm font-medium text-red-600/70 break-words">
                {errorText || 'Something went wrong'}
              </h3>
            </div>
          </div>
        </div>
      )}
      {complete && (
        <div className={`rounded-md bg-red-600/20 p-4 border border-red-500 ${className} `}>
          <div className="flex justify-center items-center flex-col gap-1">
            <div className="flex-shrink-0 flex justify-center items-center gap-x-2">
              <span className={'font-bold text-red-500'}> Oops!</span>
            </div>
            <div className="min-w-0">
              <h3 className="text-sm font-medium text-red-600/70 break-words">
                {errorText || 'Something went wrong'}
              </h3>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
