package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	"os"

	"github.com/uthng/tf-inventory/inventory"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("tf-inventory", "An utility to generate ansible inventory file from Terraform state")
	list    = app.Flag("list", "List mode").Bool()
	host    = app.Flag("host", "Host mode").String()
	tfsFile = app.Arg("tfsfile", "Terraform state file").String()
)

func main() {

	_, err := app.Parse(os.Args[1:])
	app.FatalIfError(err, "")

	if *list && (*host != "") {
		app.Fatalf("The flag --list cannot be used with --host")
	}

	if *list {
		res, err := inventory.ExecCmd("list", "", *tfsFile)
		app.FatalIfError(err, "")

		b, err := json.Marshal(res)
		app.FatalIfError(err, "")

		fmt.Fprintf(os.Stdout, string(b))
		app.Terminate(nil)
	}

	if *host != "" {
		res, err := inventory.ExecCmd("host", *host, *tfsFile)
		app.FatalIfError(err, "")

		b, err := json.Marshal(res)
		app.FatalIfError(err, "")

		fmt.Fprintf(os.Stdout, string(b))
		app.Terminate(nil)
	}

	//if *inventory {
	//fmt.Println("Flag --inventory")
	//}
}
