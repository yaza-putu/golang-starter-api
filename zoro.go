package main

import (
	"fmt"
	"github.com/yaza-putu/golang-starter-api/src/core"
	"github.com/yaza-putu/golang-starter-api/src/database"
	_ "github.com/yaza-putu/golang-starter-api/src/database/migrations"
	_ "github.com/yaza-putu/golang-starter-api/src/database/seeders"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"io"
	"os"
	"time"
)

func main() {
	if os.Args[1] != "key:generate" {
		core.Env()
		core.Database()
	}

	command := New()

	if len(os.Args) == 1 {
		fmt.Println("Hello i'm zoro, can i help you ?")
		fmt.Println("Available command :")

		// migration collection
		m := []string{
			"- key:generate",
			"- make:migration",
			"- migrate:up",
			"- migration:down",
			"- make:seeder",
			"- seed:up",
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
			case "make:seeder":
				command.newSeeder()
				break
			case "seed:up":
				command.upSeeder()
				break
			case "key:generate":
				command.keyGenerate()
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
		newSeeder() bool
		upSeeder() bool
		keyGenerate() bool
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
	logger.New(err, logger.SetType(logger.FATAL))

	defer from.Close()

	// to file
	to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
	logger.New(err, logger.SetType(logger.FATAL))

	defer to.Close()
	// copy file with template
	_, err = io.Copy(to, from)
	logger.New(err, logger.SetType(logger.FATAL))

	fmt.Printf("New migration : %s\n", fName)

	return true
}

func (z *zoroCommand) upMigration() bool {
	err := database.MigrationUp()
	if err != nil {
		logger.New(err, logger.SetType(logger.FATAL))
		return false
	} else {
		fmt.Println("Migrating collections successfully")
		return true
	}
}

func (z *zoroCommand) downMigration() bool {
	err := database.MigrationDown()
	if err != nil {
		logger.New(err, logger.SetType(logger.FATAL))
		return false
	} else {
		fmt.Println("Drop collections successfully")
		return true
	}
}

func (z *zoroCommand) upSeeder() bool {
	err := database.SeederUp()

	if err != nil {
		logger.New(err, logger.SetType(logger.FATAL))
		return false
	} else {
		fmt.Println("Run seeders successfully")
		return true
	}
}

func (z *zoroCommand) newSeeder() bool {
	if len(os.Args) != 3 {
		fmt.Println("ex : make:seeder name-of-file")
		return false
	}

	// file
	fName := fmt.Sprintf("./src/database/seeders/%s_create_%s_seeder.go", time.Now().Format("20060102150405"), os.Args[2])

	// from template
	from, err := os.Open("./src/database/seeder.stub")
	logger.New(err, logger.SetType(logger.FATAL))
	defer from.Close()

	// to file
	to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
	logger.New(err, logger.SetType(logger.FATAL))

	defer to.Close()
	// copy file with template
	_, err = io.Copy(to, from)
	logger.New(err, logger.SetType(logger.FATAL))

	fmt.Printf("New seeder : %s\n", fName)

	return true
}
func (z *zoroCommand) keyGenerate() bool {
	token := utils.Key(51)
	refresh := utils.Key(51)

	fmt.Println("Generate key successfully")
	fmt.Println("Please copy bellow to .env")
	fmt.Println(fmt.Sprintf("key_token='%s' \nkey_refresh='%s'", token, refresh))
	return true
}
