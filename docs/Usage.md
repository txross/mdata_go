
# Options
  - `-v` - Verbosity (up to -vvv)
  - `-V` - Version - show version
  - `-S` - Run Client as a Rest Server
  - `-p` - Port to run Rest Server on, default 8888

  When the client is run without a `-S` arg, it will default to the CLI implementation.
  
  Running with `-S` to initiate REST server

# CLI
Where `mdata` refers to the binary installed in `/usr/bin/mdata`

## List<br>
  - List all existing products
    `mdata list`

## Show
  - Show existing product
    `mdata show <gtin>`

## Create
  - Create a new product
    `mdata create <gtin>`

## Update
  - Requires attributes
  - Attributes supplied to update will overwrite existing attributes
  - Can provide any number of key:value pair of attributes. Keep appending with the -a flag.
  - Requires existing product
  `mdata update <gtin> -a "<key>:<value>" [-a "<key>:<value>" -a "<key>:<value>" ...]`

## Set
  - Set the state of an existing product to one of ACTIVE, INACTIVE, DISCONTINUED
  `mdata set <gtin> ["ACTIVE", "INACTIVE", "DISCONTINUED"]`

## Delete<br> 
  - Requires product to be in INACTIVE state
    `mdata delete <gtin>` 

# Rest Server
Run the exact same commands against a rest interface

## List
`curl -X GET http://localhost:8888/products`

## Show
`curl -X GET http://localhost:8888/products/<gtin>`

## Create
```
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"Gtin":"25825825825825", "Attributes": {"uom": "cases", "name": "chicken wings"}}' \
  http://localhost:8888/products
  ```

## Delete
`curl -X DELETE http://localhost:8888/products/<gtin>`

## Set 
```
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"Gtin":"25825825825825", "State": "ACTIVE"}' \
  http://localhost:8888/products/state/25825825825825
  ```
 
## Update
```
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"Gtin":"25825825825825", "Attributes": {"uom": "lbs", "name": "chicken wings"}}' \
  http://localhost:8888/products/attr/25825825825825
  ```