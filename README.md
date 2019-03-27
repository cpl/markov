# Markov Chains in Go

[![Go Report Card](https://goreportcard.com/badge/cpl.li/go/markov)](https://goreportcard.com/report/cpl.li/go/markov)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?)](https://godoc.org/cpl.li/go/markov)

If you don't know what a Markov Chain is I recommend reading up on it and checking out this [visual explanation](http://setosa.io/ev/markov-chains/). To put it simply, Markov Chain represent probabilistic state changes inside a Finite-State-Machine. The state transition probabilities can be easily represented as a matrix but in practice (software) that leaves lots of entries being 0 and taking up memory, so a more elegant solution is a nested hash map.

In the code above the following are used:

```go
// represents a grouping of individual words
// eg: []string{"I", "am", "Alex"}, this can be extracted
// from an original string of any form or shape:
// "I am Alex", "I   am  Alex", "I:am:Alex"
// and it's all up to the caller to split their strings into sequences
type Sequence []string
```


```go
// a pairs represents a possible transition between a sequence of n words
// and the next (single) word
// the Current sequence must be of an equal lenght to the chain pair size
// meaning you can't have some transitions for 2-grouped words and 1-grouped words
type Pair struct {
	Current Sequence
	Next    string
}
```

```go
// by having a
type transitionMap map[string]int
// and then nested inside
frequencyMatrix map[string]transitionMap
// we generate our mapping of all encountered
// sequences to their respective next word
// and the number of times this occurs
```

Once your wrap your head around these structures, the rest of the functions are easy to understand.

## Download & Install

If you have Go installed, you can simply run:

```shell
go get cpl.li/go/markov
```

## Usage

I provided an example main function with `stdin` input and basic flag parsing for generating n words from the input data.

```go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
  
  "cpl.li/go/markov"
)

func main() {
	// handle flags
	maxWords := flag.Int("words", 100, "max words to generate (default 100)")
	pairSize := flag.Int("pairs", 2, "size of a word pair (default 2)")
	flag.Parse()

	c := NewChain(*pairSize) // create markov chain

	// read stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// give data as sequence to chain model
	c.Add(strings.Fields(string(data)))

	b := c.NewBuilder(nil)             // create builder on top of chain
	b.Generate(*maxWords - c.PairSize) // generate new words
	fmt.Println(b.String())            // print end product
}
```






