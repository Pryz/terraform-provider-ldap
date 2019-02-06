# Hacking this plugin

If you want to hack this plugin, you need to have a working Golang build enviornment (see the [Golang website](https://golang.org/doc/install)) to get started.

## Enabling traces

Terraform has its own mechanism for enabling traces, as detailed [here](https://www.terraform.io/docs/internals/debugging.html).

This plugin logs to `stdout`, and its outputs are collected by terrafiorm and interleaved with those of the internal engine.

To enable logging, define the `TF_LOG` and optionally the `TF_LOG_PATH` environment variables before running the command, e.g.:

```bash
$> TF_LOG=TRACE TF_LOG_PATH=./terraform.log terraform plan
```

