package main

import cmd "goclisandbox/cli"

func main() {
	// name := flag.String("name", "world", "The name to greet.")
	// test := flag.String("test", "toto", "this is a test")
	// flag.Parse()

	// fmt.Printf("Hello, %s!\n", *name)
	// fmt.Printf("test value is %s\n", *test)
	cmd.Execute()
}
