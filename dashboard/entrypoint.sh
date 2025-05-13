#!/bin/sh

echo "Generating runtime env.js..."
cat <<EOF > /website/env.js
window.env = {
  API_BASE_URL: "${API_BASE_URL}"
};
EOF

echo "Starting static server..."
exec /main
