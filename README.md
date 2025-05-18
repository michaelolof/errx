# errx

Error reporting and tracking using statically defined stamps or trace ids

## Installation
To install run the following command
```sh
$ go get -u github.com/michaelolof/errx
```

## Motivation
A lot has been said about error handling in Go, but one thing that can be agreed on is that proper error reporting requires stack traces. Whether through error wrapping with human-readable texts or runtime reflection to get the file paths and line numbers. Developers want to see how their error got propagated.

A good error handling system should be able to meet the following requirements:
- What went wrong. (This is the error message)
- Where it went wrong. (Stack trace, unique identifiers, human-readable texts)
- What type of error it is.
- Store/Access any useful data needed for processing.

This [article](https://dev.to/michaelolof_/golang-error-handling-a-practical-and-robust-solution-3am1) goes into detail on my motivations for the `errx` library
<br />


## Examples

### Creating and Wrapping Errors
```go
import "github.com/michaelolof/errx"


func failerOne() error {
    ...
    return errx.New(1745397000, "something went wrong")
    ...
}
```

This will generate an error message that looks like this
```text
[ts 1745397000] something went wrong
```


You can also decorate an existing error with a new stamp
```go
import "github.com/michaelolof/errx"


func failerTwo() error {
    ...
    err := failerOne()
    if err != nil {
        return errx.Wrap(1745397994, err)
    }
    ...
}
```
This will generate an error message that looks like this:
```txt
[ts 1745397994]; [ts 1745397000] something went wrong
```
The library enforces that the stamps passed to the `New` or `Wrap` functions are literal integers. So this won't work
```go
func handleErr(ts int, msg string) error {
    return errx.New(ts, msg)
}
```
<br />

### Error Kinds and Data Kinds
There are times when we need to mark our errors depending on our use case. With `errx` its done like so:
```go
var (
    NotFoundErr  = errx.Kind("notfound")
    InvalidNoErr = errx.DataKind[int]("invalidno")
)
```

They can be used like so:
```go
import "github.com/michaelolof/errx"


func failerOne() error {
    ...
    return errx.NewKind(1745397000, NotFoundErr, "something went wrong")
    ...
}

func failerTwo() error {
    var v int = 2
    ...
    err := failerOne()
    if err != nil {
        return errx.WrapKind(1745397994, InvalidNoErr(v), err)
    }
    ...
}
```
This will generate an error message that looks like this:
```sh
[ts 1745397994 kind invalidno data 2]; [ts 1745397000 kind notfound] something went wrong
```
In `errx` error kinds are just strings. And we can check based on the error kind like so:
```go
if errx.IsKind(err, NotFoundErr) {
    // error has a kind of notfound
}

if err.IsDataKind(err, InvalidNoErr) {
    // error has a kind of invalidno
}
```

For data kinds, we can access the data using the `FindData` function.
```go
if v, ok := errx.FindData(err, InvalidNoErr); ok {
    fmt.Println(v) // v is int(2)
}
```

In `errx` every information about the error - stamp, kind, data, message are structured as part of the error string. This means your error string tells the full story about your errors. Consequently this also means you can build* back your error object from the strings by calling`ParseStampedError`
<br/>
<br/>
It might be tempting define error kinds every time you create or wrap an error, but in practice that's usually a bad idea. A general rule to decide when you need an error kind is if you need it for decisioning at some later point in your application. 
<br/>
Essentially if you're not going to check on it using `IsKind` or `IsDataKind` or retrieve data from it using `FindData` just stick to basic error creation or wrapping and don't define kinds for them.

## Why Stamps?
You might be hesitant to add random integers alongside your errors and might be wondering why not just use stack traces and pay the reflection penalty. This is perfectly valid and fine. I've used all before. No wrapping, wrapping with texts, stack traces and now stamps.
<br /><br />
These are sone of the reasons i've settled on stamps
- A simple and intituive system. Create new errors with stamps. Wrap existing errors with stamps
- Stamps are unique. This is very important for tracing.
- Stamps are shorter compared to file paths and line numbers that come with stack traces.
- Stamps are context free which means they're immune to changes and refactors
- Stamps are easier to log
- Stamps can be safely exposed to the client/public. I'm perfectly fine with adding a stamp as part of my API error response cause to the outside world, its meaningless.
- Stamps don't rely on runtime reflection, hereby pay no performance penalty.
- Stamps are suprisingly easy the generate. Using an [editor snippet](https://code.visualstudio.com/docs/editing/userdefinedsnippets#_variables), it takes me less time to generate the stamp where needed and move on than typing the perfect human-readable error context which needs to be meainingful, unique and still generic for the place where it's used.
