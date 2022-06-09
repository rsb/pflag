# Description
pflag is a drop-in replacement for Go's flag package, implementing
POSIX/GNU-style --flags.

pflag is compatible with the [GNU extensions to the POSIX recommendations
for command-line options][1]. For a more precise description, see the
"Command-line flag syntax" section below.

[1]: http://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

pflag is available under the same style of BSD license as the Go language,
which can be found in the LICENSE file.


> Much of the code is derived from "[spf13/pflag](https://github.com/spf13/pflag)" by Steve Francia. Thank you 
> for your contribution which allows me build out this package.
> 
> The key difference between the packages is in how the types are handled


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
