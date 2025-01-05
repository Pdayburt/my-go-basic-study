package main

func main() {

	serve := InitWebService()
	serve.Run(":8080")

}
