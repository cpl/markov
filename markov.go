/*
Package markov provides a markov chain implementation which allows you to
"train" a model using any form of text as input. The markov chain will split
the text sequence into pairs and then generate the transition mapping.

A Builder implementation also exists, this can be generated on top of a chain
in order to generate a continuous flow of new words.

MIT License
Copyright (c) 2019 Alexandru-Paul Copil
*/
package markov // import "cpl.li/go/markov"

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

// a random number generator used locally for Markov chain generation.
var localRand *rand.Rand

// initialize the RNG using the system timer in nanoseconds.
func init() {
	localRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Sequence represents a set of 1 or more words used to build a Markov chain.
type Sequence []string

// String returns the words of the sequence as a single string formed of space
// separeted words.
func (s Sequence) String() string {
	return strings.Join(s, " ")
}

func (s Sequence) isValid() bool {
	return len(s) > 1
}

// Pair represents a state transition between a set of 1 or more words and the
// next word in the Sequence.
type Pair struct {
	Current Sequence
	Next    string
}

// Pairs extracts pairs with states of given size transitioning to new words.
func (s Sequence) Pairs(size int) []Pair {
	// invalid sequence size for generating any pairs
	if !s.isValid() {
		return nil
	}

	// clamp size within expected limits
	size = clamp(size, 1, len(s)-1)

	// generate pairs array
	pairs := make([]Pair, len(s)-size)

	// iterate all possible pairs
	for idx := 0; idx < len(s)-size; idx++ {
		// extract and assign the pairs into the return array
		pairs[idx] = Pair{
			Current: s[idx : idx+size],
			Next:    s[idx+size],
		}
	}

	return pairs
}

// transitionMap stores the count of all possible transitions for a state.
type transitionMap map[string]int

// Sum will iterate the values inside a transition map and return their sum.
func (t transitionMap) sum() int {
	sum := 0
	for _, val := range t {
		sum += val
	}
	return sum
}

// Chain represents a Markov chain composed of given length pairs extracted from
// provided sequences.
type Chain struct {
	PairSize int

	frequencyMatrix map[string]transitionMap
}

// NewChain generates a Chain with pairs of given length.
func NewChain(pairSize int) *Chain {
	// clamp
	if pairSize < 1 {
		pairSize = 1
	}

	chain := new(Chain)
	chain.PairSize = pairSize
	chain.frequencyMatrix = make(map[string]transitionMap)
	return chain
}

// Add adds the transition counts to the chain for a given sequence of words.
func (c *Chain) Add(sequence Sequence) {
	// extract pairs from given sequence
	pairs := sequence.Pairs(c.PairSize)

	// if sequence is invalid, paris will be nil and we have nothing to do here
	if pairs == nil {
		return
	}

	// iterate pairs
	for _, pair := range pairs {
		// check if pair was encountered before
		if c.frequencyMatrix[pair.Current.String()] == nil {
			// create new transition map for pair
			c.frequencyMatrix[pair.Current.String()] = make(transitionMap, 0)
		}

		// increment transition occurrence count
		c.frequencyMatrix[pair.Current.String()][pair.Next]++
	}
}

// TransitionProbability returns the probability of transition between the
// current and next state of a pair.
func (c *Chain) TransitionProbability(p Pair) (float64, error) {
	// check the pair sizes are matching
	if c.PairSize != len(p.Current) {
		return 0, errors.New("mismatch pair size from chain and given pair")
	}

	// obtain list of transitions for current state of pair
	transitions := c.frequencyMatrix[p.Current.String()]

	// calculate the sum of all transitions
	sum := transitions.sum()

	// obtain the occurrence to next state
	frequency := transitions[p.Next]

	// compute probability
	return float64(frequency) / float64(sum), nil
}

// Next will give you the next possible token for a certain sequence based on
// a random weighted decision.
func (c *Chain) Next(seed Sequence) string {
	// check for right sequence size
	if len(seed) != c.PairSize {
		return ""
	}

	// get sequence encounter (if none, return empty string)
	arr := c.frequencyMatrix[seed.String()]
	if arr == nil {
		return ""
	}

	// compute total occurrence
	sum := arr.sum()
	// generate random value between 0 and max occurrence
	luck := localRand.Intn(sum)

	// iterate all possibilities until "luck" runs out
	for str, frq := range arr {
		luck -= frq
		if luck <= 0 {
			return str
		}
	}

	// failed to return proper token
	return ""
}
