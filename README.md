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

Wanted features are as follows:

1. Hooks (Any language as long as they speak environment variables.)
2. Declare custom "pkg-config" package configurations.
3. Cleaner and better written source code.
4. More compilers. Current support compilers are clang and rustc.
