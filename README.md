# static-server
A lightwheight web server for hosting static sties, that require index fallback (aka Single Page Web Apps)
Written in pure Go, using only inbuildt packages. The single executable is just 6 Mb.

## Usage
**Arguments:**
- `--dir`: The directory that contains the files. Default: `./public`.
- `--port`: The port of the server. Default: `8000`.
- `--logging`: Log requests. Default: `true`.
