# SlicingDice Official Go Client (v1.0)
### Build Status: [![CircleCI](https://circleci.com/gh/SlicingDice/slicingdice-go.svg?style=svg)](https://circleci.com/gh/SlicingDice/slicingdice-go)

Official Go client for [SlicingDice](http://www.slicingdice.com/), Data Warehouse and Analytics Database as a Service.  

[SlicingDice](http://www.slicingdice.com/) is a serverless, API-based, easy-to-use and really cost-effective alternative to Amazon Redshift and Google BigQuery.

## Documentation

If you are new to SlicingDice, check our [quickstart guide](http://panel.slicingdice.com/docs/#quickstart-guide) and learn to use it in 15 minutes.

Please refer to the [SlicingDice official documentation](http://panel.slicingdice.com/docs/) for more information on [analytics databases](http://panel.slicingdice.com/docs/#analytics-concepts), [data modeling](http://panel.slicingdice.com/docs/#data-modeling), [data insertion](http://panel.slicingdice.com/docs/#data-insertion), [querying](http://panel.slicingdice.com/docs/#data-querying), [limitations](http://panel.slicingdice.com/docs/#current-slicingdice-limitations) and [API details](http://panel.slicingdice.com/docs/#api-details).

## Tests and Examples

Whether you want to test the client installation or simply check more examples on how the client works, take a look at [tests and examples directory](tests_and_examples/).

## Installing

In order to install the Go client, you only need to execute the [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies) command.

```bash
go get github.com/SlicingDice/slicingdice-go/slicingdice
```

## Usage

The following code snippet is an example of how to add and query data
using the SlicingDice GO client. We entry data informing
'user1@slicingdice.com' has age 22 and then query the database for
the number of entities with age between 20 and 40 years old.
If this is the first record ever entered into the system,
 the answer should be 1.

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    // Configure client
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "YOUR_API_KEY"
    client := slicingdice.New(keys, 60)
    // If you need production end-point you can remove this
    client.Test = true

    // inserting data
    insert_data := map[string]interface{}{
        "user1@slicingdice.com": map[string]int{
            "age": 22,
        },
        "auto-create": []string{"table", "column"},
    }
    client.Insert(insert_data)

    // Querying data
    query_data := map[string]interface{}{
            "query-name": "users-between-20-and-40",
            "query": []map[string]interface{} {
                map[string]interface{}{
                    "age": map[string][]int{
                        "range": []int{20, 40},
                    },
                },
            },
    }

    result, _ := client.CountEntity(query_data)
    fmt.Println(result["status"])
}
```

## Reference

`SlicingDice` encapsulates logic for sending requests to the API. Its methods are thin layers around the [API endpoints](http://panel.slicingdice.com/docs/#api-details-api-endpoints), so their parameters and return values are JSON-like `interface{}` objects with the same syntax as the [API endpoints](http://panel.slicingdice.com/docs/#api-details-api-endpoints)

### Attributes

* `key (map[string]string)` - [API key](http://panel.slicingdice.com/docs/#api-details-api-connection-api-keys) to authenticate requests with the SlicingDice API.
* `timeout (int)` - Amount of time, in seconds, to wait for results for each request.

### Constructors

`New(key *APIKey, timeout int) *SlicingDice`
* `key (APIKey)` - [API key](http://panel.slicingdice.com/docs/#api-details-api-connection-api-keys) to authenticate requests with the SlicingDice API.
* `timeout (int)` - Amount of time, in seconds, to wait for results for each request.

### `GetDatabase()`
Get information about the current SlicingDice database. This method corresponds to a [GET request at /database](http://panel.slicingdice.com/docs/#api-details-api-endpoints-get-database). **IMPORTANT:** You can't make this request on tests end-point

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)
    fmt.Println(client.GetDatabase())
}
```

#### Output example

```json
{
    "name": "Database 1",
    "description": "My first database",
    "data-expiration": 30,
    "created-at": "2016-04-05T10:20:30Z"
}
```

### `GetColumns()`
Get all created columns, both active and inactive ones. This method corresponds to a [GET request at /column](http://panel.slicingdice.com/docs/#api-details-api-endpoints-get-column).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)
    // If you need production end-point you can remove this
    client.Test = true
    fmt.Println(client.GetColumns())
}
```

#### Output example

```json
{
    "active": [
        {
          "name": "Model",
          "api-name": "car-model",
          "description": "Car models from dealerships",
          "type": "string",
          "category": "general",
          "cardinality": "high",
          "storage": "latest-value"
        }
    ],
    "inactive": [
        {
          "name": "Year",
          "api-name": "car-year",
          "description": "Year of manufacture",
          "type": "integer",
          "category": "general",
          "storage": "latest-value"
        }
    ]
}
```

### `CreateColumn(query interface{})`
Create a new column. This method corresponds to a [POST request at /column](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-column).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    columnData := map[string]interface{}{
        "name":        "Year",
        "type":        "integer",
        "description": "Year of manufacturing",
        "storage":     "latest-value",
    }

    fmt.Println(client.CreateColumn(columnData))
}
```

#### Output example

```json
{
    "status": "success",
    "api-name": "year"
}
```

### `Insert(query interface{})`
Insert data to existing entities or create new entities, if necessary. This method corresponds to a [POST request at /insert](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-insert).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    insertData := map[string]interface{}{
        "user1@slicingdice.com": map[string]interface{}{
            "car-model": "Ford Ka",
            "year":      2016,
        },
        "user2@slicingdice.com": map[string]interface{}{
            "car-model": "Honda Fit",
            "year":      2016,
        },
        "user3@slicingdice.com": map[string]interface{}{
            "car-model": "Toyota Corolla",
            "year":      2010,
            "test-drives": []map[string]string{
                map[string]string{
                    "value": "NY",
                    "date":  "2016-08-17T13:23:47+00:00",
                },
                map[string]string{
                    "value": "NY",
                    "date":  "2016-08-17T13:23:47+00:00",
                },
                map[string]string{
                    "value": "CA",
                    "date":  "2016-04-05T10:20:30Z",
                },
            },
        },
        "user4@slicingdice.com": map[string]interface{}{
            "car-model": "Ford Ka",
            "year":      2005,
            "test-drives": []map[string]string{
                map[string]string{
                    "value": "NY",
                    "date":  "2016-08-17T13:23:47+00:00",
                },
            },
        },
        "auto-create": []string{"table", "column"},
    }
    fmt.Println(client.Insert(insertData))
}
```

#### Output example

```json
{
    "status": "success",
    "inserted-entities": 4,
    "inserted-columns": 10,
    "took": 0.023
}
```

### `ExistsEntity(ids)`
Verify which entities exist in a database given a list of entity IDs. This method corresponds to a [POST request at /query/exists/entity](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-exists-entity).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    entities := []string{
        "user1@slicingdice.com",
        "user2@slicingdice.com",
        "user3@slicingdice.com",
        "otheruser@slicingdice.com",
    }
    fmt.Println(client.ExistsEntity(entities))
}
```

#### Output example

```json
{
    "status": "success",
    "exists": [
        "user1@slicingdice.com",
        "user2@slicingdice.com",
        "user3@slicingdice.com"
    ],
    "not-exists": [
        "otheruser@slicingdice.com"
    ],
    "took": 0.103
}
```

### `CountEntityTotal()`
Count the number of inserted entities. This method corresponds to a [GET request at /query/count/entity/total](http://panel.slicingdice.com/docs/#api-details-api-endpoints-get-query-count-entity-total).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    fmt.Println(client.CountEntityTotal())
}
```

#### Output example

```json
{
    "status": "success",
    "result": {
        "total": 42
    },
    "took": 0.103
}
```

### `CountEntity(query interface{})`
Count the number of entities matching the given query. This method corresponds to a [POST request at /query/count/entity](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-count-entity).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := []interface{}{
        map[string]interface{}{
            "query-name": "corolla-or-fit",
            "query": []interface{}{
                map[string]interface{}{
                    "car-model": map[string]string{
                        "equals": "toyota corolla",
                    },
                },
                "or",
                map[string]interface{}{
                    "car-model": map[string]string{
                        "equals": "honda fit",
                    },
                },
            },
            "bypass-cache": false,
        },
        map[string]interface{}{
            "query-name": "ford-ka",
            "query": []map[string]interface{}{
                map[string]interface{}{
                    "car-model": map[string]string{
                        "equals": "ford ka",
                    },
                },
            },
            "bypass-cache": false,
        },
    }

    fmt.Println(client.CountEntity(query))
}
```

#### Output example

```json
{
   "result":{
      "corolla-or-fit":2,
      "ford-ka":2
   },
   "status":"success",
   "took":0.053
}
```

### `CountEvent(query interface{})`
Count the number of occurrences for time-series events matching the given query. This method corresponds to a [POST request at /query/count/event](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-count-event).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := []interface{}{
        map[string]interface{}{
            "query-name": "test-drives-in-ny",
            "query": []map[string]interface{}{
                map[string]interface{}{
                    "test-drives": map[string]interface{}{
                        "equals": "NY",
                        "between": []string{
                            "2016-08-16T00:00:00Z",
                            "2016-08-18T00:00:00Z",
                        },
                    },
                },
            },
            "bypass-cache": true,
        },
        map[string]interface{}{
            "query-name": "test-drives-in-ca",
            "query": []map[string]interface{}{
                map[string]interface{}{
                    "test-drives": map[string]interface{}{
                        "equals": "CA",
                        "between": []string{
                            "2016-04-04T00:00:00Z",
                            "2016-04-06T00:00:00Z",
                        },
                    },
                },
            },
            "bypass-cache": true,
        },
    }
    fmt.Println(client.CountEvent(query))
}
```

#### Output example

```json
{  
   "result":{  
      "test-drives-in-ca":1,
      "test-drives-in-ny":3
   },
   "status":"success",
   "took":0.029
}
```

### `TopValues(query interface{})`
Return the top values for entities matching the given query. This method corresponds to a [POST request at /query/top_values](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-top-values).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "car-year": map[string]interface{}{
            "year": 2,
        },
        "car models": map[string]interface{}{
            "car-model": 3,
        },
    }

    fmt.Println(client.TopValues(query))
}
```

#### Output example

```json
{
   "result":{
      "car models":{
         "car-model":[
            {
               "quantity":2,
               "value":"ford ka"
            },
            {
               "quantity":1,
               "value":"honda fit"
            },
            {
               "quantity":1,
               "value":"toyota corolla"
            }
         ]
      },
      "car-year":{
         "year":[
            {
               "quantity":2,
               "value":"2016"
            },
            {
               "quantity":1,
               "value":"2005"
            }
         ]
      }
   },
   "status":"success",
   "took":0.026
}
```

### `Aggregation(query interface{})`
Return the aggregation of all columns in the given query. This method corresponds to a [POST request at /query/aggregation](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-aggregation).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "query": []map[string]interface{}{
            map[string]interface{}{
                "year": 2,
            },
            map[string]interface{}{
                "car-model": 2,
                "equals": []string{
                    "honda fit",
                    "toyota corolla",
                },
            },
        },
    }

    fmt.Println(client.Aggregation(query))
}
```

#### Output example

```json
{
   "year":[
      {
         "car-model":[
            {
               "quantity":1,
               "value":"honda fit"
            }
         ],
         "quantity":2,
         "value":"2016"
      },
      {
         "quantity":1,
         "value":"2005"
      }
   ]
}
```

### `GetSavedQueries()`
Get all saved queries. This method corresponds to a [GET request at /query/saved](http://panel.slicingdice.com/docs/#api-details-api-endpoints-get-query-saved).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    fmt.Println(client.GetSavedQueries())
}
```

