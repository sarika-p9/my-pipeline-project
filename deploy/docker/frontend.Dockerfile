# Build the React app using Node.js
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package.json and package-lock.json for npm install
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

# Copy the rest of the React source code
COPY frontend/ ./

# Build the React app
RUN npm run build

# Use NGINX to serve the built app
FROM nginx:alpine

# Copy custom nginx config (make sure you have nginx.conf in the deploy directory)
COPY deploy/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy built React app to NGINX's html directory
COPY --from=builder /app/build /usr/share/nginx/html

# Expose the correct port for NGINX
EXPOSE 80

# Start NGINX
CMD ["nginx", "-g", "daemon off;"]
