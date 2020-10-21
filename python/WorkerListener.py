import io
import os
import socket
import select
import threading

class WorkerListener:
    def __init__(self, localHost, localPort, readChunkSize, connectionHandler):
        self.connectionHandler = connectionHandler
        self.readChunkSize = readChunkSize
        
        #Create non blocking socket
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setblocking(0)
        self.sock.bind((localHost, localPort))
        self.sock.listen()

        #Create instance of UNIX poll
        self.poller = select.poll()
        self.poller.register(self.sock, select.POLLIN)

        #Initialize thread and buffer dictionaries
        self.threads = {}
        self.threadBuffers = {}

    def run(self):
        while True:
            ev = self.poller.poll()
            for fd, event in ev:
                if fd == self.sock.fileno():
                    #Accept and register connection
                    conn, addr = self.sock.accept()
                    self.poller.register(conn, select.POLLIN)

                    #Init thread buffer and start thread
                    self.threadBuffers[conn.fileno()] = io.BytesIO()
                    t = threading.Thread(target=self.connectionHandler, args=(conn.fileno(), self.threadBuffers[conn.fileno()])) 
                    t.start()
                    self.threads[conn.fileno()] = t
                else:
                    #Find corresponding thread buffer and write to buffer
                    buf = self.threadBuffers[fd]
                    buf.write(os.read(fd, self.readChunkSize))
                    

