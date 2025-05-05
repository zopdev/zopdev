# Stage 1: Build the application
FROM node:18-alpine

# Create and set the working directory for the app inside the container
WORKDIR /app

# Copy the app's source code into the container
COPY . .

# Install dependencies
RUN npm install

# Build the Next.js application
RUN npm run build

# Expose the port the app will run on
EXPOSE 3000

# Set the default command to run when the container starts (but don't run the app yet)
CMD ["npm", "run", "dev"]