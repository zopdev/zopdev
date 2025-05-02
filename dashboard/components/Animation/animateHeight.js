import React, { useEffect, useRef, useState } from 'react';

const AnimateHeight = ({ children, className, dependency, customAnimation }) => {
  const [height, setHeight] = useState(0);
  const parentRef = useRef(null);

  let time;
  const handleChange = () => {
    if (!parentRef.current) {
      return;
    }
    const childFirst = parentRef.current.children[0];
    const heightChild = childFirst.getBoundingClientRect().height;
    setHeight(heightChild);
    time = setTimeout(() => {
      if (heightChild !== childFirst.getBoundingClientRect().height) {
        handleChange();
      }
    }, 100);
  };

  useEffect(() => {
    handleChange();
    return () => {
      clearTimeout(time);
    };
  }, [children, dependency]);

  useEffect(() => {
    if (!parentRef.current) {
      return;
    }
    const resizeOb = new ResizeObserver(handleChange);
    const childFirst = parentRef.current.children[0];
    if (childFirst) resizeOb.observe(childFirst);

    return () => {
      if (childFirst) resizeOb.unobserve(childFirst);
    };
  }, [parentRef.current]);

  return (
    <div className={`overflow-hidden ${className}`}>
      <div
        className={`ease-linear transition-all ${customAnimation || 'duration-500'}`}
        style={{
          height,
        }}
        ref={parentRef}
      >
        <div>{children}</div>
      </div>
    </div>
  );
};

export default AnimateHeight;
