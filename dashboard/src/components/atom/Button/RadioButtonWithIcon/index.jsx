import React, { useRef, useEffect } from 'react';

const DynamicFormRadioWithIcon = ({
  defaultSelected,
  options,
  name,
  className,
  onChange,
  value: passedSingleValue,
  ...props
}) => {
  const dataArrRef = useRef(null);

  useEffect(() => {
    if (defaultSelected !== undefined && defaultSelected !== null) {
      dataArrRef.current = `option-${defaultSelected}-${name}`;
    }
  }, [defaultSelected, name]);

  const handleChange = (e) => {
    const { value, id } = e.target;
    if (id === dataArrRef.current) {
      dataArrRef.current = null;
      e.target.checked = false;
      if (onChange) onChange(false);
    } else {
      dataArrRef.current = id;
      if (onChange) onChange(value);
    }
  };

  if (options && options.length > 0) {
    return (
      <div className={`flex gap-3 items-center relative flex-wrap mt-2`}>
        {options.map((option, index) => {
          const inputId = `option-${index}-${name}`;

          return (
            <div
              className="flex items-center"
              key={`index-${index}`}
              style={{
                minWidth: props.minWidth ?? '',
                maxWidth: props.maxWidth ?? '270px',
              }}
            >
              <input
                type="radio"
                className={`absolute opacity-0 -z-10 w-0 peer ${className}`}
                id={inputId}
                onClick={handleChange}
                value={option.value}
                name={name}
                defaultChecked={defaultSelected === option.value}
              />
              <label
                htmlFor={inputId}
                className="text-sm flex gap-4 cursor-pointer font-medium text-secondary-600 w-full p-2 px-4 border border-borderDefault rounded-md peer-checked:bg-primary-200/10 peer-checked:text-primary-500 peer-checked:border-primary-400"
                style={{
                  flexDirection: props.orientation === 'horizontal' ? 'row' : 'column',
                }}
              >
                {option.icon && (
                  <img
                    src={option.icon}
                    className={`w-10 self-center h-10 aspect-square object-cover`}
                  />
                )}
                <div>
                  {option.label && <div className={`font-semibold`}>{option.label}</div>}
                  {option.description && (
                    <div className={`text-secondary-500`}>{option.description}</div>
                  )}
                </div>
              </label>
            </div>
          );
        })}
      </div>
    );
  }

  const inputId = `option-${name}`;
  return (
    <div
      className="flex items-center mt-2"
      style={{
        minWidth: props.minWidth ?? '',
        maxWidth: props.maxWidth ?? '270px',
      }}
    >
      <input
        type="radio"
        className={`absolute opacity-0 -z-10 w-0 peer ${className}`}
        id={inputId}
        onClick={handleChange}
        value={props?.data?.value}
        name={name}
        defaultChecked={defaultSelected === props?.data?.value}
      />
      <label
        htmlFor={inputId}
        className="text-sm flex gap-4 cursor-pointer font-medium text-secondary-600  w-full p-2 px-4 border border-borderDefault rounded-md peer-checked:bg-primary-200/10 peer-checked:text-primary-500 peer-checked:border-primary-400"
        style={{
          flexDirection: props.orientation === 'horizontal' ? 'row' : 'column',
        }}
      >
        {props?.data?.icon && (
          <img
            src={props.data.icon}
            className={`w-10 self-center h-10 aspect-square object-cover`}
          />
        )}
        <div>
          {props?.data?.label && <div className={`font-semibold`}>{props?.data?.label}</div>}
          {props?.data?.description && (
            <div className={`text-secondary-500`}>{props?.data?.description}</div>
          )}
        </div>
      </label>
    </div>
  );
};

export default DynamicFormRadioWithIcon;
