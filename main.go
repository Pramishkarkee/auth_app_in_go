package main

import (
	"fmt"
	"auth/test/conn"
	"auth/test/router"
)



func main(){
	mydb.GetDatabase()
	router.InitializeRoute()
	fmt.Println("ok")
}