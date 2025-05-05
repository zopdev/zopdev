import React from 'react';

export default function CustomLinearProgress({ isLoading, classNames = {} }) {
  return (
    <div
      className={`w-full bg-gray-200 pt-[2px] h-[2px] rounded-md relative overflow-hidden transition-opacity duration-500 ${
        classNames.root ?? ''
      } ${isLoading ? 'opacity-100' : 'opacity-0'}`}
      data-testid="testLinearLoader"
    >
      <span
        className={
          `absolute max-w-xl w-full bg-primary-500 h-full left-0 top-0 animate-loader ` +
            classNames.loader ?? ''
        }
      ></span>
    </div>
  );
}
