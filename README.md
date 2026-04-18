# go-booba-example

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
go tool booba-wasm-build -o web/app.wasm ./cmd/example/

# 2. Populate web/ with wasm_exec.js, booba/, and ghostty-web/.
go tool booba-assets web/

# 3. Serve the web/ directory with any static file server.
npx serve web
```

Then open the URL the server prints (typically http://localhost:3000).

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

## How It Works

The example is a single `cmd/example/main.go` that calls `booba.Run(initialModel())`.
`booba.Run` is platform-polymorphic:

- **Native builds** call `tea.NewProgram(model).Run()` — a normal BubbleTea app.
- **`GOOS=js GOARCH=wasm` builds** delegate to `wasm.Run`, which installs the
  JavaScript bridge that booba's browser terminal uses.

No build tags are required in your code; the split is inside `go-booba` itself.

In the browser, [ghostty-web](https://github.com/coder/ghostty-web)
renders a real terminal emulator in a canvas, and booba's `BoobaWasmAdapter`
shuttles data between that terminal and your compiled program.

Two helper commands from `go-booba` handle the rest:

- `booba-wasm-build` — compiles to WASM, working around BubbleTea v2's missing
  `js/wasm` build tags.
- `booba-assets` — populates `web/` with `wasm_exec.js`, the booba terminal
  wrapper, and ghostty-web runtime files.

## Project Structure

```
├── .github/workflows/pages.yml   # Build WASM + deploy to GitHub Pages
├── cmd/example/
│   └── main.go                   # BubbleTea model + booba.Run entrypoint
├── web/
│   └── index.html                # Static page (customize freely)
└── go.mod
```

The following are generated and gitignored: `web/app.wasm`, `web/wasm_exec.js`,
`web/booba/`, `web/ghostty-web/`.

## Customizing

Replace the model, update, and view functions in `cmd/example/main.go` with
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
