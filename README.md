# This is a simple terminal TODO-App

The application saves its data in a local sqlite db file

## Install

1. Install go according to [the instructions on the official website](https://go.dev/doc/install)
2. Clone this repository into $GOPATH/src
3. Install according to your system specifications with the `go install` command
4. Add the $GOPATH/bin to your $PATH variable to make the binary a global command
5. Start the application from anywhere with `todo` (or the name under which the binary was installed)

## Shortcuts

`n` -- create new entry

`N` -- create new list

---

`d` -- delete entry

`D` -- delete list

---

`i` -- edit entry

`I` -- edit list name

---

`j` -- go one entry down

`J` -- switch the entry with the one below

`k` -- go one entry up

`K` -- switch the entry with the one above

---

`h` -- go one list to the left

`H` -- switch the list with the one to the left

`l` -- go one list to the right

`L` -- switch the list with the one to the right

---

`enter` -- toggle entry

`0-9` -- switch to the list if it exists

`x` -- exit
