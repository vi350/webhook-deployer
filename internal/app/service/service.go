package service

import (
	composeLoader "github.com/compose-spec/compose-go/loader"
	composeTypes "github.com/compose-spec/compose-go/types"
	cliCommand "github.com/docker/cli/cli/command"
	composeApi "github.com/docker/compose/v2/pkg/api"
	composeCompose "github.com/docker/compose/v2/pkg/compose"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/net/context"
	"os"
	"strings"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) PerformPull(repoPath string, privateKeyPath string) error {
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
	return nil
}

func (s *Service) RestartContainers(workingDir string, composeFilePath string) error {
	ctx := context.Background()

	project, err := composeLoader.Load(composeTypes.ConfigDetails{ // Load the compose file
		//Version:     composeVersion,
		WorkingDir:  workingDir,
		ConfigFiles: []composeTypes.ConfigFile{{Filename: composeFilePath}},
		Environment: s.makeEnvMap(),
	})
	if err != nil {
		return err
	}

	cli, _ := cliCommand.NewDockerCli()
	composeService := composeCompose.NewComposeService(cli)

	var services []string
	for _, service := range project.Services {
		if !strings.HasSuffix(service.Name, "deployer") { // we don't want to stop the deployer container
			services = append(services, service.Name)
		}
	}
	err = composeService.Remove(ctx, project.Name, composeApi.RemoveOptions{ // Remove (actually includes stop) containers
		Project:  project,
		Stop:     true,
		Services: services,
	})
	if err != nil {
		return err
	}

	services = make([]string, 0) // clear slice
	for i, service := range project.Services {
		if !strings.HasSuffix(service.Name, "deployer") { // we don't want to create the deployer container
			if service.Build == nil {
				continue
			}
			service.PullPolicy = composeTypes.PullPolicyBuild
			project.Services[i] = service
			services = append(services, service.Name)
		}
	}
	err = composeService.Create(ctx, project, composeApi.CreateOptions{
		Services: services,
	})
	if err != nil {
		return err
	}

	err = composeService.Start(ctx, project.Name, composeApi.StartOptions{
		Project:  project,
		Services: services,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) makeEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		split := strings.SplitN(envVar, "=", 2)
		envMap[split[0]] = split[1]
	}
	return envMap
}
