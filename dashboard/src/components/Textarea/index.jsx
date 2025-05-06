function Textarea({ error, helperText, className, ...props }) {
  return (
    <>
      <textarea
        className={`
          ${
            error
              ? 'px-4 block w-full rounded-md border-0 py-1.5 pr-10 text-red-500 ring-1 ring-inset ring-red-300 placeholder:text-red-300 focus:ring-2 focus:ring-inset focus:ring-red-500 sm:text-sm sm:leading-6 bg-white/5 '
              : 'px-4 block w-full rounded-md border-0 bg-transparent py-1.5 text-secondary-600 shadow-sm ring-1 ring-inset ring-secondary-300 focus:ring-2 focus:ring-inset focus:ring-primary-500 sm:text-sm sm:leading-6 '
          }
            ${className} ${props.disabled && ' opacity-40'}
        `}
        {...props}
      />
      {helperText != null && <p className="mt-2 text-sm text-red-500">{helperText} &nbsp; </p>}
    </>
  );
}

export default Textarea;
