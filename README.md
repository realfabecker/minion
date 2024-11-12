# kevin

## Introduction

Kevin is a dynamic command creation tool that allows you to define and execute commands based on a configuration file.
It is designed to simplify the process of running complex scripts and commands.

## Features

- Dynamic command creation
- Configuration via `kevin.yml`
- Support for flags and arguments
- Easy integration with shell scripts

## Installation

You can install or update Kevin using the installation script:

```bash
curl -so- https://raw.githubusercontent.com/realfabecker/kevin/master/install.sh | bash
```

## Usage

Here is an example of how to define and use a command with Kevin in the kevin.yml file:

```yaml
commands:
  - name: "create"
    parent: "gpg"
    short: "create a new gpg key"
    args:
      - name: kid
        usage: "gpg key name"
        required: true
    cmd: |
      echo 'Generating a new gpg key for {{ .GetArg "kid" }}'
      config="Key-Type: 1\n"
      config+="Key-Length: 2048\n"
      config+="Subkey-Type: 1\n"
      config+="Subkey-Length: 2048\n"
      config+="Name-Real: $(echo {{ .GetArg "kid" }} | sed -E 's/(.*)@.*/\1/g')\n"
      config+="Name-Email: {{ .GetArg "kid" }}\n"
      config+="Expire-Date: $(date --iso-8601=s -d '+5weeks' | tr -d ':-' | cut -c 1-15)"
      echo -e "$config" > /tmp/gpg-key.conf
      gpg --batch --gen-key /tmp/gpg-key.conf        
```

The kevin.yml configuration file can be stored globally in the user's home directory, or specifically by creating a file in the same directory as the invocation of the kevin command.

With this file ready, it will be possible to call the custom command as follows:

```bash
kevin gpg create
```

## Contributing

We welcome contributions! Please refer to the [contribution guide](./docs/CONTRIBUTING.md) for details on how to
contribute to the project.

## Versioning

This project uses [SemVer](https://semver.org/) for versioning. For all available versions, check the
[tags in this repository](https://github.com/realfabecker/kevin/tags).

## Licence

This project is licensed under the MIT License. See the [License](LICENSE.md) for more information.
