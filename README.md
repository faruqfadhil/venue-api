# VENUE BACKEND

# Running The Apps
```shell
# Run in the background:
docker-compose up -d

# Notes: The application may not start immediately,
# the application will start as soon as the all depedencies is ready (e.g database).
# Currently, we use https://github.com/ufoscout/docker-compose-wait for waiting tools. 

# We recommend to not run containers in the background.
# it will be easier for you to see the application state (ready/not).
# Use this following command:

docker-compose up --build
```
# Stopping The Apps
```shell
docker-compose down --remove-orphans --volumes
```