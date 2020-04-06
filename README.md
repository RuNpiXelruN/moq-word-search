 # gRPC Word Search Challenge
 This app enables the user to search and affect a set of pre-defined search terms as below.

 ## Useful Commands
 ```
 make help
    - prints out endpoints, and lists all available make commands    

 make proto
    - Generates proto api, grpc-gateway definition, and swagger definition    

 make build
    - Runs test suite, and builds binaries for osx, linux, and windows, and places them in the ./bin folder.    
 
 make buildrun
    - Runs tests, builds binaries and runs osx version of binary    
 
 make run
    - Runs osx version of binary    
 
 make generate
    - Generates proto output, runs tests, builds binaries, and runs osx version of binary    
```
 ### Single word search
 ```
 curl "http://localhost:8090/api/words?word=sawyer"
 ```
 
 ### Update word search list
 ```
 curl -X "POST" "http://localhost:8090/api/words/sawyer"
 ```
 
 ### Fetch top 5 most searched for words
 ```
 curl "http://localhost:8090/api/words/popular"
 ```