# hashsponge

Soak up standard input, and write to standard output if checksum matches the
given argument. That's all.

## example

```
$ echo "hello, world!" | ./hashsponge 4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc
hello, world!

$ echo "hello, world!" | ./hashsponge 12345678
error: hash mismatch (input hash: 4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc)
```