#### Output example

```json
{  
   "saved-queries":[  
      {  
         "cache-period":100,
         "name":"my-saved-query",
         "query":[  
            {  
               "car-model":{  
                  "equals":"honda fit"
               }
            },
            "or",
            {  
               "car-model":{  
                  "equals":"toyota corolla"
               }
            }
         ],
         "type":"count/entity"
      }
   ],
   "status":"success",
   "took":0.011
}
```

### `CreateSavedQuery(query interface{})`
Create a saved query at SlicingDice. This method corresponds to a [POST request at /query/saved](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-query-saved).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "name": "my-saved-query",
        "type": "count/entity",
        "query": []interface{}{
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "honda fit",
                },
            },
            "or",
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "toyota corolla",
                },
            },
        },
        "cache-period": 100,
    }

    fmt.Println(client.CreateSavedQuery(query))
}
```

#### Output example

```json
{  
   "cache-period":100,
   "name":"my-saved-query",
   "query":[  
      {  
         "car-model":{  
            "equals":"honda fit"
         }
      },
      "or",
      {  
         "car-model":{  
            "equals":"toyota corolla"
         }
      }
   ],
   "type":"count/entity"
}
```

### `UpdateSavedQuery(string queryName, query interface{})`
Update an existing saved query at SlicingDice. This method corresponds to a [PUT request at /query/saved/QUERY_NAME](http://panel.slicingdice.com/docs/#api-details-api-endpoints-put-query-saved-query-name).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)


func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "type": "count/entity",
        "query": []interface{}{
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "honda fit",
                },
            },
            "or",
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "toyota corolla",
                },
            },
        },
        "cache-period": 100,
    }

    fmt.Println(client.UpdateSavedQuery("my-saved-query", query))
}
```

