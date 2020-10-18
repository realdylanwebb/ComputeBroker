# ComputeBroker
Find workers for your distributed computing problem using ComputeBroker.

Created for Hack the U 2020.

ComputeBroker is a web service and client that allows clients to accept and distribute workloads.
It handles client connection information and can associate a number of clients to form a distributed computing network.

[API Documentation](#api-documentation)

[Client Documentation](#client-documentation)

## API Documentation
All request and response bodies are in JSON
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
