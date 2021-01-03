package main

import (
	"log"
	"os/exec"
)

func main() {
	exec.Command("mdkir", "-p", "jarFiles", "worlds", "users")

	var serverPropertiesMap = make(map[string]interface{})
	serverPropertiesMap["difficulty"] = "peaceful"
	_, err := MakeServer("1.16.4", "TestServer", true, serverPropertiesMap)

	if err != nil {
		log.Fatalln(err)
	}
}
