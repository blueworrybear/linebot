package service

import (
	"testing"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func TestSimilarity(t *testing.T) {
	similarity := strutil.Similarity("對", "對", metrics.NewHamming())
	t.Log(similarity)
	similarity = strutil.Similarity("不對", "對", metrics.NewHamming())
	t.Log(similarity)
}
