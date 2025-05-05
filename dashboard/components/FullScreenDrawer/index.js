'use client';

import React, { useState, Fragment } from 'react';
import { Dialog, DialogPanel, DialogTitle, Transition, TransitionChild } from '@headlessui/react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import IconButton from '../Button/IconButton';

export function FullScreenDrawer({
  title,
  RenderIcon,
  RenderComponent,
  disabled,
  isIcon,
  handleMenuItemClose = () => {},
  formData,
}) {
  const [open, setOpen] = useState(false);

  const handleClickOpen = (e) => {
    e.stopPropagation();
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSuccess = () => {
    setOpen(false);
  };

  return (
    <>
      {isIcon && (
        // <Tooltip title={tooltipTitle}>
        <IconButton onClick={(e) => handleClickOpen(e)} disabled={disabled}>
          {RenderIcon}
        </IconButton>
        // </Tooltip>
      )}

      <Transition show={open} as={Fragment}>
        <Dialog className="relative z-10" onClose={handleClose}>
          <TransitionChild
            as={Fragment}
            enter="ease-in-out duration-500"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in-out duration-500"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-gray-800 bg-opacity-10 backdrop-blur transition-opacity" />
          </TransitionChild>

          <div className="fixed inset-0 overflow-hidden">
            <div className="absolute inset-0 overflow-hidden">
              <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full xs:pl-0 md:pl-10">
                <TransitionChild
                  as={Fragment}
                  enter="transform transition ease-in-out duration-500 sm:duration-700"
                  enterFrom="translate-x-full"
                  enterTo="translate-x-0"
                  leave="transform transition ease-in-out duration-500 sm:duration-700"
                  leaveFrom="translate-x-0"
                  leaveTo="translate-x-full"
                >
                  <DialogPanel className="pointer-events-auto w-screen xs:min-w-screen md:max-w-3xl ">
                    <div className="flex h-full flex-col bg-gray-50  shadow-xl">
                      <div className="px-4 sm:px-8 bg-primary-500 py-3">
                        <div className="flex items-start justify-between ">
                          <DialogTitle className="font-semibold !text-white  sm:text-2xl xs:text-xl leading-6 ">
                            {title}
                          </DialogTitle>
                          <div className="ml-3 flex h-7 items-center">
                            <IconButton
                              onClick={() => {
                                handleClose();
                                handleMenuItemClose();
                              }}
                            >
                              <XMarkIcon className="h-6 w-6 text-gray-600" aria-hidden="true" />
                            </IconButton>
                          </div>
                        </div>
                      </div>
                      <div className="relative mt-4 pb-4 flex-1 px-4 sm:px-2 overflow-y-scroll">
                        <RenderComponent
                          formData={formData}
                          handleClose={handleClose}
                          setOpen={handleSuccess}
                        />
                      </div>
                    </div>
                  </DialogPanel>
                </TransitionChild>
              </div>
            </div>
          </div>
        </Dialog>
      </Transition>
    </>
  );
}
