# DVT Service

This project is used to facilitate users to operate validator sharding, SSV validator registration, etc.

The `DVT Service` provides APIs such as keystore upload, SSV cluster management, and asynchronously operates
SSV network contracts;

For ease of use, we provide some terminal command tools in package `tools`:

1. **dvt_tools**: Provide operations such as uploading keystore files in batches, setting FeeRecipient, etc., relying on
   DVT services;

The program is divided into production end and consumer end.

* dvt-service-producer: Provides an API interface and is the production end of SSV contract operation tasks.
* dvt-service-consumer: It ensures a single thread to execute the SSV contract transaction task and is the consumer side
  of the task.

## Environment

* nodejs ^16.0.0
* golang ^1.20.0
* mysql ^8.0
* redis

## Server configuration information

* dvt-service-producer: Linux 2C4G 50G
* dvt-service-consumer: Linux 2C4G 50G

> Low data volume, can share the database with other projects

## Project Structure Description

```
├── abi                 contract abi file
├── awsKms              aws kms module
├── cliCmd              dvt service cli commond active
├── common              common utils, like logs, gorm nologger.
├── conf                configuration module
├── email               send notice email by service
├── docs                swagger docs
├── init                database table initialization file
├── makeShares          Online sharding tool for ssv, used to generate shareData, implemented using nodejs
├── models              some common structures
├── process             Asynchronous service consumer, used for single-threaded processing of SSV transactions
├── service             gin api module
├── store               database operation module
├── tools               Other auxiliary tools
├── main.go    
```

## Usage

First, use this file `./init/init.sql` to initialize the database tables.

### 1. Initialize configuration info

The program supports reading configuration information by configuration files.

```json
{
  "api": {
    "port": 3000,
    "close_swagger": true,
    "keystore_secret_key": "openssl rand -hex 32"
  },
  "db": {
    "mysql": {
      "host": "mysql",
      "port": 3306,
      "username": "root",
      "password": "123456",
      "db": "dvt"
    },
    "redis": {
      "addr": "redis:6379",
      "db": 0,
      "password": "123456"
    }
  },
  "ssv": {
    "ssv_contract_addr": "0xDD9BC35aE942eF0cFa76930954a156B3fF30a4E1",
    "ssv_token_contract_addr": "0x9D65fF81a3c488d585bBfb0Bfe3c7707c7917f54",
    "contract_creation_block": 17507487,
    "amount_token_ssv": "8000000000000000000",
    "approve_check_token_ssv": "400000000000000000000",
    "make_shares": "/root/makeShares",
    "operators": [
      {
        "id": 1,
        "operatorKey": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBeVpGNUR2M2UwSkEzT25TSGwyQmMKNGFxbmpUTWFrUXNZSkY5eE55M21CVTZSQld1d2xVd1dIelJGWUFvb0FlRER3NlYxL3hRQ0JFaWJwTGx1RVdLTgoxNmRpcU5EVmY5VEZndmZlM2NHc3pNcDZCUE04bWhBdkx0c01DcHlXeDZtTEczVm0zVVRNK3hRdUJwVFZsdHNNCkV6eUZEZzNWTlphOW9hZkswbkVYRHVidlBIbkJCdWhlUW5LZThoUkJnRUo0emIrV3dncjFrM3YyWmkwTEtWNUQKYWd3c2QxK25Lb1grVktjYmJFVFBEdGRPV1AvZlpXM3dBMGp3R1pSdkhwNS8xUjBmZy91N01BUk1KTkRWVFYxQwo0Vlh1eHJkbHZWQ2JiS1pnWUIzY1ROSEMzZkVldit0NFVEeFJuQzdUcUN0WFZSYnpZQ001WHVSeUFRa3BiYU0wCjlRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"
      },
      {
        "id": 2,
        "operatorKey": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNC9NR0pJeE1vQlZuVTd5NXJKazcKeWk1VXoyL1UzMm52a2FWL0N5YlFSRU9YRWxtc1ZFTHVVckFwZGFpOTJjU29iLy8xWm5GUUEvcjFPZlNZdjJoZwp1djl3RXUvS3BSeFR0RVg2alRhbEg3RklaRC9UZXdxL1FQZ3Noc1hMVFpuenZvazlaUkNQU0RQZG9VRElBcXFDCjhsL2NtMFFDcTBWbFZWSGNXbFNhREVubUliNFFIUTdsNkFONS9INEdOZDV3YzI3eHhNR2Q4bTRIMU54U1ozNTYKOWFHOXovWEdWUC9wTVA1b0ZLb1JqNjZpeWMwL0ZOR3k1UXNFRG9lMVRHTG42enlpYUxJTUw1Q0tscEp3SHFldwpsYmtiMWRwZGRBcnl2L0dCdEhsV3l6RXpEdmFNMmJ4bWRSbGp0SmJSdTJFYkR6d2RrdDIzcmt2aGF2Qk82RVFtCjd3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"
      },
      {
        "id": 3,
        "operatorKey": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBb0ltaWhneStMUTEreWpZU01EaTEKU0R2YlRVcUs3MjAwMk1IQis1QUhpbzMveFFZY0lKNXAwUUNWNk5WZE5VREVTSnRBcjl2VU1iMSt5a0RESXQybwppVkFUTDBKVGlkVVRzdThRTjdnM3NZd0Z2eVcyTHo2b01GKzIrbDV6bVV5SWxsaTFiUXd0NUFUQ0JzK2tzRUlUCjZ5RnlUc2FCQlRUYVZkaC81NTJGaTdaUVhkcUV4eXdmYlJZb1J2TE9wTnF5YU84V1JCZ0h1WGxxcmwwZTl4eGwKUWpmQTd1MU1hSVY5RlVmL0JJVVpoRWdsVDRYc3ZYbmt6Y1h6RTU2RVNZNnlxeDA4NjMwRDRqenpPT0NsdlVHYwplUTB2NHFscDlOeFZFcFhqelpDajZna2ZaZlpkV25wSFd5VUgwZHZ4TkVGUjNZSWF3eG00bUwxU3ViTGpSZTR1CndRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"
      },
      {
        "id": 4,
        "operatorKey": "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBdk5DUVlLMmNPay96RjU0Uk1HSkUKYmR4bXRqS21UT3BvTEw5Qk0wbDNPNFBpbFFFbnN5MW5nZVdHVXZBWnpISXlxR2o0K0tlamJzZGxnbG1YaTg0awp1bGtNVHQwSDVWbzdZQUZTbGNrczdsRnZyQVRvdmM3byt0Rkl6MWlodS9zWDVyMlBuR2N2dXhWUHc5aW1CMjExCkZYUU1NZmhha0tkL2dOektrWmNFYjRtUG5ROWQwNFJuZXVCSktFaGtSbUNNTDVhNUpYY2Z2cDNCdHRQU3VYZkQKUVFmU2JjaVl5VVp6NUtFVHdUZnBLS1BpYjlvZG8rQWhLbFBEUlFWWGZNbHdBa1RGUTh1NXMvZWFFYzJ0UUxkVApWQVdSV20ySmwyR1BtOUdNZDdtOWRRS1dDUWVTeTJPTU5RTUpQNWV4RmlJTGFZTjIyTnpmTHA2amlZdkxKR25DCk1RSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"
      }
    ]
  },
  "rpc_url": "****eth rpc url****",
  "kms": {
    "access_key_id": "access_key_id",
    "access_key_secret": "access_key_secret",
    "region": "region",
    "key_id": "key_id"
  },
  "email": {
    "account": "",
    "password": "",
    "smtp": "",
    "to": [
    ],
    "cc": [
    ]
  }
}
```

