import { useState } from 'react';

const CopyToClipboardButton = ({ isYaml }) => {
  const [copied, setCopied] = useState(false);
  const [showToast, setShowToast] = useState(false);

  const handleClick = (e) => {
    if (!copied) {
      const codeElement = e.target.closest('.code-container');
      if (codeElement) {
        const codeText = codeElement.innerText || codeElement.textContent;
        navigator.clipboard.writeText(codeText).then(() => {
          setCopied(true);
          setShowToast(true);
          setTimeout(() => {
            setCopied(false);
          }, 2000);
          setTimeout(() => {
            setShowToast(false);
          }, 2000);
        });
      }
    }
  };

  return (
    <>
      <button
        onClick={handleClick}
        className={`copy-button ${isYaml ? 'yaml-code' : ''}`}
        style={{
          float: 'right',
          padding: '6px',
          background: 'transparent',
          border: 'none',
          cursor: 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        {copied ? <CheckCircleIcon color="#06b6d4" /> : <CopyIcon color="#06b6d4" />}
      </button>

      {showToast && (
        <div
          className="toast"
          style={{
            position: 'fixed',
            bottom: '20px',
            left: '20px',
            background: 'rgba(0, 0, 0, 0.7)',
            color: 'white',
            padding: '10px 20px',
            borderRadius: '4px',
            zIndex: 1000,
          }}
        >
          Copied to clipboard
        </div>
      )}
    </>
  );
};

const CopyIcon = ({ color }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
  </svg>
);

const CheckCircleIcon = ({ color }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
    <polyline points="22 4 12 14.01 9 11.01"></polyline>
  </svg>
);

export default CopyToClipboardButton;
