# Nougat

Project Nougat is an agnostic component used at Chamaconekt Kenya to help us create and send Golang http API requests. It abstracts the pain of reimplementing logic common to all clients, by storing HTTP Request properties to simplify sending requests and decoding responses.

## Features

- Method Setters: Get/Post/Put/Patch/Delete/Head
- Add or Set Request Headers
- Base/Path: Extend a Nougat for different endpoints
- Encode structs into URL query parameters
- Encode a form or JSON into the Request Body
- Receive JSON success or failure responses

## Install

```bash
go get github.com/wondenge/nougat
```
