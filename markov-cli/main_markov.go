package main // import "cpl.li/go/markov/markov-cli"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cpl/markov"
)

func main() {
	// handle flags
	maxWords := flag.Int("words", 100, "max words to generate (default 100)")
	pairSize := flag.Int("pairs", 2, "size of a word pair (default 2)")
	flag.Parse()

	c := markov.NewChain(*pairSize) // create markov chain

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
