# testcontainers-485-repro

Reproduces testcontainers/testcontainers-go#485 error by trying to run multiple postgres images at once. In my laptop 97 of the 100 containers failed to start with the error `Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.42/containers/[hash]/json": context deadline exceeded: failed to start container`
