package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikaven1024/bolt-cli/cli"
	"github.com/ikaven1024/bolt-cli/db"
	"github.com/ikaven1024/bolt-cli/version"
)

var (
	pPath    = flag.String("path", "", "path of db file")
	pWeb     = flag.Bool("web", false, "support web if set")
	pVersion = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Parse()

	if *pVersion {
		fmt.Print(version.Version)
		os.Exit(0)
	}

	if len(*pPath) == 0 {
		fmt.Println("path is not set")
		os.Exit(1)
	}

	db, err := db.New(*pPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *pWeb {
		panic("Not Implement")
	} else {
		cli.Run(db)
	}
}
