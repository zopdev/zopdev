let nextId = 0;
let toasts = [];
let listeners = [];

const DEFAULT_OPTIONS = {
  timeout: 6000,
  type: 'info',
  className: '',
  progressBar: true,
  autoClose: true,
};

// const TYPES = ['info', 'success', 'failed'];
const classNames = {
  success: ' text-green-600 ',
  info: 'text-primary-600',
  failed: 'text-red-600',
  warning: 'text-yellow-600',
};

export const toastStore = {
  addToast(toast, options) {
    const newOptions = {
      ...DEFAULT_OPTIONS,
      className: classNames[options.type],
      ...options,
    };
    const newToast = { id: nextId++, title: toast, ...newOptions };
    if (newOptions.autoClose) {
      setTimeout(() => {
        this.removeToast(newToast.id);
      }, newOptions.timeout);
    }
    toasts = [...toasts, newToast];
    emitChange();
  },
  subscribe(listener) {
    listeners = [...listeners, listener];
    return () => {
      listeners = listeners.filter((l) => l !== listener);
    };
  },
  getSnapshot() {
    return toasts;
  },
  removeToast(toastId) {
    toasts = toasts.filter((toast) => toast.id !== toastId);
    emitChange();
  },
};

function emitChange() {
  for (const listener of listeners) {
    listener();
  }
}

export const toast = {
  success: (data, option) => {
    toastStore.addToast(data || 'Action completed successfully!', {
      ...option,
      type: 'success',
    });
  },
  failed: (data, option) => {
    toastStore.addToast(data || 'An error occurred. Please try again', {
      ...option,
      type: 'failed',
    });
  },
  info: (data, option) => {
    toastStore.addToast(data || `Here’s something you should know.`, {
      ...option,
      type: 'info',
    });
  },
  warning: (data, option) => {
    toastStore.addToast(data || `Here’s something you should know.`, {
      ...option,
      type: 'warning',
    });
  },
};
