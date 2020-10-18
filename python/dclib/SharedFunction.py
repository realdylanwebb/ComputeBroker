import inspect

class SharedFunction:
    def __init__(self, exec):
        self.__exec = exec
        self.__environment = None
        self.__source = inspect.getsource(exec)



