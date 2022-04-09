# Lab 7
Due April 9th at 11:59PM

## Development Work [18 points]
- [3pt] Complete the TODOs in [main.go](course-materials/materials/lab/7/main/main.go)
- [12pt] Complete the TODOs in [hscan.go](course-materials/materials/lab/7/hscan/hscan.go)
- [3pt] Create at least one new test in [hscan_test.go](course-materials/materials/lab/7/hscan/hscan_test.go)

## Capture  details [2pts]
- Capture Timing Details (per hscan.go) for various implementation of creating the hash maps

## Submit 
1. Link to your Github Repository [16pts]
2. Report the numbers [2pts]
2. List of Collaborator

## Timing
No goroutines: 3.91s user 0.38s system 102% cpu 4.188 total
10 workers: 8.10s user 0.55s system 319% cpu 2.707 total
100 workers: 7.61s user 0.51s system 329% cpu 2.466 total
400 workers: 7.65s user 0.50s system 342% cpu 2.379 total
1000 workers: 7.88s user 0.52s system 339% cpu 2.476 total

## Time Per Password
No groutines: 0.00000251160146378 seconds/password
400 workers: 0.00000142671916961 seconds/password