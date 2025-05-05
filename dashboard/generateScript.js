const fs = require('fs');
const path = require('path');

const envFilePath = path.join(__dirname, '.env');

fs.readFile(envFilePath, 'utf8', (err, data) => {
  if (err) {
    console.error('Error reading the .env file:', err);
    return;
  }

  const envVariables = data
    .split('\n')
    .filter((line) => line.trim() && !line.startsWith('#'))
    .map((line) => {
      const [key, value] = line.split('=');
      return `-e ${key}='${value}'`;
    })
    .join(' ');

  const result = `docker run ${envVariables} -p 3000:3000 nextjs-image`;

  console.log(result);
});
