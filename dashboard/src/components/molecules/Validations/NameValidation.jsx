export default function NameValidation({ value, type, field, min, max }) {
  if (min && max) {
    return (
      <ValidationMessage
        value={value}
        maxLength={max}
        minLength={min}
        type={type}
        pattern={/^[a-zA-Z][-a-zA-Z0-9_ ]*$/}
      />
    );
  }
  if (field === 'service') {
    return (
      <ValidationMessage
        value={value}
        maxLength={52}
        minLength={1}
        type={type}
        pattern={/^[a-zA-Z][-a-zA-Z0-9_ ]*$/}
      />
    );
  }
  return (
    <ValidationMessage
      value={value}
      maxLength={16}
      minLength={6}
      type={type}
      pattern={/^[a-zA-Z][-a-zA-Z0-9_ ]*$/}
    />
  );
}

export function ValidationMessage({ value, maxLength, minLength, pattern, type }) {
  if (value === undefined || value === null || value.trim() === '') return <span></span>;
  if (value && !value.match(/^[a-zA-Z]/)) {
    return <span className={`text-yellow-500`}>Must start with alphabet</span>;
  }

  if (value.length < minLength && type !== 'azure')
    return (
      <span className={`text-yellow-500`}>Minimum {minLength} characters should be there</span>
    );

  if (value.replace(/\s+/g, '').length < minLength && type === 'azure')
    return (
      <span className={`text-yellow-500`}>Minimum {minLength} characters should be there</span>
    );

  if (value.length > 1 && !value.match(/^[a-z0-9-]+$/) && type === 'helm-name') {
    return (
      <span className={`text-yellow-500`}>
        Allowed types are alphabets in lowercase, numbers and -
      </span>
    );
  }

  if (pattern && !pattern.test(value)) {
    if (
      value.length > 1 &&
      !value.match(/^[-a-zA-Z0-9_ ][-a-zA-Z0-9_ ]*[-a-zA-Z0-9_ ]$/) &&
      type !== 'azure'
    ) {
      return (
        <span className={`text-yellow-500`}>Allowed types are alphabets, numbers, _ and space</span>
      );
    }
    return (
      <span className={`text-yellow-500`}>
        Must start with alphabet and cannot end with special characters
      </span>
    );
  }

  if (value && !value.match(/[a-zA-Z0-9]$/)) {
    return <span className={`text-yellow-500`}>Must end with alphabet or number</span>;
  }

  if (value.replace(/\s+/g, '').length === 11 && type === 'azure') {
    return (
      <span className={`text-yellow-500`}>Maximum 11 characters are allowed, limit reached</span>
    );
  }

  if (value.length >= maxLength && type !== 'azure')
    return (
      <span className={`text-yellow-500`}>
        Maximum {maxLength} characters are allowed, limit reached
      </span>
    );

  return <span className={`text-yellow-500`}></span>;
}
