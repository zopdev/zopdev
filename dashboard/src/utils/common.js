export const enforceCharLimit = (value, limit) => {
  let nonSpaceChars = value.replace(/\s+/g, '');
  if (nonSpaceChars.length > 2) {
    nonSpaceChars = nonSpaceChars.substring(0, limit);
    let newValue = '';
    let nonSpaceCount = 0;

    for (const char of value) {
      if (char !== ' ' && nonSpaceCount < limit) {
        newValue += char;
        nonSpaceCount++;
      } else if (char === ' ') {
        newValue += char;
      }
    }

    return newValue;
  }

  return value;
};

export function ParseJSON(str) {
  try {
    if (typeof str !== 'string') {
      throw new Error('input is not string');
    }
    if (typeof str === 'string' && !isNaN(str)) {
      throw new Error('Invalid JSON: Numeric string');
    }
    const result = JSON.parse(str);
    return result;
  } catch {
    return {};
  }
}
