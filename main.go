package main

import (
	"log"
	"os/exec"
)

func main() {
	exec.Command("mdkir", "-p", "jarfiles", "worlds", "users")
	var serverProperties ServerProperties
	_, err := MakeServer("1.16.4", "TestServer", true, serverProperties)

	if err != nil {
		log.Fatalln(err)
	}
}
