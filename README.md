# FetchTest
This is my attempt at the Fetch Backend Apprentice coding challenge. The challenge was written in Go on Windows 10 using Visual Studio Code as an editor, Git Bash as a terminal, and cURL to test functionality in the terminal. 

Docker was attempted to be used for containering to ensure the functionality on all OSes, but an issue on my end disallowed me from building it. Therefore, the included Dockerfile is an untested file from a developer who is new to using Docker, and most likely nonfunctional. 

# Process

Due to not being able to get Docker working, I can only ensure this process will work on a Windows 10 computer.
1. Clone repository to an accessible location
2. Open terminal of choice at repository location(I've personally tested Bash and Windows' command prompt)
3. run command "go run .". 
4. With the program running on that terminal, open a new terminal at the location that contains the test files to be fed to /receipts/process (if you're using my test files, that would be within the repository location).
5. Run the following commands to test each endpoint

# Commands
These commands use cURL for testing. cURL should be installed by default on modern Windows 10 and Mac computers, but the same can't be said for Linux based systems. If for whatever reason cURL is not installed, the download can be found at https://curl.se/

*Process Receipts*-curl http://localhost:8080/receipts/process --include --header "Content-Type: application/json" --request "POST" --data @testRec1.json (for testing different jsons, replace testRec1.json with whatever other json you wish to test with)

*Get Points*-curl http://localhost:8080/receipts/:id/points (replace :id with a valid id for testing. A processed receipt with id of 1 is provided by default for testing purposes)

*Get Processed Receipts*-curl http://localhost:8080/receipts/processed (A function I included for testing purposes. Returns a json with all processed receipt information within. Useful for retrieving ids to test Get Points with)



