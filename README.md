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
To use runit, you can run the following command:
```bash
runit <path-to-your-python-script>
```

runit will build a Docker image, run the script, and then remove the container.

## Contributing
Contributions to runit are welcome! If you have any ideas, suggestions, or bug reports, please feel free to open an issue or submit a pull request.

## License
runit is open-sourced under the MIT License - see the [LICENSE](LICENSE) file for details.

    