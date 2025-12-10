<br/>

<div align="center">
  <h3>Genesis</h3>
  <h4>A generic JSON api for small, private frontend apps</h4>
</div>

<div align="center">
  <a href="https://goreportcard.com/report/github.com/simonwep/genesis">
    <img src="https://goreportcard.com/badge/github.com/simonwep/genesis" alt="Go Report Card">
  </a>
  <a href="https://github.com/simonwep/genesis/actions/workflows/main.yml">
    <img src="https://github.com/simonwep/genesis/actions/workflows/main.yml/badge.svg" alt="CI Status">
  </a>
</div>

<br/>

### Summary

This project is designed specifically for small, personal projects requiring a straightforward, simple storage API that you can host yourself, including simplified user management.

### Usage

First, create a [.env](.env.example) and specify the initial usernames and passwords for access.
Make sure to fill out `GENESIS_JWT_SECRET` with a secure, random string, for that you can use `openssl rand -hex 32`.
You can specify the remaining values, but the defaults are good for medium-sized projects such as [ocular](https://github.com/Simonwep/ocular).

Second, start the server via `go run . start` - That's it.
Head to the [api](#api) documentation to see how to use it.
Use `go run . help` to see all available commands.

The `json` is pre-processed by the [minify](https://github.com/tdewolff/minify) package to minimize and validate it.

#### Using docker

You can run genesis using [docker](https://www.docker.com/products/docker-desktop/) by using pre-build images:

```sh
docker run -p 8080:8080 -v "$(pwd)/.data:/app/.data" --env-file .env ghcr.io/simonwep/genesis:latest start
```

Genesis should then be accessible under port `8080`.

> [!NOTE]
> You can specify the base-url via the env variable `GENESIS_BASE_URL`.

### CLI

Genesis comes with a CLI to manage users.
You can access it by running `go run . users help` or via docker using the following command:

```sh
docker run --rm -v "$(pwd)/.data:/app/.data" --env-file .env ghcr.io/simonwep/genesis:latest help
```

### API Documentation

Genesis includes interactive API documentation powered by Swagger/OpenAPI 3.0.

#### Accessing Swagger UI

Once the server is running, you can access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- Complete endpoint documentation for all 11 API endpoints
- Request/response schemas with examples
- Interactive testing capabilities
- Authentication information (cookie-based JWT)

#### Generating Swagger Documentation

If you modify the API or add new endpoints, regenerate the Swagger documentation:

```sh
# Install swag CLI tool (one-time setup)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g routes/setup.go --output docs
```

#### Disabling Swagger UI

To disable Swagger UI in production, set the environment variable:

```sh
GENESIS_SWAGGER_ENABLED=false
```

By default, Swagger is enabled.

### API

The API is kept as simple as possible; there is nothing more than user, data, and account management.

#### Authentication and account

* `POST /login` - Authenticates a user.
  - Takes either a `user` and `password` as JSON object and returns the user-data and a session cookie or, if a session-cookie exists, the current user.
  - Returns `401` the password is invalid or the user doesn't exist.
* `POST /logout` - Invalidates the current refresh token and logs out a user.
* `POST /account/update`
  - Takes a `newPassword` and `currentPassword` as JSON object.
  - Returns `200` if the password was successfully updated, otherwise `400`.

> [!NOTE]
> The JWT token is returned as a strict same-site, secure and http-only cookie!  
> When changing the password, the new password must fulfill the same requirements for adding a new user.

#### Data endpoints

* `GET /data` - Retrieves all data from the current user as object.
* `GET /data/:key` - Retrieves the data stored for the given `key`. Returns `204` if there is no content.
* `POST /data/:key` - Stores / overrides the data for `key`.
* `DELETE /data/:key` - Removes the data for `key`, always returns `200`, even if `key` doesn't exist.

> [!NOTE]
> Validation parameters for those endpoints are defined in [.env](.env.example).  
> This includes a key-pattern, the max amount per user, and a size-limit.

#### User management

> Admins can only use these endpoints!

* `GET /user` - Fetch all users as `{ name: string, admin: boolean }[]`.
* `POST /user` - Create a user, takes a JSON object with `user`, `password` and `admin` (all mandatory, `admin` is a boolean).
* `POST /user/:name` - Update a user by `name`, takes a JSON object with `password` and `admin` (both optional).
* `DELETE /user/:name` - Delete a user by `name`.

> [!NOTE]
> The username is validated against the pattern defined in [.env](.env.example).  
> The length must be between `3` and `32`, the password between `8` and `64`.
