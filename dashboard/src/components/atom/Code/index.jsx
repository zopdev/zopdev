import React from 'react';

const Code = ({ children, className, ...props }) => {
  return (
    <code className={`block px-1.5 py-0.5 rounded ${className}`} {...props}>
      {children}
    </code>
  );
};

export default Code;
