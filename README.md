# fcosctl

_A CLI for rapid Fedora CoreOS development_

When iterating on a configuration for a Fedora CoreOS machine the loop of
transpiling, launching and debugging can be a little cumbersome. `fcosctl`
tries to make that a little more convenient. To launch an local ephemeral
instance simply pull an image,

```shell
fcosctl image pull --stream=testing
```

and point to a config.

```shell
fcosctl run config.yaml
```

## Usage

`fcosctl` depends on 

- `qemu-kvm`
- `fcct`
- `coreos-installer`

Install these before proceeding. Images can be managed using the `fcosctl
image` subcommand

```
$ fcosctl image --help
Manage FCOS images

Usage:
  fcosctl image [command]

Available Commands:
  delete      Delete FCOS images
  list        List available images
  pull        Download FCOS images

Flags:
  -h, --help   help for image

Use "fcosctl image [command] --help" for more information about a command.
```

Once you have some images available locally you can run them with the `fcosctl
run` subcommand

```
$ fcosctl run --help
Runs a config using qemu-kvm in an ephemeral virtual machine

Usage:
  fcosctl run [fcos config] [flags]

Flags:
  -h, --help             help for run
      --version latest   Image version to use as base. Use latest to run the most recent (default "latest")
```
