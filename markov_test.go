package main // import "cpl.li/go/markov"

import (
	"fmt"
	"math"
	"testing"
)

func TestMarkovBasic(t *testing.T) {
	// generate unigram chain
	chain := NewChain(1)

	// train model
	chain.Add([]string{"I", "ride", "a", "bike"})
	chain.Add([]string{"I", "ride", "the", "bus"})
	chain.Add([]string{"I", "drink", "water"})

	// test invalid probability
	if _, err := chain.TransitionProbability(
		Pair{
			Current: []string{"invalid", "sequence"},
			Next:    "any",
		}); err == nil {

		t.Errorf("computed transition prob on invalid pair size\n")
	}

	// obtain probability
	prob, err := chain.TransitionProbability(
		Pair{
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
