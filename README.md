## Usage
`md-exec` could exec the commands in the Markdown files.
For instance, it will execute those commands which in Markdown block via `mde README.md`

```shell
#!title: Print ip info
ifconfig
```

### Variable input support
In some use cases, we need to change the variables or command line flags. Try the following demo:

```shell
#!title: Variable Input Hello World
name=linuxsuren
echo hello $name
```

## Limitation
Please make sure the Markdown files meet Linux end-of-line.
You could turn it via: `dos2unix your.md`
