const DashboardIcon = ({ color = '#9197B3', ...props }) => {
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
      <rect x={9} y={10} width={13} height={17} rx={1} stroke={color} strokeWidth={2} />
      <rect
        x={39}
        y={38}
        width={13}
        height={18}
        rx={1}
        transform="rotate(-180 39 38)"
        stroke={color}
        strokeWidth={2}
      />
      <rect x={9} y={31} width={13} height={7} rx={1} stroke={color} strokeWidth={2} />
      <rect
        x={39}
        y={16}
        width={13}
        height={6}
        rx={1}
        transform="rotate(-180 39 16)"
        stroke={color}
        strokeWidth={2}
      />
    </svg>
  );
};

export default DashboardIcon;
