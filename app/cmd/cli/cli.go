package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/arek-e/D7024E/app/internal"
	"github.com/manifoldco/promptui"
)

type CLI struct {
	Node *internal.Kademlia
	Net  *internal.Network
}

var port = 1337

// StartCLI initializes and starts the interactive CLI.
func (cli *CLI) StartCLI(exitCh chan<- struct{}) {
	fmt.Println("\n======Kadlab node CLI========")
	fmt.Println("Available Commands: ping, put (p), get (g), exit (q)")

	for {
		prompt := promptui.Prompt{
			Label: "Enter Command:",
		}

		input, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		// Split the entered input into command and arguments.
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "ping":
			if len(parts) > 1 {
				cli.pingCmd(parts[1]) // Use the IP address provided as an argument
			} else {
				prompt := promptui.Prompt{
					Label: "Enter ping address (IP):",
				}
				ipAddress, err := prompt.Run()
				if err != nil {
					log.Fatal(err)
				}
				cli.pingCmd(ipAddress)
			}
		case "put", "p":
			if len(parts) > 1 {
				cli.putCmd(parts[1])
			} else {
				prompt := promptui.Prompt{
					Label: "What do you want to store?",
				}
				dataToStore, err := prompt.Run()
				if err != nil {
					log.Fatal(err)
				}
				cli.putCmd(dataToStore)
			}
		case "get", "g":
			if len(parts) > 1 {
				cli.getCmd(parts[1])
			} else {
				prompt := promptui.Prompt{
					Label: "Insert hash",
				}
				hash, err := prompt.Run()
				if err != nil {
					log.Fatal(err)
				}
				cli.getCmd(hash)
			}
		case "exit", "q":
			fmt.Println("Exiting the CLI...")
			exitCh <- struct{}{}
		default:
			fmt.Println("Command not recognized. Available Commands: ping, put, get, exit")
		}
	}
}

func (cli *CLI) executeCommand(command string, exitCh chan<- struct{}) {
	switch command {
	case "ping":
		if len(os.Args) > 2 {
			cli.pingCmd(os.Args[2]) // Use os.Args[2] as the IP address
		} else {
			prompt := promptui.Prompt{
				Label: "Enter ping address (IP):",
			}
			ipAddress, err := prompt.Run()
			if err != nil {
				log.Fatal(err)
			}
			cli.pingCmd(ipAddress)
		}
	case "put", "p":
		if len(os.Args) > 2 {
			cli.putCmd(os.Args[2])
		} else {
			prompt := promptui.Prompt{
				Label: "What do you want to store?",
			}
			dataToStore, err := prompt.Run()
			if err != nil {
				log.Fatal(err)
			}
			cli.putCmd(dataToStore)
		}
	case "get", "g":
		if len(os.Args) > 2 {
			cli.getCmd(os.Args[2])
		} else {
			prompt := promptui.Prompt{
				Label: "Insert hash",
			}
			hash, err := prompt.Run()
			if err != nil {
				log.Fatal(err)
			}
			cli.getCmd(hash)
		}
	case "exit", "q":
		fmt.Println("Exiting the CLI...")
		exitCh <- struct{}{}
	default:
		fmt.Println("Command not recognized. Available Commands: ping, put, get, exit")
	}
}

func (cli *CLI) pingCmd(ipAddress string) {
	fmt.Printf("Starting to ping: %s\n", ipAddress)
	contact := internal.Contact{
		Address: ipAddress + ":" + strconv.Itoa(port),
	}

	_, err := cli.Net.SendPingMessage(&contact)
	if err != nil {
		fmt.Errorf("ERROR: %v", err)
	}
}

func (cli *CLI) putCmd(dataToStore string) {

}

func (cli *CLI) getCmd(hash string) {

}
