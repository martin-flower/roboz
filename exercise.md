## platform

Platform consists of swarm of micro services running as Docker containers. Primary development platforms are .NET Core and Node JS in conjunction with other platforms. Backed mostly using PostgreSQL as relational/document storage and Amazon S3 as blob storage.

## task

Create a new micro service that could fit into the Platform environment as described above. The service will simulate a robot moving in an office space and cleaning the places the robot visits. The path of the robot movement is described by start coordinates and move commands. After the cleaning has been done the robot reports the number of unique places cleaned, the service will store the result into the database and returns the created record in JSON format. The service listens to HTTP protocol on port 5000.

* Request method `POST`
* Request path: `/developer-test/enter-path` 

## input criteria

```0 ≤ number of commmands elements ≤ 10000
−100 000 ≤ x ≤ 100 000, x ∈ Z
−100 000 ≤ y ≤ 100 000, y ∈ Z
direction ∈ {north, east, south, west}
0 < steps < 100000, steps ∈ Z
```

## request body example
```
{
  "start": {
    "x": 10,
    "y": 22 
  },
  "commmands": [
    {
      "direction": "east",
      "steps": 2 
    },
    {
      "direction": "north",
      "steps": 1
    } 
  ]
}
```

The resulting value will be stored in a table named `executions` together with timestamp of insertion, number of command elements and duration of the calculation in seconds.

## stored record example

| id | timestamp | commmands | result | duration |
| -- | -- | -- | -- | -- |
| 1234 | 2018-05-12 12:45:10.851596 +02:00 | 2 | 4 | 0.000123 |

## notes
* You can assume, for the sake of simplicity, that the office can be viewed as a grid where the robot moves only on the vertices.
* The robot cleans at every vertex it touches, not just where it stops.
* All input should be considered well formed and syntactically correct. There is no need,
therefore, to implement elaborate input validation.
* The robot will never be sent outside the bounds of the office.
* Ensure that database connection is configurable using environment variable.
* Think about structure, readability, maintainability, performance, re-usability and test-ability of the code. Like the solution is going to be deployed into production environment. You should be proud of what you deliver.
* Use only open source dependencies if needed.
* Include Dockerfile and docker-compose configuration files in the solution.
