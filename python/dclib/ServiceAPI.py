import requests

#RequestFailedError is raised when a http request to the broker service
#fails, it contains the status code and message from the error response.
class RequestFailedError(Exception):
    def __init__(self, statusCode, message):
        self.statusCode = statusCode
        self.message = message

#ServiceAPI contains fields and methods needed to interact with the broker service
class ServiceAPI:
    def __init__(self, email, password, publicKey, localAddress, serviceAddress):
        self.email = email
        self.password = password
        self.sessionKeys = []
        self.apiKey = ""
        self.readyFor = 0
        self.publicKey = publicKey
        self.localAddress = localAddress
        self.serviceAddress = serviceAddress

    #register will post client information to the broker service, raises an error on a failed request or registration
    def register(self):
        res = requests.post(self.serviceAddress+"/client", json = {"email": self.email, "password": self.password, "pubKey": self.publicKey, "address": self.localAddress})
        if (res.status_code != 200 and res.status_code != 201):
            body = res.json()
            raise RequestFailedError(res.status_code, body["error"])

    #login attempts to login to the broker service. On success, self.apiKey will be set to the users api key
    #raises an error on a failed request or invalid login credential
    def login(self):
        res = requests.post(self.serviceAddress+"/login", json = {"email":self.email, "password": self.password})
        body = res.json()

        if res.status_code != 200:
            raise RequestFailedError(res.status_code, body["error"])
        else:
            self.apiKey = body["token"]

    #notify syncronizes the server with the client's self.readyFor
    #raises an error on a failed request
    def notify(self):
        res = requests.post(self.serviceAddress+"/client/signal/"+str(self.readyFor), headers = {"Authorization": "bearer " + self.apiKey})
        if res.status_code != 200:
            body = res.json()
            raise RequestFailedError(res.status_code, body["error"])

    #newSession requests information for a batch of workers from the server.
    #workers are of the form {pubKey: string, address: string}.
    #raises an error on a failed request or when the self.apiToken is invalid.
    def newSession(self, numWorkers):
        res = requests.post(self.serviceAddress+"/session", headers = {"Authorization": "bearer " + self.apiKey}, json = {"workers": numWorkers})
        body = res.json()
        if (res.status_code != 200 and res.status_code != 201):
            raise RequestFailedError(res.status_code, body["error"])
        else:
            self.sessionKeys.append(body["token"])
            return body["workers"]

    #getSession requests information for an already existing batch of workers,
    #associated with sessionToken.
    #workers are of the form {pubKey: string, address: string}.
    #raises an error on a failed request or when the self.apiToken is invalid.
    def getSession(self, sessionToken):
        res = requests.post(self.serviceAddress+"/session/refresh", headers = {"Authorization": "bearer " + self.apiKey}, json = {"token": sessionToken})
        body = res.json()
        if (res.status_code != 200 and res.status_code != 201):
            raise RequestFailedError(res.status_code, body["error"])
        else:
            return body["workers"]