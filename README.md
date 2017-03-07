# Smart Directory Navigation with `go`

The `go-navigate` command is a smarter way to move between directories. Just
type the base-name of the folder you want, and `go-navigate` will almost always
bring you there, even if it has to jump across the file tree. When combined
with a short alias, it makes changing directories effortless -- I rarely find
the need to type more than a basename.

The four different methods that `go-navigate` uses to determine a destination
directory, are, in order of preference:

1. `cd` - If it is possible to `cd` into the supplied relative path,
   `go-navigator` will perform the `cd` operation. This means you can start
   using `go-navigator` as a simple drop in for `cd` and it will do the right
   thing.
2. Shortcuts - Custom shortcuts (symbolic links) stored in the ~/.go_links
   allow you to navigate to tricky places with minimal keystrokes.
3. Most Recent - Navigate to the most recently visited directory with the same
   base name as the one given. This acts as a cache, making it extremely
   efficient to access your most frequently used directories.
4. Depth-limited Search [Experimental] - Search through top-level directories
   starting from the home folder, and continuing at folder depths of 1-4. This
   is a very slow search method, but will find almost any folder. Once the
   folder is found, it will be added to the history, making future accesses
   faster. _Feedback on the usefulness of this feature, compared to its slow
   speed, is appreciated._

## Installation

Since `go-navigate` changes the current directory, it must be run as a bash
_function_, not as a bash program<a href="foot-1"><sup>1</sup></a>. This means
that rather than adding the executable to your path, you must define the
`go-navigate` function. This needs to occur in your `.bashrc`,`.bash_profile`,
`.profile`, or whatever other login file you choose to use. For most people,
the following command will make sure the function is loaded when you log in to
your terminal:

```bash
echo ". $(pwd)/go-navigate" >> ~/.bashrc
```

## Usage

```bash
go-navigate [dest]
go-navigate [-h,--help]
go-navigate [-a,--add] linkname destpath
go-navigate [-r,--remove] linkname
go-navigate [-l,--list]
```

The `dest` given to `go-navigate` can be one of the following:

1. Nothing

    If no destination is given, `go-navigate` changes the current directory to
    the home directory (same as `cd`).
2. Pathname

    Both relative and absolute pathnames are accepted by `go-navigate`, and will be
    navigated to the same way that `cd` navigates to them.
3. Shortcut

    If a name does not resolve to an absolute or relative path, `go-navigate`
    will check if the name matches a shortcut. If a matching shortcut is found,
    the current directory is set to the directory the shortcut points to
    <a href="foot-2"><sup>2</sup></a>.

    If a shortcut is found, but the link is broken, the shortcut will be deleted
    to avoid future false positives.
4. Folder Base Name

    When a name that cannot be evaluated as a relative path, absolute path, or
    shortcut is given, `go-navigate` will attempt to find a directory with the
    given name. This base name should either correspond to a recent directory

### Recent Directory Example

In case there is confusion about the recently visited directory behavior, the
following is a use case showing a directory that is visited often, but that is
deep in the file tree:

```bash
$ go-navigate notes
go-navigate: notes: Directory not found.
$ go-navigate ~/Documents/CMU/F16/15410/notes
$ go-navigate # Goes to ~
$ go notes # Goes to ~/Documents/CMU/F16/15410/notes
```

## Shortcut Management

Shortcuts are stored in the form of [symbolic links](https://en.wikipedia.org/wiki/Symbolic_link#POSIX_and_Unix-like_operating_systems).
For now, `go-navigate` only looks for shortcuts in one directory. By default
this directory is `~/.go_shortcuts`, but it can be set by exporting a new
global `GO_SHORTCUT_DIR` in your `.bashrc` (or similar login file).

To save you from having to manipulate these directories yourself, `go-navigate`
has `--add|-a` and `--remove|-r` subcommands to add and remove shortcuts,
respectively<a href="foot-3"><sup>3</sup></a>:

```
# Adding a shortcut
$ go-navigate --add os ~/Documents/CMU/F16/15410
$ go-navigate os # Goes to /Users/Scott/Documents/CMU/F16/15410

# Removing a shortcut
$ go-navigate --remove os
```

There is also a `--list|-l` subcommand that allows you to view a list of your
saved shortcuts:

```
$ go-navigate --list
os -> /Users/Scott/Documents/CMU/F16/15410
```

## Replacing `cd`

It may be tempting to entirely replace `cd` with `go-navigate`. For a human,
this is totally reasonable, but before aliasing away your old `cd` command,
remember that other programs may rely on its exact behavior.

I recommend using an alias for `go-navigate` whenever you change directories in
your terminal, and leaving the `cd` command as-is to ensure you don't screw up
other programs.

If you do not use the [Go Language](https://golang.org), creating an alias `go`
for `go-navigator` convenient. `gn`, `gd` or the more verbose `goto` are also
good aliases. You can define them in your `.bashrc` like so:


```
alias gn='go-navigate'
```

---

<a href="foot-1"><sup>1</sup></a> Typical command line programs execute in their
own process (forked from the main shell process). Executing a `cd` command
inside of such a forked process will not have an effect on the parent shell
process. Functions, on the other hand, are run as part of the current process,
hence `go-navigate` is defined as a function rather than a program.

<a href="foot-2"><sup>2</sup></a> This means that if a symbolic link in the
shortcut folder named `mylink` points to `/Users/Scott/Documents/Courses/15410`,
`go-navigate mylink` changes the current directory to
`/Users/Scott/Documents/Courses/15410` not `~/.go_shortcuts/mylink`. Going to
the directory that the symbolic link points to, instead of setting the symbolic
link to the current directory, keeps file paths cleaner.

<a href="foot-3"><sup>3</sup></a> Note that `--add` and `-a`, `--remove` and
`-r`, and `--list` and `-l`, are long- and short-hand forms of the same exact
command. There is no functional difference between using the long or short form,
it is a matter of personal preference.
