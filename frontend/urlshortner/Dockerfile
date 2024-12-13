# Use an official Node.js image as the base image
FROM node:18-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the build artifacts from the local machine into the container
COPY ./dist ./build

# Expose the port the application will run on
EXPOSE 8080

# Install a simple HTTP server to serve the build (if it's a static site)
RUN npm install -g serve

# Start the application using the built files
CMD ["serve", "-s", "build", "-l", "8080"]
