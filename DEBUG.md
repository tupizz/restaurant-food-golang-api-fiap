# Update your Makefile to have the following:
```makefile
# the name of the binary when built
BINARY_NAME=fastfood-golang

# remove any binaries that are built
clean:
	rm -f ./bin/$(BINARY_NAME)*

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug main.go
```

# Update .air.toml to have the following:
```toml
# Config file for [Air](https://github.com/cosmtrek/air) in TOML format
root = "."
tmp_dir = "tmp"

[build]
cmd = "make build-debug"
bin = "./bin/my-app-debug"
full_bin = "dlv exec ./bin/my-app-debug --listen=127.0.0.1:2345 --headless=true --api-version=2 --accept-multiclient --continue --log -- "
include_ext = ["go"]
exclude_dir = [".vscode", ".github", "bin", "tmp"]
exclude_regex = ["_test.go"]
exclude_unchanged = true.
args_bin = ["server"]

[misc]
clean_on_exit = true

[screen]
clear_on_rebuild = true
keep_scroll = true
```

# Add the following to your .vscode/launch.json
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to Air",
      "type": "go",
      "mode": "remote",
      "request": "attach",
      "host": "127.0.0.1",
      "port": 2345
    }
  ]
}
```