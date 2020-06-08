package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ContainerRunParams struct {
	AWSAccess string
	AWSSecret string
	Fullname  string
	Branch    string
	Commit    string
	Bucket    string
	Repo      string
	User      string
	CDN       string
	Token     string
	buildID   string
	PrID      string
}

const (
	awsAccess string = "AWS_ACCESS_KEY_ID"
	awsSecret string = "AWS_SECRET_ACCESS_KEY"
	repoFull  string = "REPO_FULL"
	branch    string = "STACKERS_BRANCH"
	commit    string = "STACKERS_COMMIT"
	bucket    string = "BUCKET"
	repo      string = "REPO"
	user      string = "USER"
	cdn       string = "CDN"
	token     string = "OAUTH_TOKEN"
	prId      string = "STACKERS_PR_ID"
)

const (
	buildImageName = "slootsantos/own"
	buildImageHub  = "docker.io/slootsantos/own"
)

const (
	bindingContainerPort = "80"
	bindingHostPort      = "8000"
	bindingProtocol      = "tcp"
	bindingHost          = "0.0.0.0"
)

func BuildRepoSources(params ContainerRunParams) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cont := createRunnableContainer(cli, params)

	err = executeContainer(cli, params.buildID, cont)
	if err != nil {
		fmt.Println(err)
	}
}

func executeContainer(cli *client.Client, buildID string, cont container.ContainerCreateCreatedBody) error {
	log.Println("START: executing Docker container")

	err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	streamChan, _ := CreateBuild(buildID)
	streamDockerLogs(cli, cont.ID, buildID, streamChan)

	log.Println("DONE: executing Docker container")
	return err
}

func streamDockerLogs(cli *client.Client, containerID string, buildID string, streamChan chan<- string) error {
	go StreamLogs(buildID)

	logConfig := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	out, err := cli.ContainerLogs(context.Background(), containerID, logConfig)
	defer out.Close()
	if err != nil {
		return err
	}

	defer close(streamChan)
	for {
		buf := make([]byte, 32*1024+9)

		_, err := out.Read(buf)
		if err != nil {
			if err == io.EOF {
				streamChan <- StreamingCloseMessage
				break
			}

			return err
		}

		streamChan <- fmt.Sprintf("%s", buf)
	}

	_, err = cli.ContainerWait(context.Background(), containerID)
	return err
}

func createRunnableContainer(cli *client.Client, params ContainerRunParams) container.ContainerCreateCreatedBody {
	log.Println("START: pulling Docker container")

	reader, err := cli.ImagePull(context.Background(), buildImageHub, types.ImagePullOptions{})
	if err != nil {
		panic(errors.New("could not pull container" + err.Error()))
	}

	io.Copy(os.Stdout, reader)

	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: buildImageName,
			Env:   createContainerEnv(params),
			Tty:   true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				PidsLimit: 150,
			},
			PortBindings: createPortBinding(),
		}, nil, "",
	)

	if err != nil {
		panic(err)
	}

	log.Println("DONE: pulling Docker container")
	return cont
}

func createPortBinding() nat.PortMap {
	hostBinding := nat.PortBinding{
		HostIP:   bindingHost,
		HostPort: bindingHostPort,
	}

	containerPort, err := nat.NewPort(bindingProtocol, bindingContainerPort)
	if err != nil {
		panic(err)
	}

	return nat.PortMap{
		containerPort: []nat.PortBinding{hostBinding},
	}
}

func createContainerEnv(params ContainerRunParams) []string {
	return []string{
		joinEnv(awsAccess, params.AWSAccess),
		joinEnv(awsSecret, params.AWSSecret),
		joinEnv(bucket, params.Bucket),
		joinEnv(branch, params.Branch),
		joinEnv(commit, params.Commit),
		joinEnv(repo, params.Repo),
		joinEnv(user, params.User),
		joinEnv(cdn, params.CDN),
		joinEnv(repoFull, params.User+"/"+params.Repo),
		joinEnv(token, params.Token),
		joinEnv(prId, params.PrID),
	}
}

func joinEnv(key string, value string) string {
	return key + "=" + value
}
