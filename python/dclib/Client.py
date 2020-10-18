import socket
from dclib.SharedEnvironment import SharedEnvironment

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((socket.gethostname(), 13309))
env = SharedEnvironment(socket=s)

x = env.add_object(int)
x.value = 5

for i in range(0, 10):
    x.value *= 10

def add(a, b):
    return a + b

env.add_function(add)