#### Output example

```json
{
   "cache-period":100,
   "query":[
      {
         "car-model":{
            "equals":"honda fit"
         }
      },
      "or",
      {
         "car-model":{
            "equals":"toyota corolla"
         }
      }
   ],
   "status":"success",
   "took":0.011,
   "type":"count/entity"
}
```

### `GetSavedQuery(string queryName)`
Executed a saved query at SlicingDice. This method corresponds to a [GET request at /query/saved/QUERY_NAME](http://panel.slicingdice.com/docs/#api-details-api-endpoints-get-query-saved-query-name).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    fmt.Println(client.GetSavedQuery("my-saved-query"))
}
```

#### Output example

```json
{
   "query":[
      {
         "car-model":{
            "equals":"honda fit"
         }
      },
      "or",
      {
         "car-model":{
            "equals":"toyota corolla"
         }
      }
   ],
   "result":{
      "query":2
   },
   "status":"success",
   "took":0.021,
   "type":"count/entity"
}
```

### `DeleteSavedQuery(string queryName)`
Delete a saved query at SlicingDice. This method corresponds to a [DELETE request at /query/saved/QUERY_NAME](http://panel.slicingdice.com/docs/#api-details-api-endpoints-delete-query-saved-query-name).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    fmt.Println(client.DeleteSavedQuery("my-saved-query"))
}
```

