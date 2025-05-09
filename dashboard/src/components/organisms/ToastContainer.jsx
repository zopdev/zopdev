import React, { useState, useSyncExternalStore } from 'react';
import { CheckIcon, XMarkIcon } from '@heroicons/react/20/solid';
import { NoSymbolIcon } from '@heroicons/react/16/solid';
import { toastStore } from '@/components/molecules/Toast/index.jsx';
import { InformationCircleIcon } from '@heroicons/react/24/outline/index.js';

// THEME_TYPES = normal | mac

const classNames = {
  success: ' bg-green-600',
  info: 'bg-primary-500',
  failed: 'bg-red-600',
  warning: 'bg-yellow-600',
};

const toasterBgclassNames = {
  success: ' bg-green-50',
  info: 'bg-primary-50',
  failed: 'bg-red-50',
  warning: 'bg-yellow-50',
};

const ToastContainer = ({ classNameParent, classNameToast, stacked = false, theme = 'normal' }) => {
  const toast = useSyncExternalStore(toastStore.subscribe, toastStore.getSnapshot);
  const [stackState, setStackState] = useState(stacked);

  return (
    <div className={`${classNameParent}`}>
      <div
        className={`flex flex-col justify-center gap-2 w-[300px] relative !overflow-hidden z-[100]`}
      >
        {toast.map((data, index) => {
          const isProgressVisible = data.progressBar && data.timeout;
          if (theme === 'mac') {
            return (
              <div
                onMouseOverCapture={() => {
                  if (index === 0) setStackState(false);
                }}
                onMouseOutCapture={() => {
                  if (index === 0) setStackState(stacked);
                }}
                key={`toast_${data.id}`}
                id={`toast_${index}`}
                className={`z-[100] text-left w-full ${
                  toasterBgclassNames[data.type]
                } transition-all fadeIn ${
                  isProgressVisible ? 'rounded-t-md' : 'rounded-md'
                } px-2 py-2 flex justify-between items-center ${classNameToast} 
              ${stackState ? `absolute transition-all` : 'relative'} ${data.className}`}
                style={{
                  top: stackState ? index * 5 : 'auto',
                  zIndex: 500 - index,
                  scale: stackState ? `${1 - index * 0.05}` : '1',
                  fontSize: 15,
                  animation: data.autoClose ? `fadeOut 0.5s linear forwards` : '',
                  animationDelay: data.autoClose ? `${data.timeout - 500}ms` : '',
                }}
              >
                <button
                  className={`absolute flex items-center rounded-full bg-red-400 w-[10px] h-[10px] top-1 right-7`}
                  onClick={() => {
                    toastStore.removeToast(data.id);
                  }}
                ></button>

                <span
                  className={`absolute flex items-center rounded-full bg-yellow-400 w-[10px] h-[10px] top-1 right-4`}
                ></span>

                <span
                  className={`absolute flex items-center rounded-full bg-green-400 w-[10px] h-[10px] top-1 right-1`}
                ></span>

                {data.type === 'success' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-green-100`}
                  >
                    <CheckIcon className={`text-green-500 size-6`} />
                  </span>
                )}
                {data.type === 'failed' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-red-100`}
                  >
                    <NoSymbolIcon className={`text-red-500 size-5`} />
                  </span>
                )}
                {data.type === 'info' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-blue-100`}
                  >
                    <CheckIcon className={`text-blue-500`} fontSize={'small'} />
                  </span>
                )}
                {data.type === 'warning' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-red-100`}
                  >
                    <NoSymbolIcon className={`text-yellow-500 size-5`} />
                  </span>
                )}

                <span className={`p-1 w-full`}>{data.title}</span>

                {isProgressVisible && (
                  <span
                    className={`absolute bottom-0 left-0 w-0 h-[3px] bg-blue-500`}
                    style={{
                      animation: `widthAnimation ${data.timeout - 100}ms linear forwards`,
                    }}
                  ></span>
                )}
              </div>
            );
          } else
            return (
              <div
                onClick={() => {
                  toastStore.removeToast(data.id);
                }}
                onMouseOverCapture={() => {
                  if (index === 0) setStackState(false);
                }}
                onMouseOutCapture={() => {
                  if (index === 0) setStackState(stacked);
                }}
                key={`toast_${data.id}`}
                id={`toast_${index}`}
                className={`z-40 text-sm cursor-pointer text-left w-full ${
                  toasterBgclassNames[data.type]
                }  transition-all fadeIn ${
                  isProgressVisible ? 'rounded-md' : 'rounded-md'
                } px-2 py-3 flex justify-between items-center ${classNameToast} 
              ${stackState ? `absolute transition-all` : ''} ${data.className}`}
                style={{
                  top: stackState ? index * 5 : 'auto',
                  zIndex: 500 - index,
                  scale: stackState ? `${1 - index * 0.05}` : '1',
                  animation: data.autoClose ? `fadeOut 0.5s linear forwards` : '',
                  animationDelay: data.autoClose ? `${data.timeout - 500}ms` : '',
                }}
              >
                {data.type === 'success' && (
                  <>
                    <span
                      className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-green-100`}
                    >
                      <CheckIcon className={`text-green-500 size-5`} />
                    </span>
                  </>
                )}
                {data.type === 'failed' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-red-100`}
                  >
                    <NoSymbolIcon className={`text-red-500 size-5`} />
                  </span>
                )}
                {data.type === 'info' && (
                  <span
                    className={`p-1 rounded-full aspect-square w-[35px] flex justify-center items-center bg-primary-100`}
                  >
                    <InformationCircleIcon className={`text-primary-500 size-5`} />
                  </span>
                )}

                <span className={`p-1 w-full`}>{data.title}</span>

                <button
                  className={`ml-2 flex items-center h-full`}
                  onClick={() => {
                    toastStore.removeToast(data.id);
                  }}
                >
                  <XMarkIcon className={`ml-1 text-gray-400 size-4`} />
                </button>
                {isProgressVisible && (
                  <span
                    className={`absolute bottom-0 left-0 w-0 h-[3px] rounded-l-md ${
                      classNames[data.type]
                    }`}
                    style={{
                      animation: `widthAnimation ${data.timeout - 500}ms linear forwards`,
                    }}
                  ></span>
                )}
              </div>
            );
        })}
      </div>
    </div>
  );
};

export default ToastContainer;
