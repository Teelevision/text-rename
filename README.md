# text-rename
Rename files using your text editor.

## How to use ##

1. Open the program with one or more files, e.g. via drag and drop.
2. A file will open in your default text editor. Do your renaming here.
3. Open the program with the text file you just edited.

## How it works ##

The program creates two files: one containing the file names and one containing their original path. The matching is done by the line number

## Compatibility ##

Tested under Windows 7. Feel free do adapt to other OS. Asking me to do so might work as well.

## Dependencies ##

* https://github.com/skratchdot/open-golang
  * This makes it possible to open the file in the system's default text editor.

## Motivation ##

Renaming files under Windows gives me a hard time. I need search and replace, regex and column mode editing (Notepad++). So I wrote a program that creates a text file where I can do all that.

## License ##

import "LICENSE"
