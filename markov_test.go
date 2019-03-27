package markov_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"testing"

	"cpl.li/go/markov"
)

// This example shows a general usecase for the Markov Chain and the builder. It
// takes input from `stdin` and trains the markov chain then generates a given
// number of words nd prints out the fully generated string. The flags can
// configure the max number of words to generate and the sequence pairing to
// be used when "training" the markov chain.
func Example_basic() {
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

func TestMarkovBasic(t *testing.T) {
	// generate unigram chain
	chain := markov.NewChain(1)

	// train model
	chain.Add([]string{"I", "ride", "a", "bike"})
	chain.Add([]string{"I", "ride", "the", "bus"})
	chain.Add([]string{"I", "drink", "water"})

	// test invalid probability
	if _, err := chain.TransitionProbability(
		markov.Pair{
			Current: []string{"invalid", "sequence"},
			Next:    "any",
		}); err == nil {

		t.Errorf("computed transition prob on invalid pair size\n")
	}

	// obtain probability
	prob, err := chain.TransitionProbability(
		markov.Pair{
			Current: []string{"I"},
			Next:    "ride",
		})

	// fail on error
	if err != nil {
		t.Error(err)
	}

	// fail on wrong probability
	real := math.Round(prob*100) / 100
	expected := 0.67
	if expected != real {
		fmt.Printf("expected %f got %f\n", expected, real)
		t.Fail()
	}

	// create builder from chain
	builder := chain.NewBuilder([]string{"I"})
	if builder.Generate(2) != 2 {
		t.Errorf("failed to generate 2 new words\n")
	}
}
