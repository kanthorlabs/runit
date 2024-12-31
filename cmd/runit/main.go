package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

//go:embed template.tpl
var templateContent string
var tpl = template.Must(template.New("Dockerfile").Parse(templateContent))

var version = "python:3.12-slim"

var image = "kanthorlab/runit-python"

var filename = "/home/tuannguyen/Projects/kanthorlabs/runit/examples/python/ip-checker.py"

var syspkgs = []string{
	"os",
}

func main() {
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	// Create a tar archive containing the project directory
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	info, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}

	// Prepare tar header
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
	}
	header.Name = "app.py"
	if err := tw.WriteHeader(header); err != nil {
		log.Fatalf("Error writing header: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	// Copy the file content to the tar archive
	if _, err := io.Copy(tw, file); err != nil {
		log.Fatalf("Error copying file content: %v", err)
	}

	tplBuffer := new(bytes.Buffer)
	if err := tpl.Execute(tplBuffer, struct{ Version string }{Version: version}); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	imports, err := extractImports(filename)
	if err != nil {
		log.Fatalf("Error extracting imports: %v", err)
	}
	importBuffer := new(bytes.Buffer)
	importBuffer.WriteString(strings.Join(imports, "\n"))

	importHeader := &tar.Header{
		Name: "requirements.txt",
		Mode: 0600,
		Size: int64(importBuffer.Len()),
	}
	if err := tw.WriteHeader(importHeader); err != nil {
		log.Fatalf("Error writing import header: %v", err)
	}
	if _, err := importBuffer.WriteTo(tw); err != nil {
		log.Fatalf("Error writing import content: %v", err)
	}

	dockerfileHeader := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(tplBuffer.Len()),
	}
	if err := tw.WriteHeader(dockerfileHeader); err != nil {
		log.Fatalf("Error writing Dockerfile header: %v", err)
	}
	if _, err := tplBuffer.WriteTo(tw); err != nil {
		log.Fatalf("Error writing Dockerfile content: %v", err)
	}

	if err := tw.Close(); err != nil {
		log.Fatalf("Error closing tar writer: %v", err)
	}

	build, err := cli.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			Tags:       []string{image},
			Dockerfile: "Dockerfile",
			Remove:     true,
		},
	)
	if err != nil {
		log.Fatalf("Error building image: %v", err)
	}
	defer build.Body.Close()

	// Read and display the build logs
	_, err = io.Copy(os.Stdout, build.Body)
	if err != nil {
		log.Fatalf("Error reading image build response: %v", err)
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("8080/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "8080",
					},
				},
			},
		},
		&network.NetworkingConfig{}, nil, "")
	if err != nil {
		log.Fatalf("Error creating container: %v", err)
	}
	containerId := resp.ID

	if err := cli.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		log.Fatalf("Error starting container: %v", err)
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("Error waiting for container: %v", err)
		}
	case <-statusCh:

	}

	out, err := cli.ContainerLogs(ctx, containerId, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Fatalf("Error retrieving logs: %v", err)
	}
	io.Copy(os.Stdout, out)

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerId, container.RemoveOptions{}); err != nil {
		log.Fatalf("Error removing container: %v", err)
	}
}

func extractImports(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imports []string
	importRegex := regexp.MustCompile(`(?m)^(?:from\s+([a-zA-Z0-9_\.]+)\s+import|import\s+([a-zA-Z0-9_\.]+))`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := importRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			if matches[1] != "" && !slices.Contains(syspkgs, matches[1]) {
				imports = append(imports, matches[1])
			} else if matches[2] != "" && !slices.Contains(syspkgs, matches[2]) {
				imports = append(imports, matches[2])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return imports, nil
}
