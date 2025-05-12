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
