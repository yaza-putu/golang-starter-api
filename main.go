package main

import "yaza/src/core"

func main() {
	// load env
	core.Env()

	// init database
	core.Database()

	// start server
	core.HttpServe()
}
