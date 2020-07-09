<h1 align="center">Project Nougat</h1>

[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/wondenge/nougat?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/wondenge/nougat)](https://goreportcard.com/report/wondenge/nougat)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)]

Project Nougat is an agnostic component we actively use at [Chamaconekt Kenya](https://github.com/chamaconekt) to help us create and send [Golang](https://golang.org) http API requests. It abstracts reimplementing logic common to all clients, by storing HTTP request properties to simplify sending requests and decoding responses.

## Features

- **Method Setters:** Get/Post/Put/Patch/Delete/Head
- Add or Set Request Headers
- **Base/Path:** Extend a Nougat for different endpoints
- Encode structs into URL query parameters
- Encode a form or JSON into the Request Body
- Receive JSON success or failure responses

## Install

```bash
go get github.com/wondenge/nougat
```
