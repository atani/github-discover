package trending

import (
	"math"
	"sort"
	"time"

	"github.com/atani/github-discover/internal/github"
)

// ScoredRepo holds a repository with its computed velocity and z-score.
type ScoredRepo struct {
	Repo     github.Repository
	Velocity float64 // stars per day
	ZScore   float64 // standard deviations above mean velocity
}

// CalcVelocity returns the star acquisition rate (stars per day) for a repo.
// Repos created less than 1 day ago are treated as 1 day old to avoid division by zero.
func CalcVelocity(repo github.Repository) float64 {
	days := time.Since(repo.CreatedAt).Hours() / 24
	if days < 1 {
		days = 1
	}
	return float64(repo.StargazersCount) / days
}

// DetectAnomalies returns repos whose star velocity is more than 2 standard
// deviations above the mean. At least 2 repos are required for meaningful
// statistics; otherwise an empty slice is returned.
func DetectAnomalies(repos []github.Repository) []ScoredRepo {
	if len(repos) < 2 {
		return nil
	}

	velocities := make([]float64, len(repos))
	for i, r := range repos {
		velocities[i] = CalcVelocity(r)
	}

	mean, stddev := meanStddev(velocities)
	if stddev == 0 {
		return nil
	}

	threshold := mean + 2*stddev

	var results []ScoredRepo
	for i, v := range velocities {
		if v > threshold {
			results = append(results, ScoredRepo{
				Repo:     repos[i],
				Velocity: v,
				ZScore:   (v - mean) / stddev,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Velocity > results[j].Velocity
	})

	return results
}

// RankByVelocity sorts repos by star velocity in descending order.
func RankByVelocity(repos []github.Repository) []ScoredRepo {
	if len(repos) == 0 {
		return nil
	}

	scored := make([]ScoredRepo, len(repos))
	for i, r := range repos {
		scored[i] = ScoredRepo{
			Repo:     r,
			Velocity: CalcVelocity(r),
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Velocity > scored[j].Velocity
	})

	return scored
}

func meanStddev(vals []float64) (float64, float64) {
	n := float64(len(vals))
	if n == 0 {
		return 0, 0
	}

	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / n

	var variance float64
	for _, v := range vals {
		diff := v - mean
		variance += diff * diff
	}
	variance /= n

	return mean, math.Sqrt(variance)
}
