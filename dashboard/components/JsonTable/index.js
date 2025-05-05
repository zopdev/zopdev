import React, { useEffect, useRef } from 'react';

const INITIAL_CELL_STYLES = {
  minWidth: 500,
  maxWidth: 500,
  wordBreak: 'break-all',
  marginRight: 40,
  padding: 10,
  paddingBottom: 10,
};

const INITIAL_KEY_CELL_STYLES = {
  minWidth: 200,
  maxWidth: 200,
  wordBreak: 'break-all',
  marginRight: 40,
  paddingTop: 10,
};

const COLOURS = {
  BORDER: '#e5e7eb',
};

const JsonComparisonTable = ({
  data,
  keysStyles = {},
  screenOffsetTop,
  isKeysVisible = false,
  difference = false,
  isKeysScrollable = false,
  isAlternateColStyle = true,
  isAlternateRowStyle = true,
}) => {
  const dataTableRef = useRef(null);
  const dataHeaderRef = useRef(null);
  const parentRef = useRef(null);
  const keyRefs = useRef([]);
  const cellRefs = useRef([]);

  useEffect(() => {
    const handleScrollBody = (e) => {
      if (dataHeaderRef.current) {
        dataHeaderRef.current.scrollLeft = e.target.scrollLeft;
      }
    };

    const handleScrollHeader = (e) => {
      if (dataHeaderRef.current) {
        dataTableRef.current.scrollLeft = e.target.scrollLeft;
      }
    };

    if (dataTableRef.current) {
      dataTableRef.current.addEventListener('scroll', handleScrollBody);
      dataHeaderRef.current.addEventListener('scroll', handleScrollHeader);
    }

    return () => {
      if (dataTableRef.current) {
        dataTableRef.current.removeEventListener('scroll', handleScrollBody);
        dataHeaderRef.current.removeEventListener('scroll', handleScrollHeader);
      }
    };
  }, [dataTableRef]);

  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const height = entry.contentRect.height;
        const index = Number(entry.target.getAttribute('data-index'));
        const keyElement = keyRefs.current[index];
        const maxHeight = Math.max(height, keyElement?.clientHeight);
        if (keyElement) {
          // Check if keyElement exists before accessing its style
          keyElement.style.height = `${maxHeight}px`;
        }
      }
    });

    keyRefs.current.forEach((keyRef, index) => {
      let maxHeight = 0;
      if (keyRef && cellRefs.current[index]) {
        for (const elements of cellRefs.current[index]) {
          resizeObserver.observe(elements);
          if (elements.clientHeight > maxHeight) {
            maxHeight = Math.max(maxHeight, elements.clientHeight);
          }
        }
      }

      keyRef.style.height = `${maxHeight}px`;
    });

    return () => {
      resizeObserver.disconnect();
    };
  }, [data]);

  const getColumStyles = (dataColumn) => {
    let styles = { ...INITIAL_CELL_STYLES };

    if (dataColumn) {
      styles = {
        ...styles,
        ...(dataColumn.styles ?? {}),
      };
    }

    return styles;
  };

  const keys = new Set();

  for (const columns of data) {
    for (const key in columns.data) {
      keys.add(key);
    }
  }

  const keysOfData = Array.from(keys);

  const keyCellStyles = { ...INITIAL_KEY_CELL_STYLES, ...keysStyles };

  const getFormattedCell = (dataCell) => {
    if (typeof dataCell === 'string' || typeof dataCell === 'number') {
      return dataCell;
    }

    if (typeof dataCell === 'function') {
      return dataCell();
    }

    if (typeof dataCell === 'object') {
      return JSON.stringify(dataCell);
    }

    return '';
  };

  const isKeyDifferent = (key, value) => {
    if (data.length === 0) return null;
    if (!value) return null;
    if (typeof value !== 'string') return null;

    for (let i = 1; i < data.length; i++) {
      if (data[i].data[key] && data[i].data[key] !== value) {
        return true;
      }
    }

    return false;
  };

  return (
    <div
      className={`p-6 text-gray-600 ${isKeysScrollable ? 'overflow-auto scroll-hidden' : ''}`}
      ref={parentRef}
    >
      <div className={`${screenOffsetTop ? 'max-h-[calc(100vh-100px)] overflow-y-auto' : ''}`}>
        <div className={`${screenOffsetTop ? 'sticky top-0' : ''} flex bg-white font-semibold`}>
          {isKeysVisible && (
            <div
              style={{
                ...keyCellStyles,
                paddingTop: 20,
              }}
            >
              {/* Key */}
            </div>
          )}
          <div
            className={`flex ${isKeysScrollable ? '' : 'overflow-auto scroll-hidden'}`}
            ref={dataHeaderRef}
          >
            {data.map((single, i) => (
              <div
                key={i}
                style={{
                  ...getColumStyles(single),
                  paddingTop: 10,
                }}
                className={`${
                  isAlternateRowStyle
                    ? ''
                    : isAlternateColStyle
                      ? i % 2 === 0
                        ? 'bg-gray-50 rounded-t-lg border-l border-t border-r'
                        : ''
                      : 'bg-gray-50'
                } text-center text-white`}
              >
                <div
                  className={`flex bg-primary-500 rounded-md justify-center items-center text-center py-2`}
                >
                  {single.label}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className={`flex`}>
          {isKeysVisible && (
            <div className={`flex flex-col`}>
              {keysOfData.map((key, iRow) => (
                <div
                  style={keyCellStyles}
                  className={`flex`}
                  key={iRow}
                  ref={(el) => (keyRefs.current[iRow] = el)}
                >
                  <span className={`flex justify-start items-start border-b flex-1 text-gray-400`}>
                    {key}
                  </span>
                </div>
              ))}
            </div>
          )}

          <div
            className={`flex flex-col  pb-4 ${
              isKeysScrollable ? '' : 'overflow-auto scroll-hidden'
            } text-gray-600`}
            ref={dataTableRef}
          >
            {keysOfData.map((key, iRow) => (
              <div className={`flex`} key={iRow}>
                {data.map((single, i) => (
                  <div
                    key={i}
                    data-index={iRow}
                    className={`flex gap-2 ${
                      isAlternateRowStyle
                        ? ''
                        : isAlternateColStyle
                          ? i % 2 === 0
                            ? 'bg-gray-50'
                            : ''
                          : 'bg-gray-50'
                    }`}
                    style={{
                      ...getColumStyles(single),
                      padding: 0,
                      paddingLeft: 10,
                      paddingRight: 10,
                      paddingBottom: 0,
                      ...(isAlternateRowStyle
                        ? {}
                        : isAlternateColStyle
                          ? i % 2 === 0
                            ? {
                                borderLeft: `1px solid ${COLOURS.BORDER}`,
                                borderRight: `1px solid ${COLOURS.BORDER}`,
                                ...(iRow === keysOfData.length - 1
                                  ? {
                                      borderBottom: `1px solid ${COLOURS.BORDER}`,
                                      borderBottomLeftRadius: 10,
                                      borderBottomRightRadius: 10,
                                    }
                                  : {}),
                              }
                            : {}
                          : {
                              borderLeft: `1px solid ${COLOURS.BORDER}`,
                              borderRight: `1px solid ${COLOURS.BORDER}`,
                              ...(iRow === keysOfData.length - 1
                                ? {
                                    borderBottom: `1px solid ${COLOURS.BORDER}`,
                                    borderBottomLeftRadius: 10,
                                    borderBottomRightRadius: 10,
                                  }
                                : {}),
                            }),
                    }}
                    ref={(el) => {
                      if (cellRefs.current[iRow]) {
                        cellRefs.current[iRow][i] = el;
                      } else {
                        cellRefs.current[iRow] = [el];
                      }
                    }}
                  >
                    <div
                      className={`w-full flex gap-2 ${
                        difference && isKeyDifferent(key, single.data[key]) === true
                          ? 'bg-red-50'
                          : isAlternateRowStyle
                            ? iRow % 2 === 0
                              ? ''
                              : 'bg-gray-50'
                            : ''
                      }`}
                      style={{
                        ...(difference && isKeyDifferent(key, single.data[key]) === true
                          ? {
                              borderBottom: `1px solid ${'#ffb5b5'}`,
                            }
                          : {
                              borderBottom: `1px solid ${COLOURS.BORDER}`,
                            }),
                        padding: 10,
                      }}
                    >
                      {difference && isKeyDifferent(key, single.data[key]) === true && (
                        <div className={`flex justify-center items-center `}>
                          <span
                            className={`min-w-2 w-2 min-h-2 h-2 aspect-square rounded-full bg-red-500`}
                          ></span>
                        </div>
                      )}
                      <span>{getFormattedCell(single.data[key])}</span>
                    </div>
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default function JsonComparisonTableContainer({
  data,
  isKeysVisible,
  screenOffsetTop,
  keysStyles,
  difference,
  isKeysScrollable,
  isAlternateColStyle,
  isAlternateRowStyle,
}) {
  return (
    <div className="min-h-screen">
      <JsonComparisonTable
        data={data}
        isKeysVisible={isKeysVisible}
        screenOffsetTop={screenOffsetTop}
        keysStyles={keysStyles}
        difference={difference}
        isKeysScrollable={isKeysScrollable}
        isAlternateColStyle={isAlternateColStyle}
        isAlternateRowStyle={isAlternateRowStyle}
      />
    </div>
  );
}
