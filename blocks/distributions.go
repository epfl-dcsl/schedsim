package blocks

import (
	"math"
	"math/rand"
)

type randDist interface {
	getRand() float64
}

// Deterministic Distribution
type deterministicDistr struct {
	d float64
}

func newDeterministicDistr(d float64) *deterministicDistr {
	return &deterministicDistr{d}
}

func (distr *deterministicDistr) getRand() float64 {
	return distr.d
}

// Exponential Distribution
type exponDistr struct {
	lambda float64
}

func newExponDistr(l float64) *exponDistr {
	return &exponDistr{l}
}

func (distr *exponDistr) getRand() float64 {
	return float64(rand.ExpFloat64() / distr.lambda)
}

// LogNormal Distribution
type lGDistr struct {
	mu    float64
	sigma float64
}

func newLGDistr(mu, sigma float64) *lGDistr {
	return &lGDistr{mu, sigma}
}

func (distr *lGDistr) getRand() float64 {
	z := rand.NormFloat64()
	s := math.Exp(distr.mu + distr.sigma*z)
	return s
}

// Bimodel Distribution
type biDistr struct {
	v1    float64
	v2    float64
	ratio float64
}

func newBiDistr(v1, v2, ratio float64) *biDistr {
	return &biDistr{v1, v2, ratio}
}

func (distr *biDistr) getRand() float64 {
	if rand.Float64() > distr.ratio {
		return distr.v2
	}
	return distr.v1
}
