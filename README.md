# Description
pflag is a drop-in replacement for Go's flag package, implementing
POSIX/GNU-style --flags.

pflag is compatible with the [GNU extensions to the POSIX recommendations
for command-line options][1]. For a more precise description, see the
"Command-line flag syntax" section below.

[1]: http://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

pflag is available under the same style of BSD license as the Go language,
which can be found in the LICENSE file.

This package is heavily influenced by [spf13's pflag](https://github.com/spf13/plag). 
I wanted to use all the same concepts but keep all values a string since I have
other packages that will convert the string into its type based of a struct annotation


## Installation

pflag is available using the standard `go get` command.

Install by running:

    go get github.com/rsb/pflag

Run tests by running:

    go test github.com/rsb/pflag

## Usage

pflag is a drop-in replacement of Go's native flag package. If you import
pflag under the name "flag" then all code should continue to function
with no changes.

``` go
import flag "github.com/rsb/pflag"
```
