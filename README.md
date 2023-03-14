# bank

- [Get started](#get-started)
  - [Build source code](#build-source-code)
  - [Executing binary](#executing-binary)
  - [Containerizing binary](#containerizing-binary)
  - [Executing container with binary](#executing-container-with-binary)
  
- [Tasks](#tasks)
  - [Running lint](#running-lint)
  - [Running unit tests](#running-only-unit-tests)
  - [Running integration tests](#running-all-tests-including-integration-tests)

- [Source code design pattern](#source-code-design-pattern)
  - [Tree overview](#tree-overview)
  - [Concepts of source code arch](#concepts-of-source-code-arch)

## Get started

### Build source code
```sh
$ CGO_ENABLED=0 go build -v -o ./dist/api ./cmd/api/main.go
```

### Executing binary

> :balloon: It's necessary has [postgres](https://www.postgresql.org/) installed and running heathly, for more configuration details 
take a look at [docker-compose.yaml]() var environments

```sh
$ ./dist/api
```

### Containerizing binary

> :balloon: It's necessary has [docker](https://www.docker.com/get-started/) installed

```sh
$ docker build -t bank .
```

### Executing container with binary

> :balloon: It's necessary has [postgres](https://www.postgresql.org/) installed and running heathly, for more configuration details 
take a look at [docker-compose.yaml]() var environments

```sh
$ docker run -it bank
```

## Tasks

> :balloon: This project has `Makefile` as job runner

### Running lint
```sh
$ make lint
```

### Running only unit tests
```sh
$ make test
```

### Running all tests (including integration tests)

> :balloon: It's necessary has [docker](https://www.docker.com/get-started/) and [docker compose](https://docs.docker.com/compose/) installed then
pull image `postgres:14` which's app's database tech.

```sh
$ make integration
```

### Source code design pattern

#### Tree overview

```sh
.
├── cmd
│   └── api
│       ├── main.go
│       └── main_test.go
├── docker-compose.yaml
├── Dockerfile
├── domain
│   ├── account
│   │   ├── account.go
│   │   ├── errors.go
│   │   ├── port.go
│   │   ├── usecase.go
│   │   └── usecase_test.go
│   └── log
│       └── port.go
├── go.mod
├── go.sum
├── handler
│   └── http
│       └── api
│           ├── account.go
│           └── http.go
├── infra
│   ├── logger
│   │   └── log
│   │       └── logger.go
│   └── storage
│       └── postgres
│           ├── account.go
│           ├── account_transaction.go
│           ├── operation_type.go
│           ├── storage.go
│           └── storage_test.go
├── Makefile
└── pkg
    └── docker
        ├── container.go
        └── environment.go
```

#### Concepts of source code arch

- **cmd**: This directory's responsible for app's entrypoint. In this case we have a integration for API but, for example, it could there is CLI 
integration too and both existing in the same source code.

- **handler**: This directory's responsible for app's user interface protocol. It's here that'll develop all rules for handle an input then pass to use cases/domain layer.

- **domain**: This directory's responsible for app's core, where it'll be the business rule. Given hexagonal arch, it's here that all ports communicate with use cases rules and never the inverse.

- **infra**: This directory's responsible for all app's external integration which it'll helpful to use cases process the app's input. In this case it was created implementation with Postgres client.

- **pkg**: This directory's responsible for all app's modules that hasn't fit with hexagonal components. In this case we have a module called [pkg/docker](https://github.com/guiferpa/bank/tree/main/pkg/docker).
