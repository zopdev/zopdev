import React from 'react';

const DotWithProgress = ({ color }) => {
    return (
        <div
            style={{
                position: 'relative',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                width: 14,
                height: 14,
            }}
        >
      <span
          style={{
              position: 'absolute',
              width: 6,
              height: 6,
              backgroundColor: color,
              borderRadius: '50%',
              zIndex: 1,
          }}
      />
            <SimpleLoader color={color} size={17} />
        </div>
    );
};

export default DotWithProgress;
