// bhg-scanner/scanner.go modified from Black Hat Go > CH2 > tcp-scanner-final > main.go
// Code : https://github.com/blackhat-go/bhg/blob/c27347f6f9019c8911547d6fc912aa1171e6c362/ch-2/tcp-scanner-final/main.go
// License: {$RepoRoot}/materials/BHG-LICENSE

package scanner

import (
	"fmt"
	"net"
	"os"
	"time"
	"log"
)

// Server contains a map that lists which ports are open for a given server.
// Ports contains this mapping
// Size is the number of ports scanned
// Address is the server url
type Server struct {
	Ports map[int]bool
	Size int
	Address string
}

func worker(addr string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p)    
		conn, err := net.DialTimeout("tcp", address, 500 * time.Millisecond)
		if err != nil { 
			results <- -p
			continue
		}
		conn.Close()
		results <- p
	}
}

// PortScanner takes a url to scan and returns the count of open and closed ports.
// it also takes a filename and if the filename is not "", outputs results to a CSV with that name
func PortScanner(address, fname string) (openCount, closeCount int) {  
	if fname == "" {
		return portScannerWrapped(address, fname, false)
	} else {
		return portScannerWrapped(address, fname, true)
	}
}

// portScannerWrapped takes arguments handed to the public function PortScanner and then
// does the actual output, switching based on the value of shouldOutputCSV to determine
// if the function should write to a file
func portScannerWrapped(address, fname string, shouldOutputCSV bool) (openCount, closeCount int) {
	ports := make(chan int, 100) 
	results := make(chan int)

	scannedServer := Server{
		Ports: make(map[int]bool),
		Size: 1024,
		Address: address,
	}

	for i := 0; i < cap(ports); i++ {
		go worker(address, ports, results)
	}

	go func() {
		for i := 1; i <= scannedServer.Size; i++ {
			ports <- i
		}
	}()

	for i := 0; i < scannedServer.Size; i++ {
		port := <-results
		if port < 0 {
			scannedServer.Ports[port] = false
			closeCount++
		} else {
			scannedServer.Ports[port] = true
			openCount++
		}
	}

	close(ports)
	close(results)

	if !shouldOutputCSV {
		return
	}

	f, err := os.Create(fname)
	if err != nil {
		log.Println("Error creating output file")
		return
	}
	defer f.Close()

	_, err = f.WriteString("port,isOpen\n")
	if err != nil {
		panic(err)
	}

	for i := 1; i <= scannedServer.Size; i++ {
		if scannedServer.Ports[i] {
			_, err = f.WriteString(fmt.Sprintf("%v,%v\n", i, true))
		} else {
			_, err = f.WriteString(fmt.Sprintf("%v,%v\n", i, false))
		}
		if err != nil {
			panic(err)
		}
	}

	return
}
