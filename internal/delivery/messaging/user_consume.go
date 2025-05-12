package messaging

import "fmt"

func LoginConsume(msg []byte) {
	fmt.Println(string(msg))
}
