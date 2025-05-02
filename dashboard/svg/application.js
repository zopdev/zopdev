import React from 'react';

const ApplicationSvg = ({ color = '#9197B3' }) => {
  return (
    <svg
      width="22px"
      height="22px"
      viewBox="0 0 22 22"
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
    >
      <g id="New-Layout-" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
        <g id="Cloud-Account---" transform="translate(-534, -40)">
          <g id="Applications" transform="translate(534, 40)">
            <rect id="Rectangle" x="0" y="0" width="22" height="22"></rect>
            <rect
              id="Rectangle"
              stroke={color}
              strokeWidth="1.5"
              x="4"
              y="2"
              width="14"
              height="18"
              rx="2.5"
            ></rect>
            <circle id="Oval" fill={color} cx="11" cy="17" r="1"></circle>
            <circle id="Oval-Copy-3" fill={color} cx="7" cy="7" r="1"></circle>
            <circle id="Oval-Copy-4" fill={color} cx="7" cy="12" r="1"></circle>
            <rect id="Rectangle" fill={color} x="9" y="6" width="7" height="2" rx="1"></rect>
            <rect
              id="Rectangle-Copy-12"
              fill={color}
              x="9"
              y="11"
              width="7"
              height="2"
              rx="1"
            ></rect>
          </g>
        </g>
      </g>
    </svg>
  );
};

export default ApplicationSvg;
