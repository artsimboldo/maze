package maze

import (
	"testing"
	"fmt"
)

func TestPrim(t *testing.T) {
	prim0result := " _________\n|___   ___|\n|___   ___|\n| |_   ___|\n|_   ___  |\n|_______|_|\n"
	prim := Prim{Width:5, Height:5, Seed:int64(0)}
	if _, err := prim.Generate(); err == nil {
		fmt.Println(prim.String())
		if prim.String() != prim0result {
			t.Fail()
		}
	} else {
		t.Errorf("%v", err)
	}
}