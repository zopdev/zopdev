const ResourceIcon = ({ color = '#9197B3', ...props }) => {
  return (
    <svg
      width={48}
      height={48}
      viewBox="0 0 48 48"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect width={48} height={48} />
      <path
        d="M34.6667 10.6667H13.3333C11.8606 10.6667 10.6667 11.8606 10.6667 13.3333V18.6667C10.6667 20.1394 11.8606 21.3333 13.3333 21.3333H34.6667C36.1394 21.3333 37.3333 20.1394 37.3333 18.6667V13.3333C37.3333 11.8606 36.1394 10.6667 34.6667 10.6667Z"
        stroke={color}
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M34.6667 26.6667H13.3333C11.8606 26.6667 10.6667 27.8606 10.6667 29.3333V34.6667C10.6667 36.1394 11.8606 37.3333 13.3333 37.3333H34.6667C36.1394 37.3333 37.3333 36.1394 37.3333 34.6667V29.3333C37.3333 27.8606 36.1394 26.6667 34.6667 26.6667Z"
        stroke={color}
        strokeWidth={2}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M16 16H16.0133"
        stroke={color}
        strokeWidth="2.66667"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M16 32H16.0133"
        stroke={color}
        strokeWidth="2.66667"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default ResourceIcon;
