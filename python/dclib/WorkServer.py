import socket, selectors

class ClientDispatcher:
    def __init__(self, port=11033, host='127.0.0.1'):
        self.port = port
        self.host = host
        self.host_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sel = selectors.DefaultSelector()

    def enable(self):
        self.host_socket.bind((self.host, self.port))
        self.host_socket.listen()
        print('Listening on port ' + self.port)

        self.host_socket.setblocking(False)
        self.sel.register(self.host_socket, selectors.EVENT_READ, data=None)