// main.go

package main

func main() {
	a := App{}
	a.Initialize("postgres", "spoonlamp", "go-comments-api")

	a.Run(":8000")
}
