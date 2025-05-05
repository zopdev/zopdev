import { PlusCircleIcon } from '@heroicons/react/16/solid';
import { useEffect, useState } from 'react';

export default function InputWithButton({
  type = '',
  value,
  error = false,
  name,
  required,
  testExp,
  helperText,
  helperTextClass,
  errorTextClass,
  errorText,
  onChange,
  onClick,
  inputProps,
  className,
  ...props
}) {
  const [internalError, setInternalError] = useState(error);
  useEffect(() => {
    setInternalError(error);
  }, [error]);
  return (
    <div>
      <div className="mt-2 flex">
        <div className="-mr-px grid grow grid-cols-1 focus-within:relative">
          <input
            type={type}
            value={value}
            name={name}
            required={required}
            className={`col-start-1 row-start-1 block w-full rounded-l-md bg-white py-1.5 pr-3 text-gray-600 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600 sm:pl-4 sm:text-sm ${internalError && ' text-red-500 ring-1 ring-inset ring-red-300 placeholder:text-red-300'} ${className} ${props.disabled && ' opacity-50'}`}
            {...props}
            {...inputProps}
            onChange={(e) => {
              if (testExp && !new RegExp(testExp).test(e.target.value)) {
                setInternalError(true);
              } else {
                setInternalError(error);
              }
              if (onChange) onChange(e);
            }}
          />
          {/* <UsersIcon
            aria-hidden="true"
            className="pointer-events-none col-start-1 row-start-1 ml-3 size-5 self-center text-gray-400 sm:size-4"
          /> */}
        </div>
        <button
          type="button"
          className="flex shrink-0 items-center gap-x-1.5 rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-600 outline outline-1 -outline-offset-1 outline-gray-300 hover:bg-gray-50 focus:relative focus:outline focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600"
          onClick={onClick}
        >
          <PlusCircleIcon aria-hidden="true" className="-ml-0.5 size-4 text-gray-400" />
        </button>
      </div>
      {helperText != null && (
        <p className={` text-xs text-red-500 ${helperTextClass}`}>{helperText} &nbsp; </p>
      )}
      {internalError && errorText && (
        <p className={` text-xs text-red-500 ${errorTextClass}`}>{errorText} &nbsp; </p>
      )}
    </div>
  );
}
