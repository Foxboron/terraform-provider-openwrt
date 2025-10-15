terraform-provider-openwrt
==========================

openwrt provider for Terraform.

## Prerequisites

- [Terraform](http://terraform.io)
- [OpenWrt](https://openwrt.org)

## Installation

This provider is published in the [Terraform Registry](https://registry.terraform.io/providers/foxboron/openwrt).

### Quick Example

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    openwrt = {
      source = "foxboron/openwrt"
    }
  }
}
```

## Documentation

Full documentation can be found in the [`docs`](docs) directory.

## Development

#### Prerequisites

- [Golang](https://golang.org/doc/install)
- [An OpenWRT Router](https://openwrt.org)
- [Make](https://www.gnu.org/software/make/)
- [Golangci-lint](https://golangci-lint.run/)
- [GoReleaser](https://goreleaser.com/)

#### Setup

1. Checkout the repository `git clone ...`
2. Compile from sources to a development binary:

```shell
cd terraform-provider-openwrt
make build
```

3. Configure Terraform (`~/.terraformrc`) to use the development binary provider:

```shell
$ cat ~/.terraformrc
provider_installation {
  dev_overrides {
    "foxboron/openwrt" = "/<PATH_WHERE_YOU_EXECUTED_GIT_CLONE>/terraform-provider-openwrt"
  }
  direct {}
}
```
