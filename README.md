# Chain

Chain is a build system aiming to replace the frustrations with CMake
and meson but still remaining feature rich. The main language used is
YAML as it is easy to use and readable. 

## Installation

```shell
go build main.go
sudo cp main /usr/local/bin/chain 
```

You don't need to use /usr/local/bin/chain use whatever you want.

## Features

Implemented features are as follows:

1. Package lookup with "pkg-config"
2. Compiling to shared libraries (.so files)
3. Link with said libraries
4. Scoped environments for declarations hence the need for exporting.
5. Multiple compilers with clang and rustc.
6. Hooks (Any language as long as they speak environment variables.)
7. Cleaner and better written source code.

Wanted features are as follows:

1. Declare custom "pkg-config" package configurations.
2. More compilers. Current support compilers are clang and rustc.
3. Support more formats like JSON or TOML for procedure and project
   files.
4. Running **chain .** and chain should be able to find project.yml and
   run it.
