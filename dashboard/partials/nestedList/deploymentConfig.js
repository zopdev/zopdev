function deepMerge(target, source) {
  for (const key in source) {
    if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
      target[key] = target[key] || {};
      deepMerge(target[key], source[key]);
    } else {
      target[key] = source[key];
    }
  }
}

function transformArray(data) {
  return data.reduce((result, item) => {
    const keys = item.type.split('.'); // Split the type by '.' for nesting
    let currentLevel = result;

    keys.forEach((key, index) => {
      if (!currentLevel[key]) {
        currentLevel[key] = {};
      }
      if (index === keys.length - 1) {
        // Merge the item's properties, including the 'type' key, into the deepest level
        Object.assign(currentLevel[key], item);
      } else {
        currentLevel = currentLevel[key];
      }
    });

    return result;
  }, {});
}

export function transformDeploymentConfigData(input) {
  const result = {};

  for (const [key, value] of Object.entries(input)) {
    if (Array.isArray(value)) {
      const transformedArray = transformArray(value);
      deepMerge(result, transformedArray);
    } else if (value && typeof value === 'object' && 'type' in value) {
      result[value.type] = value;
    } else {
      result[key] = value;
    }
  }

  return result;
}
