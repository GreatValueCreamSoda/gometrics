package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type CorrelationMethod struct {
	Name string
	Fn   func(x, y []float64) float64
}

func printSummary(scores map[string][]float64) {
	if len(scores) == 0 {
		fmt.Fprintln(os.Stderr, "No scores to report")
		return
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Metric summary")
	fmt.Fprintln(os.Stderr, "==============")

	names := make([]string, 0, len(scores))
	for name := range scores {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		values := scores[name]
		if len(values) == 0 {
			continue
		}
		printMetricSummary(name, values)
	}

	if len(names) > 1 {
		methods := defaultCorrelationMethods()
		printCorrelations(scores, names, methods)
	}
}

func printMetricSummary(name string, values []float64) {
	n := len(values)

	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	min := sorted[0]
	max := sorted[n-1]

	var sum float64
	for _, v := range values {
		sum += v
	}
	avg := sum / float64(n)

	var median float64
	if n%2 == 1 {
		median = sorted[n/2]
	} else {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	}

	var variance float64
	for _, v := range values {
		d := v - avg
		variance += d * d
	}
	variance /= float64(n)
	stddev := math.Sqrt(variance)

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, name)
	fmt.Fprintln(os.Stderr, strings.Repeat("-", len(name)))
	fmt.Fprintf(os.Stderr, "  min     : %.6f\n", min)
	fmt.Fprintf(os.Stderr, "  max     : %.6f\n", max)
	fmt.Fprintf(os.Stderr, "  average : %.6f\n", avg)
	fmt.Fprintf(os.Stderr, "  median  : %.6f\n", median)
	fmt.Fprintf(os.Stderr, "  stddev  : %.6f\n", stddev)
}

func defaultCorrelationMethods() []CorrelationMethod {
	return []CorrelationMethod{
		{"Pearson", pearsonCorrelation},
		{"Spearman", spearmanCorrelation},
		{"Kendall", kendallTauCorrelation},
	}
}

func printCorrelations(
	scores map[string][]float64,
	names []string,
	methods []CorrelationMethod,
) {
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	formatStr := fmt.Sprintf("  %%-%ds ↔ %%-%ds : %% .6f\n", maxLen, maxLen)

	for _, method := range methods {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, method.Name, "correlations")
		fmt.Fprintln(os.Stderr, strings.Repeat("=", len(method.Name)+13))

		for i := 0; i < len(names); i++ {
			for j := i + 1; j < len(names); j++ {
				a, b := names[i], names[j]
				x, y := scores[a], scores[b]

				if len(x) == 0 || len(y) == 0 || len(x) != len(y) {
					continue
				}

				r := method.Fn(x, y)
				fmt.Fprintf(os.Stderr, formatStr, a, b, math.Abs(r))
			}
		}
	}
}

func pearsonCorrelation(x, y []float64) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return 0
	}

	var sumX, sumY float64
	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
	}

	meanX := sumX / float64(n)
	meanY := sumY / float64(n)

	var num, denomX, denomY float64
	for i := 0; i < n; i++ {
		dx := x[i] - meanX
		dy := y[i] - meanY
		num += dx * dy
		denomX += dx * dx
		denomY += dy * dy
	}

	denom := math.Sqrt(denomX * denomY)
	if denom == 0 {
		return 0
	}

	return num / denom
}

func spearmanCorrelation(x, y []float64) float64 {
	rx := ranks(x)
	ry := ranks(y)
	return pearsonCorrelation(rx, ry)
}

func kendallTauCorrelation(x, y []float64) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return 0
	}

	var concordant, discordant float64

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := x[i] - x[j]
			dy := y[i] - y[j]

			if dx*dy > 0 {
				concordant++
			} else if dx*dy < 0 {
				discordant++
			}
		}
	}

	denom := float64(n*(n-1)) / 2
	if denom == 0 {
		return 0
	}

	return (concordant - discordant) / denom
}

func ranks(values []float64) []float64 {
	type pair struct {
		value float64
		index int
	}

	n := len(values)
	pairs := make([]pair, n)
	for i, v := range values {
		pairs[i] = pair{v, i}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value < pairs[j].value
	})

	ranks := make([]float64, n)
	for i := 0; i < n; i++ {
		ranks[pairs[i].index] = float64(i)
	}

	return ranks
}
