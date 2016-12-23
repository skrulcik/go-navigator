#!/bin/bash
#
# go - A faster way to change directories.
#
# The go command must be implemented as a function, rather than a normal executable
#
# (c) Scott Krulcik 2017 MIT License
#

function go() {
    USAGE="Usage:\n\tgo [dest]"
    NOT_FOUND="Directory $@ not found. If it is not spelled incorrectly, consider using its absolute path."
    # Default shortcut location ("dirname ~/." gets absolute home path)
    DEFAULT_GO_SHORTCUT_DIR="$(dirname ~/.)/.go_shortcuts"


    # No more than 1 argument is allowed
    if [[ $# > 1 ]];
    then
        echo "$USAGE";
        return -1;
    fi



    ############################################################################
    # Method 1: Try cd
    ############################################################################

    2>/dev/null cd "$@";
    if [ "$?" = "0" ];
    then
        # Regular cd succeeded, return success
        return 0;
    fi

    ############################################################################
    # Method 2: Try a symlink
    ############################################################################

    # Establish where to look for shortcuts
    if [ -z $GO_SHORTCUT_DIR ];
    then
        # Try the default go links directory
        GO_SHORTCUT_DIR=$DEFAULT_GO_SHORTCUT_DIR;
        # The default go-links directory should exist if the environment
        # variable is not overridden
        mkdir -p $GO_SHORTCUT_DIR;
    fi

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
                return 0;
            else
                # Symbolic link is broken, delete it so we don't accumulate dead
                # links
                >&2 echo "Warning: Deleting broken symlink $dest";
                rm $dest;
                # Exit with an error code, if there is a cache hit in recent
                # history it would probably be for the broken link
                return -2;
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
    return -3;
}
