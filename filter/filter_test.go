package filter

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Printf(
		"filter result : %+v \n",
		Filter([]byte("asd,asd"), []string{"point", "coma"}))
}
