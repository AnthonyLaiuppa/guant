package guant

import (
	"fmt"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

//Details of the Derivative we want to price
type Derivative struct {
	N     distuv.Normal
	S     float64 // S - stock price
	K     float64 // K - strike price
	R     float64 //r - risk free interest rate
	Sigma float64 // sigma - std deviation of log returns (Implied volatility)
	T     float64 //T - time to exercise date in years
	Put   bool    //Set to True if Put
}

// Risk-free Interest Rate (Default 10-Year US Treasury Bond)
func DefaultRfir() float64 {
	return 0.006600000000000001
}

// D1 Value needed for Black-Scholes
func (a *Derivative) d1() float64 {
	return (math.Log(a.S/a.K) + ((a.R + (a.Sigma*a.Sigma)/2) * a.T)) / (a.Sigma * math.Sqrt(a.T))
}

// D2 Value needed for Black-Scholes
func (a *Derivative) d2() float64 {
	return (math.Log(a.S/a.K) + ((a.R - (a.Sigma*a.Sigma)/2) * a.T)) / (a.Sigma * math.Sqrt(a.T))
}

// Calculate Call value using Black-Scholes
func blackScholesCall(a Derivative) float64 {
	return (a.S * a.N.CDF(a.d1())) - ((a.K * math.Exp(-a.R*a.T)) * a.N.CDF(a.d2()))
}

// Calculate Put value using Black-Scholes
func blackScholesPut(a Derivative) float64 {
	return (a.N.CDF(-1 * a.d2()) * (a.K * math.Exp(-a.R * a.T))) - (a.S * a.N.CDF(-1 * a.d1()))
}

// Black-Scholes Method returns call or put value
func BlackScholes(a Derivative) float64 {
	if a.Put {
		return blackScholesPut(a)
	} else {
		return blackScholesCall(a)
	}
}

// Iterative Method of calculating implied volatility in lieu of an ideal closed form solution
func NewtonRaphson(a Derivative, C0 float64) float64 {
	var tol = 0.001
	var epsilon float64 = 1

	//  Variables to log and manage number of iterations
	count := 0
	maxIter := 1000

	// Starting Point for our IV calculations
	var vol = 0.50

	for epsilon > tol {
		count += 1
		if count >= maxIter {
			fmt.Println("Max iterations reached")
			break
		}

		//  Track previous value to calculate percent change
		origVol := vol

		//Update our Ingredients sigma and calculate the option value
		a.Sigma = vol
		var functionValue float64
		if a.Put {
			functionValue = blackScholesPut(a) - C0
		} else {
			functionValue = blackScholesCall(a) - C0
		}
		var vega = a.S * a.N.Prob(a.d1()) * math.Sqrt(a.T)
		vol = -functionValue/vega + vol
		epsilon = math.Abs((vol - origVol) / origVol)
	}

	return vol
}

// Return time to expiry in years
func TimeToExpiry(provided string, expiry string) (float64, error) {
	a, err := date(provided)
	if err != nil {
		return 0, err
	}
	b, err := date(expiry)
	if err != nil {
		return 0, err
	}
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	//Time to exercise date in years = DBE/Years
	return float64(days+1) / float64(365), nil
}

func date(s string) (time.Time, error) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return d, nil
	}
	return d, err
}
