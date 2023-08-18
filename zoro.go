package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"yaza/src/core"
	"yaza/src/database"
	_ "yaza/src/database/migrations"
)

func main() {
	core.Env()
	core.Database()

	command := New()

	if len(os.Args) == 1 {
		fmt.Println("Hello i'm zoro, can i help you ?")
		fmt.Println("Available command :")

		// migration collection
		m := []string{
			"- make:migration",
			"- migration:up",
			"- migration:down",
		}
		for _, v := range m {
			fmt.Println(v)
		}
		os.Exit(0)
	}

	for i, v := range os.Args {
		if i != 0 {
			switch v {
			case "make:migration":
				command.newMigration()
				break
			case "migrate:up":
				command.upMigration()
				break
			case "migrate:down":
				command.downMigration()
				break
			}
		}
	}
}

type (
	zoroCommand   struct{}
	zoroInterface interface {
		newMigration() bool
		upMigration() bool
		downMigration() bool
	}
)

func New() zoroInterface {
	return &zoroCommand{}
}

func (z *zoroCommand) newMigration() bool {
	if len(os.Args) != 3 {
		fmt.Println("ex : make:migration name-of-file")
		return false
	}

	// file
	fName := fmt.Sprintf("./src/database/migrations/%s_create_table_%s.go", time.Now().Format("20060102150405"), os.Args[2])

	// from template
	from, err := os.Open("./src/database/migration.stub")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	// to file
	to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer to.Close()
	// copy file with template
	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New migration : %s\n", fName)

	return true
}

func (z *zoroCommand) upMigration() bool {
	err := database.MigrationUp()
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		fmt.Println("Migrating collections successfully")
		return true
	}
}

func (z *zoroCommand) downMigration() bool {
	err := database.MigrationDown()
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		fmt.Println("Drop collections successfully")
		return true
	}
}
