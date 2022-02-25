// Build and Use this File to interact with the shodan package
// In this directory lab/3/shodan/main:
// go build main.go
// SHODAN_API_KEY=YOURAPIKEYHERE ./main <search term>

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"shodan/shodan"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: main {search|alert}")
	}

	apiKey := os.Getenv("SHODAN_API_KEY")
	s := shodan.New(apiKey)
	info, err := s.APIInfo()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf(
		"Query Credits: %d\nScan Credits:  %d\n\n",
		info.QueryCredits,
		info.ScanCredits)

	if os.Args[1] == "search" {
		search(s)
	} else if os.Args[1] == "alert" {
		alert(s)
	} else {
		log.Fatalln("Usage: main {search|alert}")
	}
}

func alert(c *shodan.Client) {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: main alert {new|list|triggers|add}")
	}

	// handling passed off to the appropriate handler
	// to avoid overly complex switch statement
	switch os.Args[2] {
	case "new":
		alertHandler(c)
	case "list":
		listHandler(c)
	case "triggers":
		triggersHandler(c)
	case "add":
		addTriggerHandler(c)
	default:
		log.Fatalln("Usage: main alert {new|list|triggers|add}")
	}
}

func search(c *shodan.Client) {
	if len(os.Args) != 3 {
		log.Fatalln("Usage: main search <searchterm>")
	}

	hostSearch, err := c.HostSearch(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("Host Data Dump\n")
	for _, host := range hostSearch.Matches {
		fmt.Println("==== start ", host.IPString, "====")
		h, _ := json.Marshal(host)
		fmt.Println(string(h))
		fmt.Println("==== end ", host.IPString, "====")
	}

	fmt.Printf("IP, Port\n")

	for _, host := range hostSearch.Matches {
		fmt.Printf("%s, %d\n", host.IPString, host.Port)
	}
}

func alertHandler(c *shodan.Client) {
	// validate argument
	if len(os.Args) != 6 {
		log.Fatalln("Usage: main alert new <name> <ip> <lifetime>")
	}

	name := os.Args[3]
	ip := os.Args[4]
	lifetime, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Fatalln("lifetime must be an integer")
	}

	// invoke shodan API
	res, err := c.CreateAlert(name, shodan.AlertFilter{IPs: []string{ip}}, lifetime)
	if err != nil {
		log.Fatalln(err)
	}

	// print results
	log.Printf("Alert created with name: %s\n", res.Id)
}

func listHandler(c *shodan.Client) {
	// invoke shodan API
	res, err := c.GetMyAlerts()
	if err != nil {
		log.Fatalln(err)
	}

	// print results
	if len(res) == 0 {
		fmt.Println("You have no alerts. Add one with `alert new`")
	} else {
		for _, alert := range res {
			fmt.Printf("%s\n", alert.Name)
		}
	}
}

func triggersHandler(c *shodan.Client) {
	// invoke shodan API
	res, err := c.GetAllValidTriggers()
	if err != nil {
		log.Fatalln(err)
	}

	// print results
	for _, trigger := range res {
		fmt.Printf("%s: (%s) -> %s\n", trigger.Name, trigger.Rule, trigger.Description)
	}
}

func addTriggerHandler(c *shodan.Client) {
	// validate argument
	if len(os.Args) != 5 {
		log.Fatalln("Usage: main alert add <alert name> <trigger name>")
	}

	triggerName := os.Args[3]
	alertName := os.Args[4]

	// invoke shodan API
	err := c.AddTriggerToAlert(triggerName, alertName)
	if err != nil {
		log.Panicln(err)
	}

	// print results
	log.Println("Successfuly added trigger to event")
}
