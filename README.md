# Goil : Go iseplive.fr (old version) interface 

[![Build Status](https://travis-ci.org/aabizri/goil.svg?branch=dev)](https://travis-ci.org/aabizri/goil) [![Go Report Card](https://goreportcard.com/badge/github.com/aabizri/goil)](https://goreportcard.com/report/github.com/aabizri/goil) [![GoDoc](https://godoc.org/github.com/aabizri/goil?status.svg)](https://godoc.org/github.com/aabizri/goil)

# WARNING: OBSOLETE. THIS DOESN'T WORK ON THE NEW WEBSITE. NOT MAINTAINED.

A library interface for iseplive.fr, a student website.
Warning: Still in development !

```
go get github.com/aabizri/goil
```

Documentation for the master branch is on [godoc (https://godoc.org/github.com/aabizri/goil)](https://godoc.org/github.com/aabizri/goil) !

## Getting started:

### Logging in

Start a session
```golang
import "net/http"
session,err := goil.Login(USERNAME,PASSWORD, &http.Client{})
```

### Students

Get the list of all students alors with their data 
```golang
studentList,err := session.GetStudentList()
```

Export it in csv format through whatever io.writer you prefer, here with `os.Stdout`
```golang
err = studentList.WriteCSV(os.Stdout)
```

### Publications

Publish something
```golang
publication := CreatePublication("Hello World", goil.Divers)
err := session.PostPublication(post)
```

Delete a publication
```golang
err := session.DeletePublication(1750)
```
