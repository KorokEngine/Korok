package effect

import "testing"

func TestFireSimulator(t *testing.T) {
	fire := NewFireSimulator(1024)
	fire.Initialize()

	for i := 0; i < 10; i++ {
		fire.Simulate(1.0/60)
	}

	
}
