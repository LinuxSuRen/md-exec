[![Codacy Badge](https://app.codacy.com/project/badge/Grade/5022a74d146f487581821fd1c3435437)](https://www.codacy.com/gh/LinuxSuRen/md-exec/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=LinuxSuRen/md-exec&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/5022a74d146f487581821fd1c3435437)](https://www.codacy.com/gh/LinuxSuRen/md-exec/dashboard?utm_source=github.com&utm_medium=referral&utm_content=LinuxSuRen/md-exec&utm_campaign=Badge_Coverage)

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

### Run Python Script
```python3
#!title: Python Hello World
print('hello python world');
```

## Run Golang
```golang
#!title: Golang Hello World
fmt.Println("hello golang")

items := []int{1,2,3,4}
for _, item := range items {
    fmt.Println(item)
}
```

## Limitation
Please make sure the Markdown files meet Linux end-of-line.
You could turn it via: `dos2unix your.md`
