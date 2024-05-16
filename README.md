# errx

Errx is a simple library that allows you easily create, handle and track golang errors using statically generated timestamps or trace ids

## Installation
To install run the following command
```sh
$ go get -u github.com/michaelolof/errx
```

## Motivation
Proper error handling requires you to be able to answer some basic questions about your error. what went wrong?, where it went wrong?, the kind/category of error? and any trace/recovery data to pass along with your error.

To solve this, this library enforces the definition of the following, when defining your errors
- Timestamp/Trace id - Every function to create an error requires a timestamp
- error message/error object - You can either create an new error or wrap existing ones
- error kind - an error object that defines the category of error
- data - optionally pass a recorvery data along with your error object

## Examples

You can create a new timestamped error
```go
import "github.com/michaelolof/errx"


func failerOne() error {
    ...
    return errx.New(1712851695469, "something went wrong")
    ...
}
```

This will generate an error message that looks like this
```text
[ts 1712851695469] something went wrong
```


You can also decorate an existing error with a new timestamp
```go
import "github.com/michaelolof/errx"


func failerTwo() error {
    ...
    err := failerOne()
    if err != nil {
        return errx.Wrap(1712484857431, err)
    }
    ...
}
```

This will generate an error message that looks like this:
```txt
[ts 1712484857431]; [ts 1712851695469] something went wrong
```

**Note** that the timestamp is statically generated and hardcoded when calling the New(...) function. This allows them act as a trace id when locating the source of your errors
<br />
Generating them dynamically means they cannot be used to trace errors in your code.

## Why Timestamps?
The nature of Golang error handling is one which errors are passed/returned as values across your call stack. This means no stack traces.
<br /><br />
Golang solves this by wrapping its errors with decorator texts (A descriptive text that help identify/trace where your errors are occurring). E.g
```go
data, err := logs.readFile()
if err != nil {
    return nil, fmt.Errorf("error reading logstash file: %w", err)
}
```
<br />
However in practice; I've found the use of wrapper/decorator texts to be tedious and problematic. Most online examples describe them as a way to add "more context" to an existing error, but this is usually unnecessary information. The error message already has all the context you need to know about what went wrong.
<br/><br/>
Timestamps are encouraged because:

- They are **guaranteed to be unique**. This is very important for tracing.
- They are context free. This is also good becuase they don't add noise to your error messages.
- They are immune to code refractorings/reorganization since they don't try to describe what the function is doing
- Easy to generate. Rather than thinking up descriptive texts that are unique enough and still prone to change, timestamps can be autogenerated using snippets, terminals etc.


<!-- ## Performant
The `errx` module provides a way to have stack traces without relying on reflection. This means better performance.<br>
It should be noted though that the library does use reflection in some cases when parsing tr -->

## Editor Support
You can automate the creation of timestamps on VSCode by using the [HypserSnips](https://marketplace.visualstudio.com/items?itemName=draivin.hsnips) extension that allows you execute javascript in your code snippets.