### 2. build and start `dvt-service-consumer`

Refer to file `Dockerfile-consumer` to build or run the image

### 3. build and start `dvt-service-product`

Refer to file `Dockerfile-product` to build or run the image

The api module has integrated swagger. If you start in localhost and turn `close_swagger` to `false`, you can
visit http://localhost:3000/swagger/index.html
to view the api document.

### 4. setFeeRecipient

This cmd tool can notify dvt service to modify the feeRecipient address.

```bash
./dvt_tools run
██████  ██    ██ ████████       ███████ ███████ ██████  ██    ██ ██  ██████ ███████ 
██   ██ ██    ██    ██          ██      ██      ██   ██ ██    ██ ██ ██      ██      
██   ██ ██    ██    ██    █████ ███████ █████   ██████  ██    ██ ██ ██      █████   
██   ██  ██  ██     ██               ██ ██      ██   ██  ██  ██  ██ ██      ██      
██████    ████      ██          ███████ ███████ ██   ██   ████   ██  ██████ ███████ 
                                                                                    

Please select an option: 
  > setFeeRecipient

 INFO  Selected option: setFeeRecipient
 INFO  DVT-Staking-Service api rawUrl is set to  http://localhost:3000/ssv/setFeeRecipientAddress
Input Fee Recipient Address: 0x0000000000000000000000000000000000000000

```

### 5. uploadKeystore

This cmd tool uploads the share information generated in the previous step to dvt service.

```bash
./dvt_tools run
██████  ██    ██ ████████       ███████ ███████ ██████  ██    ██ ██  ██████ ███████ 
██   ██ ██    ██    ██          ██      ██      ██   ██ ██    ██ ██ ██      ██      
██   ██ ██    ██    ██    █████ ███████ █████   ██████  ██    ██ ██ ██      █████   
██   ██  ██  ██     ██               ██ ██      ██   ██  ██  ██  ██ ██      ██      
██████    ████      ██          ███████ ███████ ██   ██   ████   ██  ██████ ███████ 
                                                                                    

Please select an option: 
  > uploadKeystore

 INFO  Selected option: uploadKeystore
 INFO  DVT-Staking-Service api rawUrl is set to  http://localhost:3000/ssv/upload
Input keystore folder path: /Github/DVT_Service/validator_keys
 SUCCESS  Check keystore folder success
 WARNING  Skip by not keystore file: deposit_data.json                                                                                                                 
 SUCCESS  Keystore 000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 packaged                                                      
Package keystore info... [2/2] ███████████████████████████████████████ 100% | 0s
 SUCCESS  Package keystore info success!
Enter keystore password: ********

```

Complete information upload according to program guidance.

## Package

The following commands can be used to package the tool for deployment to other server environments.

### dvt_service

```bash
make makeService
```

### dvt_tools

```bash
make makeTools
```

or

```shell
cd tools
go build -ldflags "-X main.apiUrl=http://localhost:3000" -o ../out/dvt_tools .
```

> Please use `product:3000` instead of the localhost address information


