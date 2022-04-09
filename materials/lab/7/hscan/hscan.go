package hscan

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

//==========================================================================\\

var shalookup map[string]string
var md5lookup map[string]string

func GuessSingle(sourceHash string, filename string) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		password := scanner.Text()

		switch len(sourceHash) {
		case 32:
			{ // md5 is 32 hex characters long
				hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
				if hash == sourceHash {
					fmt.Printf("[+] Password found (MD5): %s\n", password)
				}
			}
		case 64:
			{ // SHA-256 is 64 hex characters long
				hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
				if hash == sourceHash {
					fmt.Printf("[+] Password found (SHA-256): %s\n", password)
				}
			}
		default:
			fmt.Println("Unexpected hash length.")
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func GenHashMaps(filename string) {

	shalookup = make(map[string]string)
	md5lookup = make(map[string]string)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	fmt.Println("Creating channels")

	// create worker infrastructure
	workerCount := 400
	shaPass := make(chan string, workerCount)
	md5Pass := make(chan string, workerCount)
	shaHash := make(chan string, workerCount)
	md5Hash := make(chan string, workerCount)

	fmt.Println("Channels ready to go. Spawning workers.")

	// startup worker goroutines
	for i := 0; i < workerCount; i++ {
		go shaWorker(shaPass, shaHash)
		go md5Worker(md5Pass, md5Hash)
	}

	fmt.Println("Workers spawned. Spawning aggregators.")

	// sha aggregator
	go func(results chan string) {
		runningWorkers := workerCount
		for runningWorkers > 0 {
			res := <-results
			if res == "i have finished hashing" {
				runningWorkers--
			} else {
				splitStr := strings.Split(res, ",")
				hash, pass := splitStr[0], splitStr[1]
				shalookup[hash] = pass
			}
		}
	}(shaHash)

	// md5 aggregator
	go func(results chan string) {
		runningWorkers := workerCount
		for runningWorkers > 0 {
			res := <-results
			if res == "i have finished hashing" {
				runningWorkers--
			} else {
				splitStr := strings.Split(res, ",")
				hash, pass := splitStr[0], splitStr[1]
				md5lookup[hash] = pass
			}
		}
	}(md5Hash)

	// disperses work to the workers
	for scanner.Scan() {
		password := scanner.Text()

		shaPass <- password
		md5Pass <- password
	}

	fmt.Println("Finished writing passwords to workers")

	shaPass <- "time to stop hashing"
	md5Pass <- "time to stop hashing"

	//TODO create a test in hscan_test.go so that you can time the performance of your implementation
	//Test and record the time it takes to scan to generate these Maps
	// 1. With and without using go subroutines
	// 2. Compute the time per password (hint the number of passwords for each file is listed on the site...)
}

func shaWorker(in, out chan string) {
	for {
		password := <-in

		if password == "time to stop hashing" {
			out <- "i have finished hashing"
			return
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		out <- fmt.Sprintf("%s,%s", hash, password)
	}
}

func md5Worker(in, out chan string) {
	for {
		password := <-in

		if password == "time to stop hashing" {
			out <- "i have finished hashing"
			return
		}

		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		out <- fmt.Sprintf("%s,%s", hash, password)
	}
}

func GetSHA(hash string) (string, error) {
	password, ok := shalookup[hash]
	if ok {
		return password, nil
	} else {
		return "", errors.New("password does not exist")
	}
}

func GetMD5(hash string) (string, error) {
	password, ok := md5lookup[hash]
	if ok {
		return password, nil
	} else {
		return "", errors.New("password does not exist")
	}
}
