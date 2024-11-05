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
$ ./some-long-process | teecp | grep "dodongo"
```

For each other terminal (assumes that `alias teecp` has been applied):

```sh
$ teecp --client | grep "bomb"  | teecp --port 6668
```

```sh
$ teecp --client --port 6668 | wc -l
```

```sh
$ teecp --client | grep "[Ll]ink"
```

## Current status

- [ ] Create executable `teecp` to allow better utility experience
- [ ] Add `asdf` plugin for easiness of use
- [ ] SSL [#4](https://github.com/jeffque/teecp/issues/4)
- [ ] Auth [#5](https://github.com/jeffque/teecp/issues/5)
