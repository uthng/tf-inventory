package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/uthng/tf-inventory/inventory"
)

// printHelp print command's help
func printHelp() {
	fmt.Println("")
	fmt.Printf("%s is a utility to generate ansible inventory file from Terraform state\n", os.Args[0])
	fmt.Printf("Syntax: %s [options] [hostname] [tfstate]\n", os.Args[0])
	fmt.Printf("\t-help: Print this help\n")
	fmt.Printf("\t-list: Generate ansible inventory file under JSON format in list mode\n")
	fmt.Printf("\t-host: Generate ansible inventory file under JSON format in host mode\n")
	fmt.Printf("\t-inventory: Generate ansible inventory file\n")
}

func main() {
	var res interface{}
	var err error

	tfsFile := ""

	help := flag.Bool("help", false, "Print help")
	list := flag.Bool("list", false, "List mode")
	host := flag.String("host", "", "Host mode")
	//inventory := flag.Bool("inventory", false, "Inventory mode")

	flag.Parse()

	// Check command line flags
	if flag.NFlag() <= 0 || flag.NFlag() > 1 {
		fmt.Fprintf(os.Stderr, "Error: you must specify one and only one argument")
		printHelp()
		os.Exit(1)
	}

	if *help {
		printHelp()
	}

	if *list {

		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "Error: only 1 arguments at most are accepted!\n")
			printHelp()
			os.Exit(1)
		}

		if len(os.Args) == 3 {
			tfsFile = os.Args[2]
		}

		res, err = inventory.ExecCmd("list", "", tfsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		b, err := json.Marshal(res)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, string(b))
		os.Exit(0)
	}

	if *host != "" {
		if len(os.Args) > 4 {
			fmt.Fprintf(os.Stderr, "Error: only 2 arguments at most are accepted!\n")
			printHelp()
			os.Exit(1)
		}

		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: You must specify at least a hostname as argument!\n")
			printHelp()
			os.Exit(1)
		}

		if len(os.Args) == 4 {
			tfsFile = os.Args[3]
		}

		res, err = inventory.ExecCmd("host", os.Args[2], tfsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		b, err := json.Marshal(res)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, string(b))
		os.Exit(0)
	}

	//if *inventory {
	//fmt.Println("Flag --inventory")
	//}
}
