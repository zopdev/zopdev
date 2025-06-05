const IconButton = ({ disabled, onClick, children, style, ...rest }) => {
  return (
    <button
      type="button"
      className={`inline-flex items-center justify-center p-2 rounded-full cursor-pointer group ${
        disabled ? 'opacity-50' : 'hover:bg-secondary-200 focus:outline-none'
      } ${style}`}
      onClick={onClick}
      disabled={disabled}
      {...rest}
    >
      {children}
    </button>
  );
};

export default IconButton;
