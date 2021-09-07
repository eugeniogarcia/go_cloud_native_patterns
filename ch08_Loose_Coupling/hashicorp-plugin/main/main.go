package main

import (
	"fmt"
	"os"
	"os/exec"

	"egsmartin.com/hashicorp-plugin/commons"
	"github.com/hashicorp/go-plugin"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: run main/main.go animal")
		os.Exit(1)
	}

	// Get the animal name, and build the path where we expect to
	// find the corresponding executable file.
	name := os.Args[1]
	module := fmt.Sprintf("./%s/%s", name, name)

	// Does the file exist?
	_, err := os.Stat(module)
	if os.IsNotExist(err) {
		fmt.Println("can't find an animal named", name)
		os.Exit(1)
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"sayer": &commons.SayerPlugin{},
	}

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: commons.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(module),
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("sayer")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// We should have a Sayer now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	sayer := raw.(commons.Sayer)

	// Now we can use our loaded plugin!
	fmt.Printf("A %s says: %q\n", name, sayer.Says())
}
