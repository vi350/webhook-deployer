package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func performPull(repoPath string, privateKeyPath string) error {
	repo, err := git.PlainOpen(repoPath) // Open the repository
	if err != nil {
		return err
	}
	auth, err := ssh.NewPublicKeysFromFile("git", privateKeyPath, "") // Create the authentication
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree() // Get the working directory for the repository
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{ // Pull changes from the remote
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		return err
	}
	fmt.Println("Pull successful")
	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	err = performPull(os.Getenv("REPO_PATH"), os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		log.Panic(err)
	}
}
