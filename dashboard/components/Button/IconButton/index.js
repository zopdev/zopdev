import React from 'react';

const IconButton = ({ disabled, onClick, children, style, ...rest }) => {
  return (
    <button
      type="button"
      className={`inline-flex items-center justify-center p-2 rounded-full ${
        disabled ? 'opacity-50' : 'hover:bg-gray-200 focus:outline-none focus:bg-gray-200'
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
