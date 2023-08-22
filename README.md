# README

## About

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


## Tasks execution
Page: https://pydoit.org/contents.html
- install task manager doit: `pip3 install doit` or `pip install doit`
- look at the commands: `doit` to show help or `doit list` to show available tasks 

## Generate mocks
Install mockery: `go install github.com/vektra/mockery/v2@v2.29.0`

## Get stared
- `git submodule update --init`
- `cd frontend`
- `npm i && npm run build`
- `cd ../`
- `wails dev`
- [optional] open in a browser: http://localhost:34115

## Troubleshooting
Node runtime may fail due to lack of memory.
For such case set a var NODE_OPTIONS="--max-old-space-size=4096"
