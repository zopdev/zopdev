import React from 'react';

export const Environment = ({ color = '#9197B3', ...rest }) => {
  return (
    <svg
      width="23px"
      height="23px"
      viewBox="0 0 23 23"
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      {...rest}
    >
      <g id="Updated-sidebar" stroke="none" strokeWidth={1} fill="none" fillRule="evenodd">
        <g id="Clusters" transform="translate(-30, -337)">
          <g id="Audit-Log" transform="translate(30, 337)">
            <g id="Logs">
              <rect id="Rectangle" x={0} y={0} width={22} height={22} />
              <path
                d="M11,21 L7,21 C4.79086,21 3,19.2091 3,17 L3,5 C3,2.79086 4.79086,1 7,1 L11,1 L13.0633,1 C13.6568,1 14.2197,1.26365 14.5997,1.71963 L18.5364,6.44373 C18.836,6.80316 19,7.25623 19,7.7241 L19,11"
                id="Path"
                stroke={color}
                strokeWidth="1.5"
                strokeLinecap="round"
              />
              <path
                d="M14,1.5 L14,5 C14,6.10457 14.8954,7 16,7 L18.5,7"
                id="Path"
                stroke={color}
                strokeWidth="1.5"
                strokeLinecap="round"
              />
              <line
                x1={7}
                y1={9}
                x2={13}
                y2={9}
                id="Path"
                stroke={color}
                strokeWidth="1.5"
                strokeLinecap="round"
              />
              <line
                x1={7}
                y1={14}
                x2={10}
                y2={14}
                id="Path"
                stroke={color}
                strokeWidth="1.5"
                strokeLinecap="round"
              />
            </g>
            <g
              id="Group-54"
              transform="translate(11.8174, 11.8174)"
              stroke={color}
              strokeWidth="1.5"
            >
              <circle
                id="Oval"
                transform="translate(4.5, 4.5) rotate(-43) translate(-4.5, -4.5)"
                cx="4.5"
                cy="4.5"
                r="3.18392007"
              />
              <line
                x1="7.18257686"
                y1="7.18257686"
                x2="9.18257686"
                y2="9.18257686"
                id="Path-7"
                strokeLinecap="round"
              />
            </g>
          </g>
        </g>
      </g>
    </svg>
  );
};
