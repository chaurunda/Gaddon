# Gaddon

A Go CLI program that update your wow addon and install new one

Usage: gaddon [OPTIONS]

Options:

* u, update <FOLDER>       check if the given FOLDER is a git one and if there is an update
* -h, --help               Print help
* -V, --version            Print version


exemple 

check for update : 
```
gaddon u /path/to/wow/addon
```

install a new addon (incoming)

```
gaddon --folder /path/to/wow/addon --install url/to/repo.git
```


This repo is a portage in Go of a program I have written in Rust : https://github.com/chaurunda/Raddon
