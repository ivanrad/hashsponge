package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/blake2b"
)

var (
	hashAlg = flag.String("a", "sha256", "hash algorithms: md5, sha1, sha256, sha384, sha512, blake2b")
	quiet   = flag.Bool("q", false, "quiet; don't output error message on hash mismatch")
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage:  %s [-a HASH_ALGORITHM] HASH_VALUE\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Soak up standard input, and write to standard output if checksum matches the given argument.\n\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var hashFn func([]byte) []byte

	switch strings.ToLower(*hashAlg) {
	case "md5":
		hashFn = func(buf []byte) []byte {
			h := md5.Sum(buf)
			return h[:]
		}
	case "sha1":
		hashFn = func(buf []byte) []byte {
			h := sha1.Sum(buf)
			return h[:]
		}
	case "sha256":
		hashFn = func(buf []byte) []byte {
			h := sha256.Sum256(buf)
			return h[:]
		}
	case "sha384":
		hashFn = func(buf []byte) []byte {
			h := sha512.Sum384(buf)
			return h[:]
		}
	case "sha512":
		hashFn = func(buf []byte) []byte {
			h := sha512.Sum512(buf)
			return h[:]
		}
	case "blake2b":
		hashFn = func(buf []byte) []byte {
			h := blake2b.Sum512(buf)
			return h[:]
		}
	default:
		fmt.Fprintf(os.Stderr, "error: unknown hash algorithm: %s\n", *hashAlg)
		os.Exit(1)
	}

	hashArg := flag.Arg(0)
	if hashArg == "" {
		flag.Usage()
		os.Exit(1)
	}

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	inHash := fmt.Sprintf("%x", hashFn(buf))
	if hashArg != inHash {
		if !*quiet {
			fmt.Fprintf(os.Stderr, "error: hash mismatch (input hash: %s)\n", inHash)
		}
		os.Exit(1)
	}

	nwritten := 0
	for nwritten < len(buf) {
		n, err := os.Stdout.Write(buf[nwritten:])
		if err != nil && !errors.Is(err, syscall.EINTR) {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		nwritten += n
	}
}
