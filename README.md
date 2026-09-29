# go-booba-example

[Link the GitHub Pages Example](https://nimblemarkets.github.io/go-booba-example/)

A complete example of using [go-booba](https://github.com/NimbleMarkets/go-booba) to host a [BubbleTea](https://github.com/charmbracelet/bubbletea) program on GitHub Pages.

The BubbleTea program compiles to WebAssembly and runs entirely in the browser —
no backend server needed. No npm dependencies. No shell scripts.

## Quick Start

1. Fork this repo
2. Enable GitHub Pages: **Settings → Pages → Source → GitHub Actions**
3. Push to `main` — the workflow builds and deploys automatically

## Local Development

Run the program natively in your terminal:

```sh
go run ./cmd/example/
```

Build and serve the WASM version:

```sh
# 1. Compile the program to WASM.
GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/example/

# 2. Populate web/ with wasm_exec.js, booba/, and ghostty-web/.
go tool booba-assets web/

# 3. Serve the web/ directory with any static file server.
npx serve web
```

Then open the URL the server prints (typically http://localhost:3000).

## Serving it over HTTP

The same BubbleTea program can also be served over HTTP+WebSocket via
go-booba's `serve` package, so users can connect from a browser or a native
terminal without compiling the app to WASM themselves.

Build and run the server entrypoint:

```sh
task build-native-server
./bin/example-server --listen 127.0.0.1:8080
```

Open `http://127.0.0.1:8080/` in a browser to use the built-in terminal UI
(BoobaTerminal), or connect from a native terminal using the companion
`booba-sip-client`:

```sh
booba-sip-client ws://localhost:8080/ws
```

`booba-sip-client` ships with go-booba (`cmd/booba-sip-client/`) for headless
or SSH-tunneled use. Press Ctrl-] then type `quit` to disconnect.

### Taskfile

A [Taskfile](https://taskfile.dev) wraps the commands above. With
[`task`](https://taskfile.dev/installation/) installed:

```sh
task list         # show all tasks
task build-native # build bin/example
task build-web    # compile WASM + populate web/
task serve        # build web assets and serve on :3000
task test         # run go test ./...
task go-lint      # run golangci-lint
task clean        # remove build artifacts
```

`task` (no args) runs `test` + `build`. `task dev-deps` installs the Go dev
tools (golangci-lint, errcheck, godoc).

To run in the terminal or serve via webpage:
  * `task run`
  * `task serve`

## How It Works

The example is a `cmd/example/main.go` that calls `booba.Run(model.InitialModel())`.
`booba.Run` is platform-polymorphic:

- **Native builds** call `tea.NewProgram(model).Run()` — a normal BubbleTea app.
- **`GOOS=js GOARCH=wasm` builds** delegate to `wasm.Run`, which installs the
  JavaScript bridge that booba's browser terminal uses.

No build tags are required in your code; the split is inside `go-booba` itself.

In the browser, [ghostty-web](https://github.com/coder/ghostty-web)
renders a real terminal emulator in a canvas, and booba's `BoobaWasmAdapter`
shuttles data between that terminal and your compiled program.

WASM compilation uses the standard Go toolchain (`GOOS=js GOARCH=wasm go build`) —
no custom build tool is needed. The `booba-assets` helper (declared as a
`tool` in `go.mod`) populates `web/` with `wasm_exec.js`, the booba terminal
wrapper, and ghostty-web runtime files.

### BubbleTea fork

Building for the browser depends on a
[patched BubbleTea fork](https://github.com/neomantra/bubbletea/tree/nm-wasm)
that adds `GOOS=js` support (not yet merged upstream). Go does not propagate
`replace` directives from dependencies, so this repo's `go.mod` carries its own:

```
replace charm.land/bubbletea/v2 => github.com/neomantra/bubbletea/v2 v2.0.0-20260928192001-1b36865b418a
```

Keep this line in sync with
[go-booba's `go.mod`](https://github.com/NimbleMarkets/go-booba/blob/main/go.mod)
when upgrading go-booba. Compile errors inside `charm.land/bubbletea/v2` such as
`p.listenForResize undefined` mean the directive is missing or stale.

Go 1.26 or newer is required.

## Project Structure

```
├── .github/workflows/pages.yml   # Build WASM + deploy to GitHub Pages
├── cmd/example/
│   └── main.go                   # booba.Run entrypoint (native + WASM)
├── cmd/example-server/
│   └── main.go                   # HTTP+WebSocket server entrypoint
├── internal/model/
│   └── model.go                  # Shared BubbleTea model
├── web/
│   └── index.html                # Static page (customize freely)
└── go.mod
```

The following are generated and gitignored: `web/app.wasm`, `web/wasm_exec.js`,
`web/booba/`, `web/ghostty-web/`.

## Customizing

Replace the model, update, and view functions in `internal/model/model.go` with
your own BubbleTea program. The model should handle `tea.WindowSizeMsg` to
adapt to the browser window size.

Customize `web/index.html` for your own title, theme, or layout — the
`booba-assets` tool won't overwrite it unless you pass `--force`.

Keep the `booba.Run(initialModel())` call in `main` as-is — it's what
makes the same source build for both native terminals and the browser.

## See Also

- [go-booba](https://github.com/NimbleMarkets/go-booba) — the library

## License

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).

Copyright (c) 2026 [Neomantra Corp](https://www.neomantra.com).   

----
Made with :heart: and :fire: by the team behind [Nimble.Markets](https://nimble.markets).
