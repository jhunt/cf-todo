package main

type File struct {
	b []byte
	t string
}

var files map[string]File
