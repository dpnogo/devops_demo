package main

import "iam/internal/apiserver"

func main() {
	apiserver.NewApp(".keep-server").Run()
}
