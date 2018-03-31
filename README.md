# fumble
Fumble is a popular app that lets friends know when they cross paths. The mobile client sends the users long,lat to a post endpoint once a second, all day long. The user should be able to open the app at any point and see the people they've crossed paths with. 

# API:
fumble will run on http://localhost:3000 and will support the following REST APIs:
1. GET /fumble/
    returns "Fumble, a cross path service for friends"

2. POST /fumble/location with data userId={userId}/long={longitude}/lat={latitude}
    will update the coordinate of users

    example: 


    ```
    $ curl -d '{"userId": 1, "long": 70, "lat": -35}' -X POST -H "application/json" http://localhost:3000/fumble/location
    http 200

    $ curl -d '{"userId": 2, "long": 70, "lat": -35}' -X POST -H "application/json" http://localhost:3000/fumble/location
    http 200
    
    ```

3.  GET /fumble/friends
    returns a list of all users

    ```
    $ curl http://localhost:3000/fumble/friends
    {
    All users: [
    1,
    2
    ]
    }

    ```

4.  GET /fumble/friends/{userId}
    returns a list of users who cross paths with {userId}

    ```
    $ curl http://localhost:3000/fumble/friends/1
    {              
    friends: [    
    2             
    ]              
    }
                  
    ```
 
the server should respond with 404 to all other requests not listed above
 
 # environment & build
 require Go and mySql
 dbinit.sql should be run after the installation of mysql to set up the db environment
 
$ go build ../src/github.com/gengwensu/fumble/fumble.go ../src/github.com/gengwensu/fumble/queryDB.go

$./fumble &

...
 
# Implementation
A cache with MAXCACHEENTRY entries are added to improve the performance. fumble will check the cache before querying the location table in the database.
