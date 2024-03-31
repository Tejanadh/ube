# cloc-tool

A command line tool to count lines of code (CLOC), written in Go.

## Example
<img width="505" alt="Screen Shot 2024-03-31 at 1 25 03 AM" src="https://github.com/ramirezfernando/cloc-tool/assets/91701930/a667ac6b-9a37-4a7f-b901-84c25edeed44">

## Features
- cloc-tool has a huge range of languages, supporting over **220** language extensions.
- cloc-tool is **accurate**, and **consistent** as it counts the number of newline characters `/n` present in a specified path. This ensures consistency across different platforms and text editors. Different text editors may interpret line endings differently (e.g., `\n` in Unix-like systems, `\r\n` in Windows), which could lead to discrepancies in line counts if you try to match the exact number of lines displayed in a specific editor.

## What's next?
- Add unit tests to double check the accuracy
- Making your command available to Homebrew users
