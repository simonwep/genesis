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
You can specify the remaining values, but the defaults are good for medium-sized projects such as [ocular](https://github.com/Simonwep/ocular).

Second, start the server via `go run .`. That's it.
Head to the [api](#api) documentation to see how to use it.

The `json` is pre-processed by the [minify](https://github.com/tdewolff/minify) package to minimize and validate it.

#### Using docker

This API can also be deployed by using docker.
For this you can build and run the container using the following command:

```sh
docker build -t genesis .
docker run -p 8088:8080 -v "$(pwd)/.data:/app/.data" genesis
```

Genesis should then be accessible under port `8088`.

### API

The API is kept as simple as possible, there is nothing more than simple data-validation, json-storage and user-authentication.
It comes with the following endpoints (so far):

* `POST /login` - Authenticates a user via [JWT](https://jwt.io/).
  - Takes a `user` and `password` as json object.
  - Returns `{ expiresAt: number, token: string}` as body (`expiresAt` is a unix timestamp) and a refresh-token in the form of a cookie.
  - Returns `401` the password is invalid or the user doesn't exist.
* `GET /login/refresh` - Refreshes the access token and rotates the refresh token.
  - Returns `200` on success and `{ expiresAt: number, token: string}` as body. It also returns a new refresh-token in the form of a cookie.
  - Returns `401` if the refresh token is invalid / expired, or the user doesn't exist.
* `POST /logout` - Invalidates the current refresh token and logs out a user.
* `POST /account/update`
  - Takes a `newPassword` and `currentPassword` as json object.
  - Returns `200` if the password was successfully updated, otherwise `400`.
* `GET /data` - Retrieves all data from a user as object.
* `GET /data/:key` - Retrieves the data stored for the given `key`. Returns `204` if there is no content.
* `POST /data/:key` - Stores / overrides the data for `key`.
* `DELETE /data/:key` - Removes the data for `key`, always returns `200`, even if `key` doesn't exist.
