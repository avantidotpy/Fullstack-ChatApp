# Use a Node.js base image
FROM node:18.16.1

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci --only=production

# Copy the rest of the application source code
COPY . .

# Build the frontend application
RUN npm run build

# Expose the port the server will run on
EXPOSE 4000

# Start the frontend server
CMD ["npm", "start"]
