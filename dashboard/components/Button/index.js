const Button = ({ variant, children, className: customClasses, ...props }) => {
  const variantClasses = {
    primary:
      'inline-flex items-center gap-x-1.5 rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 ',
    secondary:
      'inline-flex items-center gap-x-1.5 rounded-md bg-black/5 px-3 py-2 text-sm font-semibold text-gray-500 shadow-sm hover:bg-black/10 ',
    danger:
      'inline-flex items-center gap-x-1.5 rounded-md bg-red-500 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-400 ',
  };

  const className = `${variantClasses[variant] || variantClasses.primary} ${
    props?.disabled && 'opacity-40 hover:bg-primary-700 pointer-events-none'
  }`;
  return (
    <button
      className={customClasses ? className + customClasses : className}
      type="button"
      {...props}
    >
      {children}
    </button>
  );
};

export default Button;
