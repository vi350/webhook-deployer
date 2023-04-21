package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	// Open the repository
	repo, err := git.PlainOpen(os.Getenv("REPO_PATH"))
	if err != nil {
		log.Panic(err)
	}

	// Set up the SSH authentication
	auth, err := ssh.NewPublicKeysFromFile("git", os.Getenv("PRIVATE_KEY_PATH"), "")
	if err != nil {
		log.Panic(err)
	}

	// Get the worktree of the repository
	worktree, err := repo.Worktree()
	if err != nil {
		log.Panic(err)
	}

	// Pull the latest changes from the origin repository using key-based authentication
	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Pull successful")
}
