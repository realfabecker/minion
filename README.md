# kevin

kevin is a command-line tool designed to assist in the parallel execution of scripts.

## Usage

The simplest usage scenario involves executing a replicated script across multiple workers

```bash
kevin -w 10 -c "echo 'Hello, World!' $RANDOM"
```

It is also possible to use a CSV file as a basis, where each column will be passed as an argument to the script

```bash
kevin -w 10 -c "echo 'Hello, World!'" -f args.csv
```

## Installing and Updating

Installation or update of kevin can be done through its [installation script](./install.sh).

For this, you can either download and run it manually or use cURL as follows:

```bash
curl -so- https://raw.githubusercontent.com/realfabecker/kevin/master/install.sh | bash
```

The script above will download the latest stable release and extract it to the user's base directory

## Custom commands

kevin allows the creation of custom commands from the configuration file kevin.yml

```yaml
commands:
- name: "hello"    
  short: "say hello"
  cmd: |
    echo "Hello World"
  shell: "/bin/bash"
```

## Contributing

Refer to the [contribution guide](./docs/CONTRIBUTING.md) for details on how to contribute to the project.

## Versioning

This project uses [SemVer](https://semver.org/) for versioning. For all available versions, check the
[tags in this repository](https://github.com/realfabecker/kevin/tags).

## Licença

This project is licensed under the MIT License. See the [License](LICENSE.md) for more information.
