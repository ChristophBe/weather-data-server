
# You should always specify a full version here to ensure all of your developers
# are running the same version of Node.
FROM node:latest
# The base node image sets a very verbose log level.
ENV NPM_CONFIG_LOGLEVEL warn


RUN mkdir /app
WORKDIR /app
# Copy all local files into the image.
COPY package.json /app/package.json
COPY src /app/src
COPY public /app/public

# Build for production.
RUN npm install --save --silent
RUN npm audit fix --silent
RUN npm run build --production

# Install `serve` to run the application.
RUN npm install -g serve

# Set the command to start the node server.
CMD serve -s build

# Tell Docker about the port we'll run on.
EXPOSE 5000