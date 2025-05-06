import { Link } from 'react-router-dom';

const CustomButton = ({
  variant = 'primary',
  children,
  className: customClasses,
  href,
  size = 'md',
  loading = false,
  icon,
  startEndornment,
  endEndornment,
  fullWidth = false,
  disabled = false,
  ...props
}) => {
  const variantClasses = {
    primary: `inline-flex items-center ${children?.length > 0 ? 'gap-x-1.5' : ''} rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600`,
    'primary-outline': `inline-flex items-center  ${children?.length > 0 ? 'gap-x-1.5' : ''} rounded-md bg-transparent px-3 py-2 text-sm font-semibold text-primary-600 border border-primary-500 hover:border-primary-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600`,
    secondary: `inline-flex items-center ${children?.length > 0 ? 'gap-x-1.5' : ''} rounded-md bg-black/5 px-3 py-2 text-sm font-semibold text-secondary-500 shadow-sm hover:bg-black/10`,
    'secondary-outline': `inline-flex items-center ${children?.length > 0 ? 'gap-x-1.5' : ''}  rounded-md bg-white px-3 py-2 text-sm font-semibold text-secondary-600 shadow-sm hover:bg-secondary-50 border border-secondary-200`,
    danger: `inline-flex items-center ${children?.length > 0 ? 'gap-x-1.5' : ''} rounded-md bg-red-500 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-400`,
  };
  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-2 text-sm',
    lg: 'px-4 py-3 text-base',
  };
  const widthClass = fullWidth ? 'w-full' : 'min-w-fit';

  // const justifyClass = fullWidth ? 'justify-center' : '';
  const className = `${children ? 'min-w-24 ' : ''} flex justify-center items-center ${variantClasses[variant] || variantClasses.primary} ${sizeClasses[size]} ${widthClass}  ${
    (disabled || loading) && 'opacity-60 hover:bg-primary-700 pointer-events-none'
  } ${loading ? 'relative' : ''} ${customClasses || ''}`;

  const spinnerColor =
    variant === 'primary' || variant === 'danger'
      ? 'border-white/80 border-t-transparent'
      : 'border-secondary-500/60 border-t-transparent';

  const content = (
    <>
      {loading ? (
        <div className={'flex items-center justify-center'}>
          <span className="opacity-0 flex items-center justify-center">{children}</span>
          <span className="absolute inset-0 flex items-center justify-center">
            <div
              style={{ borderTopColor: 'transparent' }}
              className={`w-4 h-4 border-2 rounded-full animate-spin ${spinnerColor}`}
            ></div>
          </span>
        </div>
      ) : (
        <>
          {icon && startEndornment && <span className="inline-flex">{icon}</span>}
          <span className={'flex justify-center items-center gap-x-2'}>{children}</span>
          {icon && endEndornment && <span className="inline-flex">{icon}</span>}
        </>
      )}
    </>
  );

  if (href) {
    return (
      <Link
        href={href}
        {...props}
        className={`${className} ${children ? 'min-w-24 flex justify-center items-center' : ''}`}
      >
        {content}
      </Link>
    );
  }

  return (
    <button
      className={`
    ${className} 
    ${children ? 'min-w-24 flex justify-center items-center' : ''} 
    ${!disabled && !loading ? 'hover:cursor-pointer' : 'cursor-not-allowed'}
  `}
      type="button"
      disabled={disabled || loading}
      {...props}
    >
      {content}
    </button>
  );
};

export default CustomButton;
