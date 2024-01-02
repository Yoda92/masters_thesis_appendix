---
description: This section describes the configuration parameters and their types for WASP.
keywords:
- IOTA Node 
- Hornet Node
- WASP Node
- Smart Contracts
- L2
- Configuration
- JSON
- Customize
- Config
- reference
---


# Core Configuration

WASP uses a JSON standard format as a config file. If you are unsure about JSON syntax, you can find more information in the [official JSON specs](https://www.json.org).

You can change the path of the config file by using the `-c` or `--config` argument while executing `wasp` executable.

For example:
```shell
wasp -c config_defaults.json
```

You can always get the most up-to-date description of the config parameters by running:

```shell
wasp -h --full
```

