#!/usr/bin/env python3

import random
import sys
import os
import argparse

def locate_file():
    if "BULLSHIT_FILE" in os.environ:
        return os.environ["BULLSHIT_FILE"]
    if "HOME" in os.environ:
        return os.environ["HOME"] + "/.bullshit"
    raise ValueError("$HOME or $BULLSHIT_FILE not set")

# Data Structures
protocols = []
ends = []
suffixes = []
starts = []
noends = []
words = []


def load_data(file_path):
    """Load data from the input file into respective categories."""
    with open(file_path) as file:
        for line in file:
            params = line.strip().split(' ')
            if len(params) == 1:
                params.append("")  # Add empty type if not provided
            word, typ = params

            match typ:
                case '*': protocols.append(word)
                case '$': ends.append(word)
                case '%': suffixes.append(word)
                case '^': starts.append(word)
                case '|':
                    words.append(word)
                    noends.append(word)
                case _: words.append(word)


def random_suffix():
    """Randomly choose a suffix or return an empty string."""
    if random.random() < 0.2:
        return random.choice(suffixes)
    return ""


def generate_bullshit():
    """Generate and print a single bullshit sentence."""
    last = None
    output_count = 0
    total_words = random.randint(3, 10)
    result = []

    # Add starting words
    num_starts = random.randint(0, 3)
    for _ in range(num_starts):
        result.append(random.choice(starts))
        output_count += 1

    # Add main words with optional suffixes
    remaining = min(total_words - output_count, 3)
    num_words = random.randint(0, remaining) if output_count < total_words else 0
    for _ in range(num_words):
        word = random.choice(words)
        result.append(word + random_suffix())
        last = word
        output_count += 1

    # Add protocol section
    if random.random() > 0.5:
        num_protocols = random.randint(0, 3)
        for i in range(num_protocols):
            result.append(random.choice(protocols))
            if i != num_protocols - 1:
                result.append("over")
            output_count += 1
        last = None

    # Add more words
    remaining = min(total_words - output_count, 3)
    num_more_words = random.randint(0, remaining) if output_count < total_words else 0
    if output_count + num_more_words <= 1 or last is None:
        num_more_words += 2
    for _ in range(num_more_words):
        word = random.choice(words)
        result.append(word + random_suffix())
        last = word

    # Optionally add an ending
    if random.random() < 0.1 or (last in noends) or any(suffix in result[-1] for suffix in suffixes):
        result.append(random.choice(ends))

    print(" ".join(result))


def main():
    """Main function to handle arguments and generate bullshit sentences."""
    parser = argparse.ArgumentParser(description="Generate random bullshit sentences.")
    parser.add_argument(
        "n", nargs="?", type=int, default=1, help="Number of sentences to generate (default: 1)"
    )
    parser.add_argument(
        "-f", "--file", type=str, default=locate_file(), help="Path to the data file (default: $BULLSHIT_FILE or ~/.bullshit)"
    )
    args = parser.parse_args()

    # Load data from the specified file
    try:
        load_data(args.file)
    except FileNotFoundError:
        print(f"Error: File '{args.file}' not found.")
        sys.exit(1)

    # Generate the specified number of sentences
    for _ in range(args.n):
        generate_bullshit()


if __name__ == "__main__":
    main()
