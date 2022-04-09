package main

import (
	"fmt"
	"hscan/hscan"
	"os"
)

func main() {

	//To test this with other password files youre going to have to hash
	var md5hash = "77f62e3524cd583d698d51fa24fdff4f"
	var sha256hash = "95a5e1547df73abdd4781b6c9e55f3377c15d08884b11738c2727dbd887d4ced"

	var drmike1 = "90f2c9c53f66540e67349e0ab83d8cd0"                                 // p@ssword
	var drmike2 = "1c8bfe8f801d79745c4631d09fff36c82aa37fc4cce4fc946683d7b336b63032" // letmein

	if len(os.Args) != 2 {
		fmt.Println("Usage main.go <wordlist file>")
	}

	file := os.Args[1]

	hscan.GuessSingle(md5hash, file)
	hscan.GuessSingle(sha256hash, file)
	hscan.GuessSingle(drmike1, file)
	hscan.GuessSingle(drmike2, file)
	hscan.GenHashMaps(file)
	fmt.Println(hscan.GetSHA(sha256hash))
	fmt.Println(hscan.GetMD5(md5hash))
	fmt.Println(hscan.GetSHA(drmike2))
	fmt.Println(hscan.GetMD5(drmike1))
}
