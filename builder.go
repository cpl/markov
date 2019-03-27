package main // import "cpl.li/go/markov"

import (
	"strings"
)

// Builder spawns from a Markov chain and using the generated tokens from it,
// creates sequences of words.
type Builder struct {
	words Sequence
	chain *Chain
}

// NewBuilder creates a Markov sequence builder form the current chain.
func (c *Chain) NewBuilder(seed Sequence) Builder {
	// check for no given seed, invalid size seed or non-existent sequence
	if seed == nil || len(seed) < c.PairSize ||
		c.frequencyMatrix[seed.String()] == nil {

		// pick a random sequence from the chain as start
		rn := localRand.Intn(len(c.frequencyMatrix))
		for sequence := range c.frequencyMatrix {
			rn--
			if rn == 0 {
				seed = strings.Split(sequence, " ")
			}
		}
	}

	return Builder{
		chain: c,
		words: seed[:],
	}
}

// String will return the word sequence as a single string of all words
// separated by spaces.
func (b *Builder) String() string {
	return b.words.String()
}

// Generate will tell the builder to poll the markov chain for at most `count`
// new words to append to the builders sequence. The function will return the
// real number of generated words, 0 meaning no new words could be generated.
func (b *Builder) Generate(count int) int {
	initialCount := count

	// iterate for `count` generations
	for count > 0 {
		// obtain new word from chain and check for empty string,
		// always use the last chain.PairSize words from the builder
		// sequence as the seed for the chain generator
		next := b.chain.Next(b.words[len(b.words)-b.chain.PairSize:])
		if next == "" {
			break
		}

		// append new word to sequence
		b.words = append(b.words[:], next)

		count--
	}

	// return real number of generated words
	return initialCount - count
}
