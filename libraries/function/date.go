package function

import (
	"fmt"
	"time"
)

func NowMicro() string {
	return "|" + fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
}
