import selectors, socket, types, pickle
from dclib.SharedEnvironment import SharedEnvironment

clientsock = None
address = None
HEADERSIZE=10
COMMANDSIZE=10

# Called when an object is modified and we are the client
# Send updates to all connected workers
def client_on_mutate(self, attr, value, uuid):
    pass

def worker_on_receive(msg):
    pass

class Worker:
    def __init__(self):
        self.lsock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.shared_env = None

    def worker_on_mutate(self, attr, value, uuid, sock):
        if sock is not None:
            sock.send(bytes('MUT {} {} {}'.format(attr, value, uuid)))

    def listen(self, host='127.0.0.1', port=13309):
        self.lsock.bind((host, port))
        self.lsock.listen()

        while True:
            clientsocket, address = self.lsock.accept()
            print('Accepted connection from ', address)
            clientsocket.send(bytes('CONN_ACCEPT', encoding='utf-8'))

            self.shared_env = SharedEnvironment(isworker=True)
            self.shared_env.worker_on_mutate = self.worker_on_mutate
            self.shared_env.worker_on_receive = worker_on_receive

            full_msg = b''
            new_msg = True
            while True:
                msg = clientsocket.recv(16)
                if new_msg:
                    print("new msg len:", msg[COMMANDSIZE:HEADERSIZE + COMMANDSIZE])
                    msglen = int(msg[COMMANDSIZE:HEADERSIZE + COMMANDSIZE])
                    new_msg = False

                print(f"full message length: {msglen}")

                full_msg += msg

                print(len(full_msg))

                if len(full_msg) - HEADERSIZE - COMMANDSIZE == msglen:
                    print("full msg recvd")
                    print(full_msg[COMMANDSIZE + HEADERSIZE:])
                    command = full_msg[:COMMANDSIZE].decode('utf-8')

                    print('Command: ' + command)

                    # Add object
                    if command == 'ADD OBJ   ':
                        obj = pickle.loads(full_msg[COMMANDSIZE + HEADERSIZE:])
                        print(obj)
                        self.shared_env.add_object(obj)

                    if command == 'ADD FUNC  ':
                        func = eval(full_msg[COMMANDSIZE + HEADERSIZE:].decode(encoding='utf-8'))
                        print(func)

                    new_msg = True
                    full_msg = b""






