package crypto

import "github.com/ZZMarquis/gm/sm3"


// SumSM3 returns the SM3 hash of the data.
func SumSM3(data []byte) []byte {
	d := sm3.New()
	d.Write(data)
	return d.Sum(nil)
}
