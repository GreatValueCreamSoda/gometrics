package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/GreatValueCreamSoda/gometrics/video/metrics"
)

const (
	jodA   = 0.0439569391310215
	jodExp = 0.9302042722702026
)

func jod(a float64) float64 {
	return 10.0 - jodA*math.Pow(a, jodExp)
}

func inverseJOD(a float64) float64 {
	return math.Pow((10.0-a)/jodA, 1.0/jodExp)
}

// ────────────────────────────────────────────────────────────────────────────────
// Metric presentation abstraction
// ────────────────────────────────────────────────────────────────────────────────

type MetricPresenter interface {
	DisplayName() string
	// TransformForStats: space in which min/avg/median/stddev are computed
	TransformForStats(v float64) float64
	// TransformForDisplay: space in which values are shown to the user
	TransformForDisplay(v float64) float64
}

type DefaultPresenter struct {
	name string
}

func (p DefaultPresenter) DisplayName() string {
	return p.name
}

func (p DefaultPresenter) TransformForStats(v float64) float64 {
	return v
}

func (p DefaultPresenter) TransformForDisplay(v float64) float64 {
	return v
}

type CVVDPPresenter struct{}

func (p CVVDPPresenter) DisplayName() string {
	return metrics.CVVDPName
}

func (p CVVDPPresenter) TransformForStats(v float64) float64 {
	return inverseJOD(v)
}

func (p CVVDPPresenter) TransformForDisplay(v float64) float64 {
	return jod(v)
}

// ────────────────────────────────────────────────────────────────────────────────
// Main printing logic
// ────────────────────────────────────────────────────────────────────────────────

func getPresenter(name string) MetricPresenter {
	if name == metrics.CVVDPName {
		return CVVDPPresenter{}
	}
	return DefaultPresenter{name: name}
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

func printMetricSummary(name string, rawValues []float64) {
	presenter := getPresenter(name)

	// Transform all values into the space where we want statistics
	values := make([]float64, len(rawValues))
	for i, v := range rawValues {
		values[i] = presenter.TransformForStats(v)
	}

	n := len(values)
	if n == 0 {
		return
	}

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
	variance /= float64(n) // population stddev; use n-1 for sample if preferred
	stddev := math.Sqrt(variance)

	// Output ─ all displayed values go through TransformForDisplay
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, presenter.DisplayName())
	fmt.Fprintln(os.Stderr, strings.Repeat("-", len(presenter.DisplayName())))

	fmt.Fprintf(os.Stderr, "  min     : %.6f\n", presenter.TransformForDisplay(min))
	fmt.Fprintf(os.Stderr, "  max     : %.6f\n", presenter.TransformForDisplay(max))
	fmt.Fprintf(os.Stderr, "  average : %.6f\n", presenter.TransformForDisplay(avg))
	fmt.Fprintf(os.Stderr, "  median  : %.6f\n", presenter.TransformForDisplay(median))
	fmt.Fprintf(os.Stderr, "  stddev  : %.6f\n", presenter.TransformForDisplay(stddev))
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

// ────────────────────────────────────────────────────────────────────────────────
// Correlation implementations (unchanged)
// ────────────────────────────────────────────────────────────────────────────────

type CorrelationMethod struct {
	Name string
	Fn   func(x, y []float64) float64
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
		ranks[pairs[i].index] = float64(i + 1) // typically ranks start from 1
	}

	return ranks
}
