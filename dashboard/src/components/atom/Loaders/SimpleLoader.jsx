import React from 'react';

export default function SimpleLoader({
  size = 34,
  thickness = 2,
  color,
  secondaryCol,
  secondaryColStr = 'text-secondary-200',
  colorStr = 'text-primary-500',
}) {
  return (
    <div
      className={`relative flex justify-center items-center ${secondaryColStr} animate-spin`}
      data-testid="testSimpleLoader"
      style={{
        height: size,
        width: size,
      }}
    >
      <div
        className="border border-borderDefault rounded-full aspect-square"
        style={{
          border: `${thickness}px solid ${secondaryCol ?? 'currentColor'}`,
          height: size,
        }}
      >
        <span
          className={`flex rounded-full ${colorStr} absolute left-0 top-0 w-full h-full`}
          style={{
            borderTop: `${thickness}px solid ${color ?? 'currentColor'}`,
            borderLeft: `${thickness}px solid ${color ?? 'currentColor'}`,
            borderBottom: `${thickness}px solid ${color ?? 'currentColor'}`,
            borderRight: `${thickness}px solid transparent`,
          }}
        ></span>
      </div>
    </div>
  );
}
