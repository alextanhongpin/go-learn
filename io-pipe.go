package main

import (
	"io"
	"log"
	"os"
)

func main() {
	pr, pw := io.Pipe()
	defer pw.Close()

	// cmd := exec.Command("cat", "out.json")
	// cmd.Stdout = pw
	go func() {
		defer pr.Close()
		if _, err := io.Copy(os.Stdout, pr); err != nil {
			log.Fatal(err)
		}
	}()
	n, err := pw.Write([]byte("hello world, this is amazing"))
	if err != nil {
		log.Println(err)
	}
	log.Println("n", n)

	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }

}
