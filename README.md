# CC-Shell - Custom Command-Line Shell

This project is a custom implementation of a command-line shell built with Go. It was developed as part of a coding challenge [here](https://codingchallenges.fyi/challenges/challenge-shell). The shell provides basic command execution capabilities, along with history management and directory navigation.

## Features

- Execute common shell commands (`ls`, `pwd`, `cd`, etc.).
- Maintain a persistent command history across sessions.
- Use `cd -` to switch to the previous directory.
- Handle `Ctrl+C` gracefully without terminating the shell.

## Getting Started

These instructions will help you set up and run the project on your local machine for development and testing.

### Prerequisites

- You need to have Go installed on your machine (Go 1.18 or later is recommended).
- You can download and install Go from [https://golang.org/dl/](https://golang.org/dl/).

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/nullsploit01/cc-shell.git
cd cc-shell
```

### Building

Compile the project using:

```bash
go build -o cc-shell
```

### Usage

To start the shell, run the compiled binary:

```bash
./cc-shell
```

### Example Commands

```bash
> pwd
/home/user

> ls
file1.txt  file2.txt  Documents/

> cd Documents
> pwd
/home/user/Documents

> cd ~
> pwd
/home/user

> cd -
/home/user/Documents

> history
1 pwd
2 ls
3 cd Documents
4 pwd
5 cd ~
6 pwd
7 cd -

> Ctrl+C
Ignored Ctrl+C. Use 'exit' to quit.
> exit
Exiting... Bye!
```

### Running the Tests

```bash
go test ./...
```
