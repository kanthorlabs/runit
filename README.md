# runit
> Run arbitrary python script you have, no setup, no configuration, just run it

## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction
runit is a simple tool that allows you to run arbitrary Python scripts without any setup or configuration. It is designed to be lightweight and easy to use, making it perfect for quick tasks and prototyping.

## Installation
To install runit, you need to have Go (Golang) installed on your system. Follow these steps:

1. **Install Go**: If you haven't installed Go yet, you can download it from the official Go website: [https://golang.org/dl/](https://golang.org/dl/). Follow the installation instructions for your operating system.

2. **Set up your Go workspace**: Make sure your Go workspace is set up correctly. You can set the `GOPATH` environment variable to your desired workspace directory. For example, you can add the following lines to your `.bashrc` or `.zshrc` file:
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```

3. **Install runit**: Open your terminal and run the following command to install runit:
   ```bash
   go install github.com/kanthorlabs/runit/cmd/runit@latest
   ```

4. **Verify the installation**: After the installation is complete, you can verify that runit is installed correctly by running:
   ```bash
   runit --version
   ```

Now you are ready to use runit to execute your Python scripts!

## Usage
runit provides a simple command-line interface to run Python scripts. Here's the basic syntax:

```bash
runit [options] <path-to-your-python-script>
```

### Command Options

| Option             | Description                                    | Default            | Example                               |
|--------------------|------------------------------------------------|--------------------|---------------------------------------|
| `--platform-version` | Specify the Python Docker image version      | `python:3.13-slim` | `--platform-version python:3.11-slim`   |
| `--ports`          | Expose ports (can be specified multiple times) |                    | `--ports 3000 --ports 8000`           |

### Examples

1. Run a script with default settings:
```bash
runit script.py
```

2. Run with a specific Python version:
```bash
runit --platform-version python:3.11-slim script.py
```

3. Expose specific ports:
```bash
runit --ports 3000 script.py
```

4. Multiple port exposure:
```bash
runit --ports 3000 --ports 8000 script.py
```

5. Combine multiple options:
```bash
runit --platform-version python:3.11-slim --ports 3000 --ports 8000 script.py
```

### Docker Images
runit uses official Python Docker images. You can specify any valid Python image tag from Docker Hub. Some common options:

- `python:3.13-slim` (default)
- `python:3.11-slim`
- `python:3.12-slim`
- `python:3.11-alpine`
- `python:3.12-alpine`

### Port Binding
When you specify ports using the `--ports` flag, runit will:
1. Expose these ports in the Docker container
2. Map them to the same port numbers on your host machine
3. Make them accessible via localhost/127.0.0.1

## Contributing
Contributions to runit are welcome! If you have any ideas, suggestions, or bug reports, please feel free to open an issue or submit a pull request.

## License
runit is open-sourced under the MIT License - see the [LICENSE](LICENSE) file for details.

