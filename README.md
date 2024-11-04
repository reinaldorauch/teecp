# TeeCP

The power of TCP over `tee`

## Design goals

Allow to create a replacement of `tee`, capable of receiving new clients
on the fly. Does not limit the process of stdout to be solely to be the
only stream to apply further processing.

One should be able to run further processing in a long running process on
the fly:

```sh
$ alias teecp='go run github.com/jeffque/teecp@latest'
$ ./some-long-process | teecp --server=:6667 | grep "dodongo"
```

For each other terminal (assumes that `alias teecp` has been applied):

```sh
$ teecp --client=localhost:6667 | grep "bomb"  | teecp --server=:6668
```

```sh
$ teecp --client=localhost:6668 | wc -l
```

```sh
$ teecp --client=localhost:6667 | grep "[Ll]ink"
```

## Current status

- [ ] Command line 
- [ ] User selected port (for now is 6667)
- [ ] Alternate between client and server with something more beautiful
