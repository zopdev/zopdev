import React, { useState, useEffect, useRef } from 'react';

import FullScreenOverlay from '@/components/atom/FullScreenOverlay/index.jsx';
import Button from '@/components/atom/Button/index.jsx';

const RenderDeleteContent = ({
  title,
  renderComponent,
  isConfirmation,
  confirmText,
  setConfirmText,
  deleteKey,
  isDeleteEnabled,
  isLoading,
  handleClose,
  handleDelete,
  deleteLabel,
}) => {
  return (
    <div className="p-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-4">{title}</h2>
      {renderComponent && <div className="mb-6 text-gray-600">{renderComponent()}</div>}
      {isConfirmation && (
        <div className="mb-6">
          <p className="text-gray-600 font-medium mb-3">
            To confirm, type <span className="text-red-500 font-semibold">{deleteKey}</span> in the
            box below.
          </p>
          <input
            type="text"
            value={confirmText}
            onChange={(e) => setConfirmText(e.target.value)}
            placeholder={deleteKey}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none text-red-500 placeholder-gray-400/75 italic"
          />
        </div>
      )}
      <div className="flex justify-end gap-3">
        <Button
          variant="secondary"
          onClick={handleClose}
          className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-colors"
        >
          Cancel
        </Button>

        <Button
          variant="danger"
          loading={isLoading}
          onClick={handleDelete}
          disabled={!isDeleteEnabled || isLoading}
          className={`px-4 py-2 text-white rounded-md min-w-[80px] flex items-center justify-center transition-all ${
            !isDeleteEnabled || isLoading
              ? 'bg-gray-400 cursor-not-allowed'
              : 'bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-1 focus:ring-red-500'
          }`}
        >
          {deleteLabel}
        </Button>
      </div>
    </div>
  );
};
const DeleteModal = ({
  isOpen = false,
  onClose,
  onDelete,
  customCTA,
  deleteKey = 'CONFIRM',
  isConfirmation = false,
  deleteTitle,
  title = 'Are you sure you want to delete?',
  renderComponent,
  deleteLabel = 'Delete',
  isLoading = false,
}) => {
  const [showModal, setShowModal] = useState(false);
  const [confirmText, setConfirmText] = useState('');
  const [isDeleteEnabled, setIsDeleteEnabled] = useState(false);

  useEffect(() => {
    setShowModal(isOpen);
  }, [isOpen]);

  useEffect(() => {
    if (isConfirmation) {
      setIsDeleteEnabled(confirmText === deleteKey);
    } else {
      setIsDeleteEnabled(true);
    }
  }, [confirmText, deleteKey, isConfirmation]);

  useEffect(() => {
    if (!showModal) {
      setConfirmText('');
    }
  }, [showModal]);

  const prevLoadingRef = useRef(isLoading);
  useEffect(() => {
    if (prevLoadingRef.current && !isLoading && showModal) {
      setTimeout(() => {
        handleClose();
      }, 500);
    }
    prevLoadingRef.current = isLoading;
  }, [isLoading, showModal]);

  const handleOpen = (e) => {
    e.stopPropagation();
    setShowModal(true);
  };

  const handleClose = () => {
    setShowModal(false);
    if (onClose) onClose();
  };

  const handleDelete = async () => {
    if (!isDeleteEnabled || isLoading) return;
    if (onDelete) {
      await onDelete();
    }
  };

  return (
    <>
      {customCTA &&
        React.Children.map(customCTA, (child) =>
          React.cloneElement(child, {
            onClick: (e) => {
              e.stopPropagation();
              handleOpen(e);
              if (child.props.onClick) {
                child.props.onClick(e);
              }
            },
          }),
        )}

      {showModal && (
        <FullScreenOverlay
          isOpen={showModal}
          onClose={handleClose}
          title={deleteTitle}
          size="md"
          maxHeight="90vh"
          renderContent={() => (
            <RenderDeleteContent
              title={title}
              renderComponent={renderComponent}
              isConfirmation={isConfirmation}
              confirmText={confirmText}
              setConfirmText={setConfirmText}
              deleteKey={deleteKey}
              isDeleteEnabled={isDeleteEnabled}
              isLoading={isLoading}
              handleClose={handleClose}
              handleDelete={handleDelete}
              deleteLabel={deleteLabel}
            />
          )}
        />
      )}
    </>
  );
};

export default DeleteModal;
