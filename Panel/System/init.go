package System

import "fmt"

func init() {
	var err = error(nil)
	Data, err = GetOSData()
	temp, err := getPublicIP()
	if err != nil {
		fmt.Println("GetPublicIP() Error: ", err.Error())
	}
	PublicIP = temp
}
