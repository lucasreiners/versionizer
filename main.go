package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jessevdk/go-flags"
	"github.com/snabb/isoweek"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	_, err := flags.Parse(&CommandlineArguments)

	if err != nil {
		log.Fatal(err)
	}

	if CommandlineArguments.Type == "number" || CommandlineArguments.Type == "" {
		fmt.Println(versionNumber())
	} else if CommandlineArguments.Type == "code" {
		fmt.Println(versionCode())
	} else {
		log.Fatal("Unsupported output type. Only number|code allowed")
	}
}

func versionNumber() string {
	return fmt.Sprintf("%d.%d.%d", majorVersion(), minorVersion(), patchVersion())
}

func versionCode() int {
	val, _ := strconv.Atoi(fmt.Sprintf("%d%02d%03d", majorVersion(), minorVersion(), patchVersion()))
	return val
}

func majorVersion() int {
	return time.Now().Year()
}

func minorVersion() int {
	_, week := time.Now().ISOWeek()
	return week
}

func patchVersion() int {
	year, week := time.Now().ISOWeek()
	startOfWeek := isoweek.StartTime(year, week, time.UTC)

	var workingDirectory string

	if CommandlineArguments.Path != "" {
		workingDirectory = CommandlineArguments.Path
	} else {
		workingDirectory, _ = os.Getwd()
	}

	repo, err := git.PlainOpen(workingDirectory)

	if err != nil {
		log.Fatal(err)
	}

	_, err = repo.Head()

	if err != nil {
		log.Fatal(err)
	}

	numberOfCommits := 0

	gitLog, err := repo.Log(&git.LogOptions{Since: &startOfWeek})
	if err != nil {
		log.Fatal(err)
	}

	err = gitLog.ForEach(func(c *object.Commit) error {
		numberOfCommits++
		return nil
	})

	return numberOfCommits
}

var CommandlineArguments struct {
	Path string `short:"p" long:"path" description:"The path that contains the GIT repository"`
	Type string `short:"t" long:"type" description:"The type of version number to generate"`
}
