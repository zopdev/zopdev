import { Fragment, useState, cloneElement } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { XMarkIcon } from '@heroicons/react/24/outline/index.js';

const FullScreenOverlay = ({
  isOpen,
  onClose,
  customCTA,
  title,
  renderContent,
  variant = 'popup',
  position,
  size = 'md',
  hasCloseIcon = true,
  closeOnOutsideClick = true,
  overlayColor = 'light',
  maxHeight = '80vh',
  renderContentProps,
}) => {
  const defaultPosition = variant === 'popup' ? 'center' : 'right';
  const effectivePosition = position || defaultPosition;
  const isExternallyControlled = isOpen !== undefined && onClose !== undefined;
  const [internalIsOpen, setInternalIsOpen] = useState(false);
  const effectiveIsOpen = isExternallyControlled ? isOpen : internalIsOpen;

  const closePopup = () => {
    if (isExternallyControlled) {
      onClose?.();
    } else {
      setInternalIsOpen(false);
    }
  };

  const openPopup = () => {
    if (!isExternallyControlled) {
      setInternalIsOpen(true);
    }
  };

  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
    '3xl': 'max-w-3xl',
    '4xl': 'max-w-4xl',
    full: 'max-w-full mx-4',
  };

  const overlayClasses = {
    dark: 'bg-black/20',
    light: 'bg-white/20',
  };

  const drawerPositionClasses = {
    left: 'left-0 top-0 bottom-0 h-full w-full',
    right: 'right-0 top-0 bottom-0 h-full w-full',
  };

  const popupPositionClasses = {
    center: 'items-center justify-center',
    top: 'items-start justify-center pt-20',
    bottom: 'items-end justify-center pb-20',
    left: 'items-center justify-start pl-20',
    right: 'items-center justify-end pr-20',
  };

  const drawerEnterFrom = {
    left: '-translate-x-full',
    right: 'translate-x-full',
  };

  const drawerLeaveTo = drawerEnterFrom;

  const renderTrigger = () => {
    if (!customCTA) return null;

    return cloneElement(customCTA, {
      onClick: (e) => {
        if (customCTA.props.onClick) {
          customCTA.props.onClick(e);
        }
        openPopup();
      },
    });
  };

  return (
    <>
      {renderTrigger()}
      <Transition show={!!effectiveIsOpen} as={Fragment}>
        <Dialog
          as="div"
          className="fixed inset-0 z-50 overflow-y-auto"
          onClose={closeOnOutsideClick ? closePopup : () => {}}
          static
        >
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className={`fixed inset-0 ${overlayClasses[overlayColor]}`} />
          </Transition.Child>
          <div
            className={`fixed inset-0 ${overlayClasses[overlayColor]} bg-opacity-10 backdrop-blur transition-opacity `}
          />
          <div
            className={`fixed inset-0 flex ${variant === 'popup' ? popupPositionClasses[effectivePosition] : ''}`}
          >
            {variant === 'popup' ? (
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel
                  className={`flex flex-col bg-white border border-gray-300 rounded-sm shadow-xl ${sizeClasses[size]} w-full m-4 overflow-hidden`}
                  style={{ maxHeight }}
                >
                  <div className="flex-shrink-0 flex items-center justify-between px-5 py-3 border-b-2 border-gray-300">
                    <Dialog.Title className="font-medium text-l text-gray-900 pr-8">
                      {title}
                    </Dialog.Title>
                    {hasCloseIcon && (
                      <button
                        onClick={closePopup}
                        className="text-gray-400 cursor-pointer hover:text-gray-600 rounded-full hover:bg-gray-200 p-1 focus:outline-none transition-colors duration-200"
                      >
                        <XMarkIcon className="h-5 w-5" />
                        <span className="sr-only">Close</span>
                      </button>
                    )}
                  </div>
                  <div className="flex-grow px-6 py-4 overflow-y-auto">
                    {renderContent && renderContent({ onClose: closePopup, ...renderContentProps })}
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            ) : (
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-300"
                enterFrom={drawerEnterFrom[effectivePosition]}
                enterTo="translate-x-0 translate-y-0"
                leave="transform transition ease-in-out duration-300"
                leaveFrom="translate-x-0 translate-y-0"
                leaveTo={drawerLeaveTo[effectivePosition]}
              >
                <Dialog.Panel
                  className={`fixed flex flex-col bg-white shadow-xl ${sizeClasses[size]} ${drawerPositionClasses[effectivePosition]} overflow-hidden`}
                >
                  <div className="flex-shrink-0 flex items-center justify-between px-6 py-4 border-b border-gray-200">
                    <Dialog.Title className="text-lg font-semibold text-gray-900">
                      {title}
                    </Dialog.Title>
                    {hasCloseIcon && (
                      <button
                        onClick={closePopup}
                        className="text-gray-400 hover:text-gray-600 focus:outline-none transition-colors duration-200"
                      >
                        <XMarkIcon className="h-5 w-5" />
                        <span className="sr-only">Close</span>
                      </button>
                    )}
                  </div>
                  <div className="flex-grow px-6 py-4 overflow-y-auto">
                    {renderContent && renderContent({ onClose: closePopup, ...renderContentProps })}
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            )}
          </div>
        </Dialog>
      </Transition>
    </>
  );
};
export default FullScreenOverlay;
