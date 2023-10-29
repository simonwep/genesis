<br/>

<div align="center">
  <h3>Genesis</h3>
  <h4>A generic json api for small, private frontend apps (WIP)</h4>
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

This project is currently work-in-progress. _Anything_ can happen, hell, this might even disappear!
It is designed specifically for small, personal projects requiring a straightforward, simple storage API that you can host yourself.
It also requires you to specify a list of known-users, simplifying it greatly as there is no need for handling any sign-up process of new users.

### Usage

First, create a [.env](.env.example) and specify the initial usernames and passwords for access.
Make sure to fill out `GENESIS_JWT_SECRET` with a secure, random string, for that you can use `openssl rand -hex 32`.
You can specify the remaining values, but the defaults are good for medium-sized projects such as [ocular](https://github.com/Simonwep/ocular).

Second, start the server via `go run .`. That's it.
Head to the [api](#api) documentation to see how to use it.

The `json` is pre-processed by the [minify](https://github.com/tdewolff/minify) package to minimize and validate it.

#### Using docker

You can run genesis using [docker](https://www.docker.com/products/docker-desktop/) by using pre-build images:

```sh
docker run -p 8080:8080 -v "$(pwd)/.data:/app/.data" --env-file .env ghcr.io/simonwep/genesis:latest
```

Genesis should then be accessible under port `8080`.

### API

The API is kept as simple as possible, there is nothing more than user, data and account management.

#### Authentication and account

* `POST /login` - Authenticates a user.
  - Takes either a `user` and `password` as json object and returns the user-data and a session cookie or, if a session-cookie exists, the current user.
  - Returns `401` the password is invalid or the user doesn't exist.
* `POST /logout` - Invalidates the current refresh token and logs out a user.
* `POST /account/update`
  - Takes a `newPassword` and `currentPassword` as json object.
  - Returns `200` if the password was successfully updated, otherwise `400`.

> The JWT token is returned as strict same-site, secure and http-only cookie!  
> When changing the password, the new password must fulfill the same requirements for adding a new user.

#### Data endpoints

* `GET /data` - Retrieves all data from the current user as object.
* `GET /data/:key` - Retrieves the data stored for the given `key`. Returns `204` if there is no content.
* `POST /data/:key` - Stores / overrides the data for `key`.
* `DELETE /data/:key` - Removes the data for `key`, always returns `200`, even if `key` doesn't exist.

> Validation parameters for those endpoints are defined in [.env](.env.example).  
> This includes a key-pattern, the max amount per user and a size-limit.

#### User management

> These endpoints can only be used by admins!

* `GET /user` - Fetch all users as `{ name: string, admin: boolean }[]`.
* `POST /user` - Create a user, takes a json object with `user`, `password` and `admin` (all mandatory, `admin` is a boolean).
* `POST /user/:name` - Update a user by `name`, takes a json object with `password` and `admin` (both optional).
* `DELETE /user/:name` - Delete a user by `name`.

> The username is validated against the pattern defined in [.env](.env.example).  
> The length must be between `3` and `32`, the password between `8` and `64`.
