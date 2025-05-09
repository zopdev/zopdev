import React from 'react';
import Button from '@/components/atom/Button/index.jsx';

const IconButton = ({ disabled, onClick, children, style, ...rest }) => {
  return (
    <Button
      type="button"
      className={`inline-flex items-center justify-center p-2 rounded-full ${
        disabled ? 'opacity-50' : 'hover:bg-secondary-200 focus:outline-none focus:bg-secondary-200'
      } ${style}`}
      onClick={onClick}
      disabled={disabled}
      {...rest}
    >
      {children}
    </Button>
  );
};

export default IconButton;
