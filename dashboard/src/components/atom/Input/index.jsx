import { useEffect, useState } from 'react';

const Input = ({
                   error,
                   helperText,
                   helperTextClass,
                   errorTextClass,
                   className,
                   type,
                   value,
                   name,
                   required,
                   inputProps,
                   endAdornment,
                   testExp,
                   errorText,
                   ...props
               }) => {
    const [internalError, setInternalError] = useState(error);
    useEffect(() => {
        setInternalError(error);
    }, [error]);

    return (
        <>
            <div className="relative flex items-center">
                <input
                    type={type}
                    value={value}
                    name={name}
                    required={required}
                    className={`
            ${
                        internalError
                            ? 'px-4 block w-full rounded-md border-0 py-1.5 text-red-500 ring-1 ring-inset ring-red-300 placeholder:text-red-300 focus:ring-2 focus:ring-inset focus:ring-red-500 sm:text-sm sm:leading-6 bg-white/5 '
                            : 'px-4 block w-full rounded-md border-0 bg-transparent py-1.5 text-secondary-600 shadow-sm ring-1 ring-inset ring-secondary-300 focus:ring-2 focus:ring-inset focus:ring-primary-500 sm:text-sm sm:leading-6 '
                    }
            ${className} ${props.disabled && ' opacity-50'}
          `}
                    {...props}
                    {...inputProps}
                    onChange={(e) => {
                        if (testExp && !new RegExp(testExp).test(e.target.value)) {
                            setInternalError(true);
                        } else {
                            setInternalError(error);
                        }
                        if (props.onChange) props.onChange(e);
                    }}
                />
                {endAdornment && <div className="absolute right-3 text-secondary-400">{endAdornment}</div>}
            </div>
            {helperText != null && (
                <p className={` text-xs text-red-500 ${helperTextClass}`}>{helperText} &nbsp; </p>
            )}
            {internalError && errorText && (
                <p className={` text-xs text-red-500 ${errorTextClass}`}>{errorText} &nbsp; </p>
            )}
        </>
    );
};

export default Input;
