# ComputeBroker
Find workers for your distributed computing problem using ComputeBroker.

Created for Hack the U 2020.

ComputeBroker is a web service and client that allows clients to accept and distribute workloads.
It handles client connection information and can associate a number of clients to form a distributed computing network.

##### The Alpha Server is currently hosted at vps295572.vps.ovh.ca

[API Documentation](#api-documentation)

[Client Documentation](#client-documentation)

## API Documentation

Before you get into the in depth API documentation, the client ships with a nice wrapper for API requests.
### Client API Wrapper
### Example Usage
```python
from ServiceAPI import *

# Getting a batch of 5 workers
try:
    service = ServiceAPI("example@example.com", "password", "examplePublicKey", "127.0.0.1:3000", "http://vps295572.vps.ovh.ca")
    service.register()
    service.login()
    workers = service.newSession(5)

# Notifying that you are ready to accept 3 workloads
    service.readyFor = 3
    service.notify()
        
except RequestFailedError as err:
    #Handle failed requests here
```

### Class ServiceAPI(email::string, password::string, publicKey::string, localAddress::string, serviceAddress::string)
* email: Email used for login credentials
* password: password used for login credentials
* publicKey: the local client's public key
* localAddress: the local service's IP address and port
* serviceAddress: the broker service's IP or URL
* readyFor: indicates the amount of jobs the service is ready for
### ServiceAPI.register()
Registers a the user credentials with the broker service
### ServiceAPI.login()
Retrieves the API key associated with the credentials and stores it internally
### ServiceAPI.notify()
Syncronizes local and broker service ready for values
### ServiceAPI.newSession(workers::int) returns [{address: string, pubKey: string}]
Creates a new session in the broker service and returns a session key and associated worker addresses and public keys
### ServiceAPI.getSession(sessionKey::string) returns [{address: string, pubKey: string}]
Retrieves the worker addresses and public keys associated with an existing sessionKey


## API Endpoints
If you'd rather make your own API requests, here's the raw API endpoints.
All request and response bodies are in JSON format.
### Registering a client
```
POST /client
```
Body:
```
{
  email: string
  password: string
  pubKey: string
  address: string
}
```
* email: email used your account
* password: account password, this is hashed using SHA-256 on the backend so rest assured :)
* pubKey: client public key
* address: public IP address and port the service is running at in form address:port

Response:
```
{
  email: string
  password: string
  pubKey: string
  address: string
}
```
This is just an echo of the submission for confirmation purposes.
### Retrieving your API key
`POST /login`
Body:
```
{
  email: string
  passord: string
}
```
* email: email associated with your account
* password: password associated with your account

Response:
```
{
  token: string
}
```
* token: the API key for your account, this has a TTL embedded so it will only be valid for a set amount of time.

### Signaling that you are ready for a workload
```
POST /client/signal/{available}
```
Headers:
```
Authorization: Bearer {apiKey}
```
* available: the amount of workloads you are ready to accept
* apiKey: your API key, obtainable using POST /login

### Creating a work session
`POST /session`
Body:
```
{
  workers: int
}
```
Headers:
```
Authorization: Bearer {apiKey}
```
* workers: the amount of workers to connect to
* apiKey: your API key, obtainable using POST /login

Response:
```
{
  workers: [
    {
      pubKey: string
      address: string
    }
  ]
  token: string
}
```
* pubKey: The public key for a worker
* address: The public IP address for a worker
* token: A session token that can be used to access the list of workers again
### Getting a work session
```
POST /session/refresh
```
Body:
```
{
  token: string
}
```
Headers:
```
Authorization: Bearer {apiKey}
```
* token: the session token obtained when creating the session using POST /session
* apiKey: your API key, obtainable using POST /login

Response:
```
{
  workers: [
    {
      pubKey: string
      address: string
    }
  ]
}
```
* pubKey: The public key for a worker
* address: The public IP address for a worker


## Client Documentation
