# Smart Directory Navigation with `go`

`go` is a smarter way to move between directories. Just type the name of the
folder you want, and `go` will almost always bring your there, even if it has to
jump across the file tree.

The four different methods that `go` uses to determine a destination directory,
are, in order of preference:

1. `cd` - If you just want to `cd` somewhere, `go` will follow the path
   appropriately.
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

Since `go` changes the current directory, it must be run as a bash _function_,
not as a bash program<a href="foot-1"><sup>1</sup></a>. This means that rather
than adding the executable to your path, you must define the `go` function. This
needs to occur in your `.bashrc`,`.bash_profile`, `.profile`, or whatever other
login file you choose to use. For most people, the following command will make
sure the function is loaded when you log in to your terminal:

```bash
echo ". $(pwd)/go" >> ~/.bashrc
```

## Usage

```bash
go [dest]
go [-h,--help]
go [-a,--add] linkname destpath
go [-r,--remove] linkname
go [-l,--list]
```

The `dest` given to `go` can be one of the following:

1. Nothing

    If no destination is given, `go` changes the current directory to the home
    directory (same as `cd`).
2. Pathname

    Both relative and absolute pathnames are accepted by `go`, and will be
    navigated to the same way that `cd` navigates to them.
3. Shortcut

    If a name does not resolve to an absolute or relative path, `go` will check
    if the name matches a shortcut. If a matching shortcut is found, the current
    directory is set to the directory the shortcut points to<a
    href="foot-2"><sup>2</sup></a>.

    If a shortcut is found, but the link is broken, the shortcut will be deleted
    to avoid future false positives.
4. Folder Base Name

    When a name that cannot be evaluated as a relative path, absolute path, or
    shortcut is given, `go` will attempt to find a directory with the given
    name. This base name should either correspond to a recent directory


### Recent Directory Example

In case there is confusion about the recently visited directory behavior, the
following is a use case showing a directory that is visited often, but that is
deep in the file tree:

```bash
$ go notes
go: notes: Directory not found.
$ go ~/Documents/CMU/F16/15410/notes
$ go # Goes to ~
$ go notes # Goes to ~/Documents/CMU/F16/15410/notes
```

## Shortcut Management

Shortcuts are stored in the form of [symbolic links](https://en.wikipedia.org/wiki/Symbolic_link#POSIX_and_Unix-like_operating_systems).
For now, `go` only looks for shortcuts in one directory. By default this
directory is `~/.go_shortcuts`, but it can be set by exporting a new
global `GO_SHORTCUT_DIR` in your `.bashrc` (or similar login file).

To save you from having to manipulate these directories yourself, `go` has
`--add|-a` and `--remove|-r` subcommands to add and remove shortcuts,
respectively<a href="foot-3"><sup>3</sup></a>:

```
# Adding a shortcut
$ go --add os ~/Documents/CMU/F16/15410
$ go os # Goes to /Users/Scott/Documents/CMU/F16/15410

# Removing a shortcut
$ go --remove os
```

There is also a `--list|-l` subcommand that allows you to view a list of your
saved shortcuts:

```
$ go --list
os -> /Users/Scott/Documents/CMU/F16/15410
```

## Replacing `cd`

It may be tempting to entirely replace `cd` with `go`. For yourself, this is
totally reasonable, but before aliasing away your old `cd` command, remember
that other programs may rely on its exact behavior.

I recommend typing `go` whenever you are navigating anywhere in the terminal,
and leaving the `cd` command as-is to ensure you don't screw up other programs.

If you use the [Go Language](https://golang.org), renaming or aliasing `go` is
a good idea. But in this case, it's still probably better to rename or alias it
to `goto` (or similar) rather than `cd`.

---

<a href="foot-1"><sup>1</sup></a> Typical command line programs execute in their
own process (forked from the main shell process). Executing a `cd` command
inside of such a forked process will not have an effect on the parent shell
process. Functions, on the other hand, are run as part of the current process,
hence `go` is defined as a function rather than a program.

<a href="foot-2"><sup>2</sup></a> This means that if a symbolic link in the
shortcut folder named `mylink` points to `/Users/Scott/Documents/Courses/15410`,
`go mylink` changes the current directory to
`/Users/Scott/Documents/Courses/15410` not `~/.go_shortcuts/mylink`. Going to
the directory that the symbolic link points to, instead of setting the symbolic
link to the current directory, keeps file paths cleaner.

<a href="foot-3"><sup>3</sup></a> Note that `--add` and `-a`, `--remove` and
`-r`, and `--list` and `-l`, are long- and short-hand forms of the same exact
command. There is no functional difference between using the long or short form,
it is a matter of personal preference.
