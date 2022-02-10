package scanner

import (
	"os"
	"testing"
)

// THESE TESTS ARE LIKELY TO FAIL IF YOU DO NOT CHANGE HOW the worker connects (e.g., you should use DialTimeout)
func TestOpenPort(t *testing.T) {

	open, _ := PortScanner("scanme.nmap.org", "") // Currently function returns only number of open ports
	want := 2                                     // default value when passing in 1024 TO scanme; also only works because currently PortScanner only returns
	//consider what would happen if you parameterize the portscanner address and ports to scan

	if open != want {
		t.Errorf("got %d, wanted %d", open, want)
	}
}

func TestTotalPortsScanned(t *testing.T) {

	open, closed := PortScanner("scanme.nmap.org", "") // Currently function returns only number of open ports
	want := 1024                                       // default value; consider what would happen if you parameterize the portscanner ports to scan
	got := open + closed

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

// I, Andey Robins, being the owner and administrator of all servers located
// at andeyrobins.org and its subdomains, hereby grant you permission to use
// this test to verify my grade for this assignment and any other course related
// tasks. This does not extend to the general public and is not intended to be
// used for any other purpose.
func TestAndeyRobinsOrg(t *testing.T) {
	open, _ := PortScanner("andeyrobins.org", "")
	want := 3

	if open != want {
		t.Errorf("got %d, wanted %d", open, want)
	}
}

func TestOutput(t *testing.T) {
	PortScanner("scanme.nmap.org", "out.csv")

	if _, err := os.Stat("out.csv"); err != nil {
		t.Errorf("file not found")
	}
}
