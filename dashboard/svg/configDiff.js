import React from 'react';

export const ConfigDiff = ({ color = '#9197B3', ...rest }) => {
  return (
    <svg
      width="22px"
      height="22px"
      viewBox="0 0 22 22"
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      {...rest}
    >
      <g id="New-Layout-" stroke="none" strokeWidth={1} fill="none" fillRule="evenodd">
        <g id="Application--Summary-page" transform="translate(-30, -683)">
          <g id="Metrics" transform="translate(30, 683)">
            <rect id="Rectangle" x={0} y={0} width={22} height={22} />
            <rect
              id="Rectangle"
              stroke={color}
              strokeWidth="1.5"
              x={1}
              y={1}
              width={20}
              height={20}
              rx={5}
            />
            <polyline
              id="Path"
              stroke={color}
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              points="5 12 7 12 9 15 13 8 15 12 17 12"
            />
          </g>
        </g>
      </g>
    </svg>
  );
};
