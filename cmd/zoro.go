package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yaza-putu/golang-starter-api/internal/core"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	_ "github.com/yaza-putu/golang-starter-api/internal/database/migrations"
	_ "github.com/yaza-putu/golang-starter-api/internal/database/seeders"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"github.com/yaza-putu/golang-starter-api/pkg/unique"
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
			"- migrate:down",
			"- make:seeder",
			"- seed:up",
			"- configure",
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
			case "configure:module":
				command.configure()
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
		configure() bool
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
	fName := fmt.Sprintf("./internal/database/migrations/%s_create_table_%s.go", time.Now().Format("20060102150405"), os.Args[2])

	// from template
	from, err := os.Open("./internal/database/migration.stub")
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
	fName := fmt.Sprintf("./internal/database/seeders/%s_create_%s_seeder.go", time.Now().Format("20060102150405"), os.Args[2])

	// from template
	from, err := os.Open("./internal/database/seeder.stub")
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
	token := unique.Key(51)
	refresh := unique.Key(51)
	passphrase := unique.Key(32)

	findAndReplaceByKey("key_token", fmt.Sprintf("'%s'", token))
	findAndReplaceByKey("key_refresh", fmt.Sprintf("'%s'", refresh))
	findAndReplaceByKey("key_passphrase", fmt.Sprintf("'%s'", passphrase))

	fmt.Println("Generate key successfully")
	return true
}

func (z *zoroCommand) configure() bool {

	if len(os.Args) != 3 {
		fmt.Println("ex : configure:module module-name")
		return false
	}

	baseDir := "../." // root directory
	oldModule := "github.com/yaza-putu/golang-starter-api"
	newModule := os.Args[2]

	fmt.Printf("Processing all .go files in directory: %s\n", baseDir)
	err := replaceInAllDirectories(baseDir, oldModule, newModule)
	if err != nil {
		fmt.Printf("Error processing directories: %v\n", err)
		return false
	}

	fmt.Println("Configure module successfully.")
	return true
}

func findAndReplaceByKey(key, newValue string) error {
	filename := ".env"
	// Read the entire file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("file .env not found")
		return err
	}

	// Split content by lines
	lines := strings.Split(string(content), "\n")

	// Find and replace the key if found
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + newValue
			found = true
			break
		}
	}

	// If key not found, return an error
	if !found {
		return fmt.Errorf("key '%s' not found in file", key)
	}

	// Write the modified content back to the file
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

func replaceInFile(filePath string, oldString string, newString string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	output := strings.ReplaceAll(string(input), oldString, newString)

	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

func replaceInAllDirectories(baseDir string, oldString string, newString string) error {
	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), ".stub") || strings.HasSuffix(info.Name(), ".mod") {
				err = replaceInFile(path, oldString, newString)
				if err != nil {
					fmt.Printf("Failed to replace in file %s: %v\n", path, err)
				}
			}
		}

		return nil
	})
}
