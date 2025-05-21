import { Field, Label, Switch } from '@headlessui/react';
import { useState, useId } from 'react';

function classNames(...classes) {
  return classes.filter(Boolean).join(' ');
}

export default function SwitchButton({
  isEnabled,
  onChange,
  title,
  titleList,
  disabled,
  name,
  labelPosition,
}) {
  const [enabled, setEnabled] = useState(
    isEnabled === 'true' ? true : isEnabled === 'false' ? false : !!isEnabled,
  );
  const switchId = useId();

  const handleChange = (value) => {
    if (disabled) return;
    setEnabled(value);
    onChange?.(value);
  };

  const toggleSwitch = () => {
    handleChange(!enabled);
  };

  const labelTitle = titleList?.[enabled] ?? title;

  return (
    <Field
      as="div"
      className={classNames(
        'flex items-center space-x-3',
        labelPosition === 'right' && 'flex-row-reverse  justify-end',
      )}
    >
      <Label
        as="span"
        htmlFor={switchId}
        onClick={toggleSwitch}
        className={classNames(
          'text-sm text-gray-600 cursor-pointer',
          disabled && 'opacity-50 cursor-default',
        )}
      >
        {labelTitle}
      </Label>
      <Switch
        data-name={name}
        datatype="boolean"
        id={switchId}
        checked={enabled}
        onChange={handleChange}
        className={classNames(
          enabled ? (disabled ? 'bg-primary-500' : 'bg-primary-600') : 'bg-gray-300',
          'relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
          disabled ? 'opacity-50 cursor-auto' : 'cursor-pointer',
        )}
        disabled={disabled}
      >
        <span
          className={classNames(
            enabled ? 'translate-x-5' : 'translate-x-0',
            'pointer-events-none relative inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
          )}
        >
          <span
            className={classNames(
              enabled ? 'opacity-0 duration-100 ease-out' : 'opacity-100 duration-200 ease-in',
              'absolute inset-0 flex h-full w-full items-center justify-center transition-opacity',
            )}
            aria-hidden="true"
          >
            <svg className="h-3 w-3 text-gray-500" fill="none" viewBox="0 0 12 12">
              <path
                d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
                stroke="currentColor"
                strokeWidth={2}
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </span>
          <span
            className={classNames(
              enabled ? 'opacity-100 duration-200 ease-in' : 'opacity-0 duration-100 ease-out',
              'absolute inset-0 flex h-full w-full items-center justify-center transition-opacity',
            )}
            aria-hidden="true"
          >
            <svg className="h-3 w-3 text-primary-600" fill="currentColor" viewBox="0 0 12 12">
              <path d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z" />
            </svg>
          </span>
        </span>
      </Switch>
    </Field>
  );
}
