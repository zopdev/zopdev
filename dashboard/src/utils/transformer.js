export function transformResourceAuditPayload(inputPayload) {
  console.log(inputPayload);
  try {
    if (!inputPayload || !inputPayload['0'] || !inputPayload['0'].credentials) {
      throw new Error('Invalid input payload structure');
    }
    let parsedCredentials;
    try {
      parsedCredentials = JSON.parse(inputPayload['0'].credentials);
    } catch (error) {
      throw new Error('Failed to parse credentials JSON: ' + error.message);
    }
    return {
      name: inputPayload['0'].name,
      provider: inputPayload['0'].provider,
      credentials: parsedCredentials,
    };
  } catch (error) {
    console.error('Transformation error:', error);
    return { error: error.message };
  }
}
