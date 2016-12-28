#!/bin/bash
#
# go - A faster way to change directories.
#
# The go command must be implemented as a function, rather than a normal executable
#
# Implementation = 30% parsing + 60% error handling + 10% actual useful code
#
# (c) Scott Krulcik 2017 MIT License
#

function go() {

    USAGE="Usage:
    go [dest]
    go [-h,--help]
    go [-a,--add] linkname destpath
    go [-r,--remove] linkname";

    HELP="
    [dest] can be empty, a relative or local path, a pre-defined shortcut in
    $GO_SHORTCUT_DIR, or the name of a recently visited directory.

--add or -a:
    Adds a new shortcut for the go command to follow.

    linkname - An alphanumeric name for the shortcut. This should be a name
               that is short and easy to type. Adding a new shortcut with the
               same name as one that already exists will override the old
               shortcut, but print a warning indicating it is doing so.

    destpath - The path to the shortcut destination.

--remove or -r:
    Removes a shortcut from the go command's list.

    linkname - The name of an existing shortcut.

More information can be found at https://github.com/skrulcik/go-navigator"

    NOT_FOUND="Directory $@ not found. If it is not spelled incorrectly, consider using its absolute path.";

    # Success and error will be binary, different error codes were becoming too
    # cumbersome
    ERR=-1;
    SUCCESS=0;

    # Default shortcut location ("dirname ~/." gets absolute home path)
    DEFAULT_GO_SHORTCUT_DIR="$(dirname ~/.)/.go_shortcuts";

    # Establish where to look for shortcuts
    if [ -z "$GO_SHORTCUT_DIR" ];
    then
        # Try the default go links directory
        GO_SHORTCUT_DIR=$DEFAULT_GO_SHORTCUT_DIR;
        # The default go-links directory should exist if the environment
        # variable is not overridden
        mkdir -p $GO_SHORTCUT_DIR;
    fi
    # Use the full path for the shortcut directory
    GO_SHORTCUT_DIR=$(dirname $GO_SHORTCUT_DIR/.);

    ############################################################################
    # Argument parsing, --help, --add and --delete options
    ############################################################################
    if [[ $# > 0 ]];
    then
        case "$1" in
        # Allow the user to explicitly ask for help
        "--help" | "-h")
            echo "$USAGE";
            echo "$HELP";
            return $SUCCESS;
        ;;
        ########################################################################
        # Add option - quickly create new shortcuts
        ########################################################################
        "--add" | "-a")
            # After --add or -a, there should be two arguments: a shortcut name and a path to a directory
            if [[ $# != 3 ]];
            then
                >&2 echo "Error: expected 2 arguments after $1.";
                >&2 echo "$USAGE";
                return $ERR;
            fi

            # Validate that the shortcut name is alphanumeric
            if [[ "$2" =~ [^a-zA-Z0-9] ]];
            then
                # Shortcut name is not alphanumeric, don't allow this funny business
                echo "Shortcut could not be created. Shortcut names must be alphanumeric."
                return $ERR;
            fi

            # Validate that the directory exists
            if [ ! -d $3 ];
            then
                echo "Shortcut could not be created. Directory $dirpath does not exist.";
                return $ERR;
            fi

            # Combine the shortcut name given in the argument with the shortcut
            # directory to gets its full path
            shortname="$GO_SHORTCUT_DIR/$2";
            # Store the full path of the directory for symbolic linking
            dirpath=$(realpath $3);

            # Validate that the shortcut name is not taken
            if [ -e $shortname ];
            then
                if [ -L $shortname ];
                then
                    # shortcut exists, but can be overridden, so print warning and continue
                    >&2 echo "Warning: Overriding shortcut to $(cd $shortname; pwd)"
                else
                    echo "Shortcut could not be created. $shortname exists and is not a symbolic link.";
                    return $ERR;
                fi
            fi

            # At this point, shortname is a valid shortcut name, and dirpath is
            # a valid directory to link to
            ln -s $dirpath $shortname;
            if [ "$?" = "0" ];
            then
                echo "Created shortcut \"$2\" to $dirpath";
                return $SUCCESS;
            fi

            >&2 echo "Error creating symbolic link \"$shortname\" to $dirpath"
        ;;
        ########################################################################
        # Remove option - Quickly remove shortcuts
        ########################################################################
        "--remove" | "-r")
            if [ $# != 2 ];
            then
                >&2 echo "Error: expected 1 argument after $1.";
                >&2 echo "$USAGE";
                return $ERR;
            fi

            if [ ! -L $GO_SHORTCUT_DIR/$2 ];
            then
                >&2 echo "Error: $2 is not an existing shortcut.";
                >&2 echo "Error: $GO_SHORTCUT_DIR/$2 is not an existing shortcut.";
                return $ERR;
            fi

            rm $GO_SHORTCUT_DIR/$2
            if [ "$?" = "0" ];
            then
                echo "Removed shortcut $2";
                return $SUCCESS;
            fi

            >&2 echo "Error removing symbolic link for shortcut $2";
            return $ERR;
        ;;
        esac

        # If no options are given, there should not be more than one argument
        # passed to the go function
        if [[ $# > 1 ]];
        then
            >&2 echo "$USAGE";
            return $ERR;
        fi
    fi


    ############################################################################
    # Method 1: Try cd
    ############################################################################

    2>/dev/null cd "$@";
    if [ "$?" = "0" ];
    then
        # Regular cd succeeded, print location information then return success
        pwd;
        ls;
        return $SUCCESS;
    fi

    ############################################################################
    # Method 2: Try a symlink
    ############################################################################
    # If the shortcut directory exists, check to see if the argument matches
    # any shortcuts
    if [ -d $GO_SHORTCUT_DIR ];
    then
        # $1 exists, otherwise cd would not have failed
        dest="$GO_SHORTCUT_DIR/$1"
        if [ -L $dest ];
        then
            # Follow the symlink, but check the result because the link itself
            # could be broken
            2>/dev/null cd $dest;
            if [ "$?" = "0" ];
            then
                # Go to the absolute path location (pwd -P works on Mac and
                # Linux, readlink does not
                2>/dev/null cd `pwd -P`;
                # Show the user where they are, and what files are available
                pwd;
                ls;
                # Return successfully
                return $SUCCESS;
            else
                # Symbolic link is broken, delete it so we don't accumulate dead
                # links
                >&2 echo "Warning: Deleting broken symlink $dest";
                rm $dest;
                # Exit with an error code, if there is a cache hit in recent
                # history it would probably be for the broken link
                return $ERR;
            fi
        fi
    else
        # Print a warning message that the go-links directory does not exist
        >&2 echo "Warning: Could not locate shortcut directory $GO_SHORTCUT_DIR"
        # Continue to recent history method
    fi

    ################################################################################
    # Method 3: Recent History
    ################################################################################

    # TODO

    # All three methods failed failed :/
    >&2 echo "go: $@: Directory not found.";
    return $ERR;
}

