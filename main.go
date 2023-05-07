package main

import (
	"fmt"
	composeLoader "github.com/compose-spec/compose-go/loader"
	composeTypes "github.com/compose-spec/compose-go/types"
	cliCommand "github.com/docker/cli/cli/command"
	composeApi "github.com/docker/compose/v2/pkg/api"
	composeCompose "github.com/docker/compose/v2/pkg/compose"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
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

func restartContainers(workingDir string, composeFilePath string, env map[string]string) error {
	ctx := context.Background()
	project, err := composeLoader.Load(composeTypes.ConfigDetails{ // Load the compose file
		//Version:     composeVersion,
		WorkingDir:  workingDir,
		ConfigFiles: []composeTypes.ConfigFile{{Filename: composeFilePath}},
		Environment: env,
	})
	if err != nil {
		return err
	}

	cli, _ := cliCommand.NewDockerCli()
	minute := time.Minute
	composeService := composeCompose.NewComposeService(cli)

	// TODO: split into stop and remove, check if containers are running
	// TODO: make exclusion for specific container (need to split)
	err = composeService.Down(ctx, project.Name, composeApi.DownOptions{ // Stop and remove containers
		RemoveOrphans: true,
		Project:       project,
		Timeout:       &minute,
		Images:        "local",
		Volumes:       false,
	})
	if err != nil {
		return err
	}

	var services []string
	for i, service := range project.Services {
		if service.Build == nil {
			continue
		}
		service.PullPolicy = composeTypes.PullPolicyBuild
		project.Services[i] = service
		services = append(services, service.Name)
	}
	err = composeService.Create(ctx, project, composeApi.CreateOptions{
		Services: services,
	})
	if err != nil {
		return err
	}

	err = composeService.Start(ctx, project.Name, composeApi.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	err = performPull("/repo", "/.ssh/id_rsa")
	if err != nil {
		log.Panic(err)
	}

	envVarsMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		split := strings.SplitN(envVar, "=", 2)
		envVarsMap[split[0]] = split[1]
	}
	err = restartContainers("/repo", "/repo/compose.yml", envVarsMap)
	if err != nil {
		log.Panic(err)
	}
}
