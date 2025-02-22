package dockerx

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/kanthorlabs/runit/runtime/pythonx"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

func Exec(filepath string, vars *pythonx.DockerfileVars) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := BuildLockfile(tw, file); err != nil {
		return errors.Wrap(err, "failed to build lockfile")
	}

	if err := BuildApplication(tw, file); err != nil {
		return errors.Wrap(err, "failed to build application")
	}

	if err := BuildDockerfile(tw, vars); err != nil {
		return errors.Wrap(err, "failed to build Dockerfile")
	}

	if err := tw.Close(); err != nil {
		return err
	}

	// Create a Docker client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}

	name, err := BuildName(filepath, vars, file)
	if err != nil {
		return err
	}

	if err := BuildImage(cli, buf, name); err != nil {
		return err
	}

	return RunContainer(cli, name, vars)
}

func BuildLockfile(tw *tar.Writer, file *os.File) error {
	// Reset the file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	maps, err := pythonx.Scan(scanner, pythonx.PackageSystem)
	if err != nil {
		return err
	}
	lockfile := pythonx.Lockfile(maps)

	header := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "requirements.txt",
		Size:     int64(lockfile.Len()),
		Mode:     0755,
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, lockfile); err != nil {
		return err
	}

	return nil
}

func BuildApplication(tw *tar.Writer, file *os.File) error {
	// Reset the file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "main.py",
		Size:     stat.Size(),
		Mode:     0755,
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}

func BuildDockerfile(tw *tar.Writer, vars *pythonx.DockerfileVars) error {
	tpl, err := pythonx.Dockerfile(vars)
	if err != nil {
		return err
	}

	header := &tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "Dockerfile",
		Size:     int64(tpl.Len()),
		Mode:     0755,
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, tpl); err != nil {
		return err
	}

	return nil
}

func BuildName(filepath string, vars *pythonx.DockerfileVars, file *os.File) (string, error) {
	// Reset the file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	hash := sha256.New()

	// Write filepath to hash
	if _, err := hash.Write([]byte(filepath)); err != nil {
		return "", err
	}

	// Write file content to hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Write vars to hash
	varsString := fmt.Sprintf("%s-%s-%s-%s",
		vars.Version,
		vars.Arguments,
		vars.Params,
		vars.Ports,
	)
	if _, err := hash.Write([]byte(varsString)); err != nil {
		return "", err
	}

	sha256sum := fmt.Sprintf("%x", hash.Sum(nil))
	name := fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), sha256sum[0:6])
	return name, nil
}

func BuildImage(cli *client.Client, buf *bytes.Buffer, name string) error {
	ctx := context.Background()
	buildops := types.ImageBuildOptions{
		Tags:       []string{fmt.Sprintf("kanthorlab/runit-python:%s", name)},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}
	build, err := cli.ImageBuild(ctx, buf, buildops)
	if err != nil {
		return err
	}
	defer build.Body.Close()

	_, err = io.Copy(os.Stdout, build.Body)
	if err != nil {
		return err
	}

	return nil
}

func RunContainer(cli *client.Client, name string, vars *pythonx.DockerfileVars) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Build command array
	cmd := []string{"python", "main.py"}
	if vars.Arguments != "" {
		cmd = append(cmd, vars.Arguments)
	}
	if vars.Params != "" {
		cmd = append(cmd, vars.Params)
	}

	conf := &container.Config{
		Image: fmt.Sprintf("kanthorlab/runit-python:%s", name),
		Cmd:   cmd,
	}
	hostconf := &container.HostConfig{
		PortBindings: nat.PortMap{},
	}
	networkconf := &network.NetworkingConfig{}
	for _, port := range vars.Ports {
		hostconf.PortBindings[nat.Port(port+"/tcp")] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port,
			},
		}
	}
	platformconf := &ocispec.Platform{}

	cont, err := cli.ContainerCreate(ctx, conf, hostconf, networkconf, platformconf, name)
	if err != nil {
		return err
	}

	defer func() {
		if err := cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{}); err != nil {
			log.Fatal(err)
		}
	}()

	if err := cli.ContainerStart(ctx, cont.ID, container.StartOptions{}); err != nil {
		return err
	}

	waitc, errc := cli.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errc:
		if err != nil {
			return err
		}
	case waitresp := <-waitc:
		if err := waitresp.Error; err != nil {
			return fmt.Errorf("exit code: %d - %s", waitresp.StatusCode, err.Message)
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	out, err := cli.ContainerLogs(ctx, cont.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		log.Fatalf("Error retrieving logs: %v", err)
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}