#### Output example

```json
{
   "cache-period":100,
   "deleted-query":"my-saved-query",
   "query":[
      {
         "car-model":{
            "equals":"honda fit"
         }
      },
      "or",
      {
         "car-model":{
            "equals":"toyota corolla"
         }
      }
   ],
   "status":"success",
   "took":0.011,
   "type":"count/entity"
}
```

### `Result(query interface{})`
Retrieve inserted values for entities matching the given query. This method corresponds to a [POST request at /data_extraction/result](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-data-extraction-result).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "query": []interface{}{
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "ford ka",
                },
            },
            "or",
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "honda fit",
                },
            },
        },
        "columns":     []string{"car-model", "year"},
        "limit":      2,
    }

    fmt.Println(client.Result(query))
}
```

#### Output example

```json
{
   "data":{
      "user2@slicingdice.com":{
         "car-model":"honda fit",
         "year":"2016"
      },
      "user4@slicingdice.com":{
         "car-model":"ford ka",
         "year":"2005"
      }
   },
   "next-page":null,
   "page":1,
   "status":"success",
   "took":0.055
}
```

### `Score(query interface{})`
Retrieve inserted values as well as their relevance for entities matching the given query. This method corresponds to a [POST request at /data_extraction/score](http://panel.slicingdice.com/docs/#api-details-api-endpoints-post-data-extraction-score).

#### Request example

```go
package main

import (
    "fmt"
    "github.com/SlicingDice/slicingdice-go/slicingdice"
)

func main() {
    keys := new(slicingdice.APIKey)
    keys.MasterKey = "MASTER_OR_READ_API_KEY"
    client := slicingdice.New(keys, 60)

    // If you need production end-point you can remove this
    client.Test = true

    query := map[string]interface{}{
        "query": []interface{}{
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "toyota corolla",
                },
            },
            "or",
            map[string]interface{}{
                "car-model": map[string]string{
                    "equals": "honda fit",
                },
            },
        },
        "columns":     []string{"car-model", "year"},
        "limit":      2,
    }

    fmt.Println(client.Score(query))
}
```

#### Output example

```json
{
   "data":{
      "user2@slicingdice.com":{
         "car-model":"honda fit",
         "score":1,
         "year":"2016"
      },
      "user3@slicingdice.com":{
         "car-model":"toyota corolla",
         "score":1,
         "year":"2010"
      }
   },
   "next-page":null,
   "page":1,
   "status":"success",
   "took":0.036
}
```

## License

[MIT](https://opensource.org/licenses/MIT)
