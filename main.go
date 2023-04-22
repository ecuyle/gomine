package main

import (
	"github.com/ecuyle/gomine/internal/servermanager"
	"log"
	"os/exec"
)

func main() {
	exec.Command("mdkir", "-p", "jarFiles", "worlds", "users")

	var serverPropertiesMap = make(map[string]interface{})
	serverPropertiesMap["difficulty"] = "peaceful"
	_, err := servermanager.MakeServer("1.19.4", "TestServer", true, serverPropertiesMap)

	if err != nil {
		log.Fatalln(err)
	}
}
