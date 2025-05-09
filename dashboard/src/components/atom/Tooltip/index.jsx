import React, { useState, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';

const Tooltip = ({ title, children, arrow = true, className = {} }) => {
  const [showTooltip, setShowTooltip] = useState(false);
  const [tooltipStyle, setTooltipStyle] = useState({});
  const tooltipRef = useRef(null);
  const triggerRef = useRef(null);

  const [recalculate, setRecalculate] = useState(false);

  const calculatePosition = () => {
    if (!tooltipRef.current || !triggerRef.current) {
      return;
    }

    const tooltipRect = tooltipRef.current.getBoundingClientRect();
    const triggerRect = triggerRef.current.getBoundingClientRect();
    const viewportHeight = window.innerHeight;
    const viewportWidth = window.innerWidth;

    let top, left, isAbove;

    if (
      viewportHeight - triggerRect.bottom < tooltipRect.height &&
      triggerRect.top > tooltipRect.height
    ) {
      top = triggerRect.top - tooltipRect.height - 8;
      isAbove = true;
    } else {
      top = triggerRect.bottom + 8;
      isAbove = false;
    }

    if (triggerRect.left + tooltipRect.width / 2 > viewportWidth) {
      left = viewportWidth - tooltipRect.width - 8;
    } else if (triggerRect.left < tooltipRect.width / 2) {
      left = 8;
    } else {
      left = triggerRect.left + triggerRect.width / 2 - tooltipRect.width / 2;
    }

    setTooltipStyle({
      position: 'fixed',
      top: `${top}px`,
      left: `${left}px`,
      zIndex: 50,
      arrowPosition: isAbove ? 'bottom' : 'top',
    });
  };

  useEffect(() => {
    if (showTooltip) {
      calculatePosition();
      setRecalculate(true);
    }
  }, [showTooltip]);

  useEffect(() => {
    if (recalculate) {
      calculatePosition();
      setRecalculate(false);
    }
  }, [recalculate]);

  return (
    <div className={`relative ${className?.root ?? ''}`}>
      <div
        ref={triggerRef}
        className="inline-block"
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
      >
        {children}
      </div>

      {showTooltip &&
        title &&
        createPortal(
          <div
            ref={tooltipRef}
            style={tooltipStyle}
            className="bg-secondary-300 text-xs text-secondary-600 px-4 py-2 rounded-md transition-opacity duration-300 text-wrap max-w-80"
          >
            {title}
            {arrow && (
              <div
                className="absolute w-2 h-2 bg-secondary-300"
                style={{
                  bottom: tooltipStyle.arrowPosition === 'bottom' ? '-4px' : 'auto',
                  top: tooltipStyle.arrowPosition === 'top' ? '-4px' : 'auto',
                  left: '50%',
                  transform: 'translateX(-50%) rotate(45deg)',
                }}
              ></div>
            )}
          </div>,
          document.body,
        )}
    </div>
  );
};

export default Tooltip;
