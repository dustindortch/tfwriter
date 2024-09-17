# tfwriter

tfwriter is a CLI to that writes templated Terraform blocks.

## Usage

```bash
tfwriter [subcommand] [flags]
```

## Subcommands

* list
* resource

### List

```bash
tfwriter list [<provider>]
```

Lists the available resources, optionally filtered by provider.

### Resource

```bash
tfwriter resource <resource_type.label1> [<resource_type.labelN>...]
```

Outputs Terraform blocks for the provided resource types.
