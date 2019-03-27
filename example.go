package markov // import "cpl.li/go/markov"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ExampleMain shows a general usecase for the Markov Chain and the builder. It
// takes input from `stdin` and trains the markov chain then generates a given
// number of words nd prints out the fully generated string. The flags can
// configure the max number of words to generate and the sequence pairing to
// be used when "training" the markov chain.
func ExampleMain() {
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
