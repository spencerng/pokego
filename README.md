# pokego

A deliberately-poor Go wrapper to write Go programs with Pokémon names and types as keywords in the language.

Created by Jannotta House for UChicago's 2020 Scav item #85:

> pokemon-go. No, DEFINITELY not Pokémon Go. pokemon-go: your programming language that is just a crappy wrapper for Go, where all the keywords are now Pokémon. It should be at least able to execute your equivalent of `go run`, but we wouldn’t object if it had more functionality.

## Installation

1. Install the [Go distribution from this link](https://golang.org/doc/install) and add the binaries to your PATH.
2. Clone this repo:
    ````
    git clone https://github.com/spencerng/pokego.git
    ````
3. Build the `pokego` binary and make it executable (on Linux):
    ```
    go build pokego.go
    chmod +x pokego
    ```
4. (optional) Add the `pokego` binary to your PATH to call it from anywhere using `pokego` instead of `./pokego`.

## Usage

`pokego` works identically to the `go` binary, except `.pgo` source files require you to use the keyword names specified in [`dict.json`](./dict.json). To run a program, execute
```
./pokego run myprogram.pgo
```

Additional Pokego commands can be viewed with 
```
./pokego help
```

### Sample Usage

To run our versions of HelloWorld and FizzBuzz programs, execute the following:

```
./pokego run samples/helloworld.pgo
./pokego run samples/mewtwo.pgo
```

The Pokego version of the interpreter can also be compiled with Pokego by running `./pokego build samples/pokego.pgo`

Feel free to create your own programs with the `.pgo` extension!

## Questions?

Contact spencerng [at] uchicago [dot] edu if you have any questions about entering the world of `pokego`.