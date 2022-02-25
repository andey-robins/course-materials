# Lab 3 - PART 2 
Due Feb 25th at 11:59PM


Lab3Dir : ~/course_materials/materials/lab/3

Lab3Dir has two main folders main and shodan
main folder, where you'll run the program from (read top of main.go for details); you'll need to modify this based on which options you choose below
shodan folder, where you'll be adding additional code, do not try to run this code independently

## Extension

For my extension I added in the functionality to add alerts to a specified ip. You are then able to query a list of triggers and add a trigger to the alert IP.

To use the old search functionality use `./main search <query>`. The new query is accessed with `./main alert {alert queries}`.

There are four supported operands to alert: `new, list, triggers, and add`. These create a new alert, list all active alerts, list all valid triggers, and add a trigger to an alert respectively. Each subcommand provides information on usage when being invoked.

```bash
./main alert new andeyrobins.org 45.33.22.9
./main alert list 
./main alert triggers
./main alert add <trigger name> <alert id>
```

## Part 1(12 points)

Within {Lab3Dir}/shodan/shodan there are three files.

Two access different [Shodan API methods](https://developer.shodan.io/api)
 - host.go →Host/Search
 - api.go → APIInfo

The third file `shodan.go` contains a helper function to create a new client (with the API you provide on the command line; see README's for more details on usage)

### Option 1: 
Create another method file (like host.go/api.go) to access another method.
### Option 2: 
Extend host.go to build up more complex queries
### Option 3: 
Create a helper file to pull currently available facets and filters

### Required:
Create/Update the README in the shodan directory to discuss which option you selected and what files you created/modified.

## Part 2 (8 points)
Extend {Lab3Dir}/main/main.go to use your new functionality

### Required

 Create/Update the README in the main directory to discuss HOW to use your program (provide example command-line usages)

 # Submission

 ## Option 1 Preferred*
* Update / Push Your Code to your github repo and 
* then submit Link to your Public Repo on WyoCourses

## Option 2: 
* Download Lab 3 
  * find the lab/3 folder in the navigation explorer, 
  * right click, download
* Rename it to YOURLASTNAME_3.tar 
* upload it to WyoCourses