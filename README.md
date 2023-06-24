<br/>

<div align="center">
  <h3>Genesis</h3>
  <h4>A generic json api for small, private frontend apps (WIP)</h4>
</div>

<br/>

### Summary

This project is currently work-in-progress. _Anything_ can happen, hell, this might even disappear!
It is designed specifically for small, personal projects requiring a straightforward, simple storage API that you can host yourself.
It also requires you to specify a list of known-users, simplifying it greatly as there is no need for handling any sign-up process of new users.

### Usage

First, create a [.env](.env.example) and specify the usernames you want to have access to it.
You can specify the remaining values, but the defaults are good for medium-sized projects such as [ocular](https://github.com/Simonwep/ocular).

Second, start the server via `go run .`. That's it.
Head to the [api](#the-api) documentation to see how to use it.


### The API

The API is kept as simple as possible, there is nothing more than simple data-validation, json storage and user authentication.
It comes with the following endpoints (so far):


* `POST /register` - Registers a user initially.
  - Takes a `user` and `password` as json object.
  - Returns `201` on success, `401` if the user already exists or is invalid.
* `POST /login` - Authenticates a user via [JWT](https://jwt.io/).
  - Takes a `user` and `password` as json object.
  - Returns `200` on success including the token as header, `401` the password is invalid or the user is invalid.
* `GET /data` - Retrieves all data from a user as object.
* `GET /data/:key` - Retrieves the data stored for the given `key`.
* `PUT /data/:key` - Stores / overrides the data for `key`.
* `DELETE /data/:key` - Removes the data for `key`, always returns `200`, even if `key` doesn't exist.
