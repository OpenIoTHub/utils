package main

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
)

func main() {
	tk, err := models.DecodeUnverifiedToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSdW5JZCI6IjA2ZmIxYmUyLWNhNWYtNDM0Ni04NDNmLWYzNmRhMTc5NGJjNCIsIkhvc3QiOiIxMjcuMC4wLjEiLCJUY3BQb3J0IjozNDMyMCwiS2NwUG9ydCI6MzQzMjAsIlRsc1BvcnQiOjM0MzIxLCJHcnBjUG9ydCI6MzQzMjIsIlAyUEFwaVBvcnQiOjM\n0MzIxLCJQZXJtaXNzaW9uIjoyLCJleHAiOjIwMTU4ODE3ODg4NSwibmJmIjoxNTg4MTUwMDg1fQ.VouW4wwODqqg2pCmghARfWPHrUtEPLmhkT_8hhCdEHc")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tk.Host)
	fmt.Println(tk.GrpcPort)
}
