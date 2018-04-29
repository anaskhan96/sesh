# sesh

![](https://raw.githubusercontent.com/anaskhan96/sesh/master/preview.png?token=AWKfRu3IiZ6pI7l5w1BkShc5PbCrnqsNks5a7ub7wA%3D%3D)

`sesh` is a simple (read basic), elegant shell written in Go. Built as a school project for the course _Unix Systems Programming_, it supports the following:
+ Aliasing
+ Piping and I/O redirection
+ Aliasing
+ Arrow keys up and down for history
+ Tab autocompletion

Apart from this, it has two custom builtins:
+ `walk`: walks through the directory specified as an argument recursively. Takes the current directory as input if no argument is specified.
+ `show`: lists the commands in the PATH having the given argument as its prefix. Lists all the commands in the PATH if no argument is specified.

### Installation

```bash
go install github.com/anaskhan96/sesh
```

It can be run by invoking `sesh` from anywhere in the terminal.
