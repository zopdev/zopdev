const sizeClasses = {
  xl: {
    wrapper: 'gap-3 p-2 text-base',
    button: 'px-4 py-2',
    active: 'p-2',
  },
  lg: {
    wrapper: 'gap-2.5 p-1.5 text-sm',
    button: 'px-3 py-1.5',
    active: 'p-1.5',
  },
  md: {
    wrapper: 'gap-2 p-1 text-sm',
    button: 'px-2 py-1',
    active: 'p-1',
  },
  sm: {
    wrapper: 'gap-1.5 p-0.5 text-xs',
    button: 'px-1.5 py-0.5',
    active: 'p-0.5',
  },
  xs: {
    wrapper: 'gap-1 p-0.5 text-[10px]',
    button: 'px-1 py-0.5',
    active: 'p-0.5',
  },
};

export function Tabs({ tabs, activeTab, onTabChange, size = 'xl' }) {
  const { wrapper, button, active } = sizeClasses[size] || sizeClasses.xl;

  return (
    <div className={`flex   space-x-8  rounded-sm ${wrapper} `}>
      {tabs.map(({ label }) => {
        const isActive = activeTab === label;
        return (
          <button
            key={label}
            type="button"
            onClick={() => onTabChange(label)}
            className={`transition-all cursor-pointer duration-200 ease-in-out text-black hover:border-primary-600 hover:text-primary-600 group inline-flex items-center border-b-2  border-borderDefault py-3 px-1 font-medium ${button} ${
              isActive ? ` text-primary-600 border-primary-600 ${active}` : ''
            }`}
          >
            {label}
          </button>
        );
      })}
    </div>
  );
}
