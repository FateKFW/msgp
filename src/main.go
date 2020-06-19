package main

import "msgp/comn"

func main() {
	gw := comn.BindParam()
	gw.Start()
}

