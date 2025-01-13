# Bullshit Generator

The **Bullshit Generator** is a program inspired by Plan 9's [`bullshit`](http://man.9front.org/1/bullshit) script. It generates random nonsense phases by combining words and phrases from a customizable data file, making it perfect for creating humorous or jargon-filled text.

## Features

- Randomized sentence generation using categories like:
  - Starting words
  - Suffixes
  - Protocols
  - Endings
- Support for custom word files via the `-f` option or environment variables.
- Generate multiple phrases with a single command.

## Requirements

- [Go Compiler](https://go.dev/)


## Installation

Clone the repository and build `bullshit` using the Go Compiler:
```bash
go build bullshit.go
```

## Usage

```bash
./bullshit [-s] [-f file] [times]
```

### Options

- `times`: (Optional) The number of sentences to generate. Default is `1`.
- `-f file`: (Optional) Specify a custom data file. If omitted, the program uses:
  1. The path specified in the environment variable `$BULLSHIT_FILE`.
  2. The file `$XDG_CONFIG_HOME/bullshit.txt` or `~/.config/bullshit.txt`.
  3. The file `/usr/share/bullshit.txt`.
- `-s`: (Optional) Sort input records and print sorted file to stdout

### Examples

#### Generate a Single Nonsense Phrase
```bash
./bullshit
```

#### Generate 5 Phrases
```bash
./bullshit 5
```

#### Use a Custom Word File
```bash
./bullshit -f custom_bullshit.txt
```

#### Set a Default File via Environment Variable
```bash
export BULLSHIT_FILE=/path/to/bullshit.txt
./bullshit
```

## License

This project is licensed under the Zlib License. See the [LICENSE](LICENSE) file for details.

## Credits

This program is inspired by the `bullshit` command from [9front's](https://git.9front.org/plan9front/plan9front/HEAD/rc/bin/bullshit/f.html).
