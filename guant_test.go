package guant

import (
	"fmt"
	"gonum.org/v1/gonum/stat/distuv"
	"testing"
	"time"
)

func getTestDerivative() Derivative {
	//Baseline Values to test against
	return Derivative{
		N: distuv.Normal{Mu: 0, Sigma: 1},
		S: 183.25,
		K: 190.00,
		R: DefaultRfir(),
		T: 0.01643835616438356,
	}
}

func TestD1(t *testing.T) {
	ltd := getTestDerivative()
	ltd.Sigma = 0.23979117049551366
	d1 := ltd.d1()
	ans := "-1.157672"
	d := fmt.Sprintf("%f", d1)
	if d != ans {
		t.Errorf("D1 Calculation was incorrect, got %s, want %s", d, ans)
	}
}
func TestD2(t *testing.T) {
	ltd := getTestDerivative()
	ltd.Sigma = 0.23979117049551366
	var d2  = ltd.d2()
	ans := "-1.188416"
	d := fmt.Sprintf("%f", d2)
	if d != ans {
		t.Errorf("D2 Calculation was incorrect, got %s, want %s", d, ans)
	}
}
func TestDefaultRfir(t *testing.T) {
	dfr := DefaultRfir()
	ans := 0.006600000000000001
	if dfr != ans {
		t.Errorf("Default RFIR incorrect, got %f, want %f", dfr, ans)
	}
}

func TestBlackScholesCall(t *testing.T) {
	ltd := getTestDerivative()
	ltd.Sigma = 0.23979117049551366
	bsc := blackScholesCall(ltd)
	ans := 0.340000008961411
	if bsc != ans {
		t.Errorf("Black-Scholes Call Calculation was incorrect, got %f, want %f", bsc, ans)
	}
}

func TestBlackScholesPut(t *testing.T) {
	ltd := getTestDerivative()
	ltd.K = 170
	ltd.Sigma = 0.3570574310795081
	ltd.Put = true
	bsp := blackScholesPut(ltd)
	ans := 0.1700000113961444
	if bsp != ans {
		t.Errorf("Black-Scholes Call Calculation was incorrect, got %f, want %f", bsp, ans)
	}
}

func TestBlackScholes(t *testing.T) {
	ltd := getTestDerivative()
	//Test Call
	ltd.Sigma = 0.23979117049551366
	bsc := BlackScholes(ltd)
	bscAns := 0.340000008961411
	//Test Put
	ltd.K = 170.00
	ltd.Sigma = 0.3570574310795081
	ltd.Put = true
	bsp := BlackScholes(ltd)
	bspAns := 0.1700000113961444
	if bsc != bscAns || bsp != bspAns {
		t.Errorf("Black-Scholes Calculation was incorrect\nCall: got %f, want %f\nPut: got %f, want %f\n", bsc, bscAns, bsp, bspAns)
	}
}

func TestNewtonRaphson(t *testing.T) {
	ltd := getTestDerivative()
	callMid := 0.33999999999999997
	nrcAns := 0.23979117049551366
	nrc := NewtonRaphson(ltd, callMid)
	//Calculate Put IV
	putMid := 0.17
	nrpAns := 0.3570574310795081
	ltd.K = 170
	ltd.Put = true
	nrp := NewtonRaphson(ltd, putMid)
	if nrc != nrcAns || nrp != nrpAns {
		t.Errorf("Newton-Raphson Calculation was incorrect\nCall: got %f, want %f\nPut: got %f, want %f\n", nrc, nrcAns, nrp, nrpAns)
	}
}

// Test against 6/365 value
func TestTimeToExpiry(t *testing.T) {
	currTime := time.Now()
	currentDate := currTime.Format("2006-01-02")
	daysLater := currTime.AddDate(0, 0, 5)
	expiry := daysLater.Format("2006-01-02")
	tte, err  := TimeToExpiry(currentDate, expiry)
	if err != nil {
		t.Errorf("Time to expiry failed parsing %v", err)
	}
	ans := 0.01643835616438356
	if tte != ans {
		t.Errorf("TimeToExpiry incorrect calculation, got %f, wanted %f", tte, ans)
	}
}

func TestDate(t *testing.T) {
	_, err := date("2020-06-05")
	if err != nil {
		t.Errorf("%v", err)
	}
}
