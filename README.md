# Bullshit Generator

The **Bullshit Generator** is a Python program inspired by the [Plan 9](https://9p.io/plan9/) [`bullshit`](http://man.9front.org/1/bullshit) command. It generates random nonsense phrases by combining words and phrases from a customizable data file, making it perfect for creating humorous or jargon-filled text.

## Features

- Randomized sentence generation using categories like:
  - Starting words
  - Suffixes
  - Protocols
  - Endings
- Support for custom word files via the `-f` option or environment variables.
- Default file locations (`$BULLSHIT_FILE` or `~/.bullshit`).
- Generate multiple phrases with a single command.

## Installation

Clone the repository and make the script executable:
```bash
git clone https://github.com/yourusername/bullshit-generator.git
cd bullshit-generator
chmod +x bullshit.py
```

## Usage

```bash
./bullshit.py [n] [-f file]
```

### Options

- `n`: (Optional) The number of sentences to generate. Default is `1`.
- `-f file`: (Optional) Specify a custom data file. If omitted, the program uses:
  1. The path set in the environment variable `$BULLSHIT_FILE`.
  2. The file `~/.bullshit`.

### Examples

#### Generate a Single Nonsense Phrase
```bash
./bullshit.py
```

#### Generate 5 Phrases
```bash
./bullshit.py 5
```

#### Use a Custom Word File
```bash
./bullshit.py -f custom_bullshit.txt
```

#### Set a Default File via Environment Variable
```bash
export BULLSHIT_FILE=/path/to/bullshit.txt
./bullshit.py
```

## Data File Format

The data file must contain one word or phrase per line. Words can be optionally categorized using a suffix:
- `*`: Protocols (e.g., `HTTP`, `SSL`)
- `$`: Endings (e.g., `manager`, `interface`)
- `%`: Suffixes (e.g., `-driven`, `-enabled`)
- `^`: Starting words (e.g., `cloud`, `agile`)
- `|`: Words that cannot end a phrase (e.g., `realtime`, `cache`)
- No suffix: General words.

### Example Data File
```plaintext
cloud ^
engine
metadata |
HTTP *
manager $
-driven %
```

## Requirements

- Python 3.6 or higher

## License

This project is licensed under the Zlib License. See the [LICENSE](LICENSE) file for details.

## Credits

This program is inspired by the `bullshit` command from Plan 9's 9front distribution (August 2011).
