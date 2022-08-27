# Cake Store API

Restful API for testing purposes only
- Get cake list
- Get cake
- Create cake
- Update cake
- Delete cake

## Running the migrator

```sh
$ make migrate_up
```

## Running the app locally

```sh
$ go build
$ ./main
⇨ http server started on [::]:8080
``` 
- [route test](http://localhost:8080)

```sh
$ curl http://localhost:8080
{"message":"API OK"}
```

## Building and running the docker image

```sh
$ make docker_build
$ make docker_run
⇨ http server started on [::]:8080
```

## For detailed documentation

- [Docs](http://localhost:8080/docs/)
