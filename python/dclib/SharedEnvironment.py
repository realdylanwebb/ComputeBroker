import subprocess, sys, pickle

from dclib.SharedObject import SharedObject
import inspect

HEADERSIZE=10


def obj_setattr(self, attr, value):
    print('SET {} = {}'.format(attr, value))
    super(type(self), self).__setattr__(attr, value)




class SharedEnvironment:
    def __init__(self, requirements=None, version=3.8, isworker=False, socket=None):
        self.__requirements = requirements
        self.__version = version
        self.__functions = []
        self.__shared_objects = []
        self.__isworker = isworker
        self.__socket = socket

        self.client_on_mutate = None
        self.worker_on_mutate = None
        self.client_on_receive = None
        self.worker_on_receive = None

    def get_requirements(self, logging='verbose'):
        """
        Installs the packages required for the SharedEnvironment.

        :param logging: The logging level to use during package installation.
        :return: None
        """

        for requirement in self.__requirements:
            if requirement['version'] is not None:
                if logging == 'verbose':
                    print('Installing package {}=={}'.format(requirement['name'], requirement['version']))
                subprocess.check_call([sys.executable, '-m', 'pip', 'install',
                                       '{}=={}'.format(requirement['name'], requirement['version'])])
            else:
                if logging == 'verbose':
                    print('Installing package {}'.format(requirement['name']))
                subprocess.check_call([sys.executable, '-m', 'pip', 'install',
                                       '{}'.format(requirement['name'])])

    def add_object(self, obj):
        """
        Create a new SharedObject. SharedObjects will be synchronized with connected worker instances.
        SharedObjects should be capable of wrapping most python classes.

        :param obj: The item from which to create the new shared object. Can be an instance, or an object type.
        :return:    The new SharedObject instance.
        """
        if self.__isworker:
            self.__shared_objects.append(SharedObject(obj, uuid=len(self.__shared_objects), on_mutate=self.worker_on_mutate))
            self.__socket.send(bytes('ADD OBJ   ', encoding='utf-8'))
            self.__socket.send(bytes(obj.__class__))
            self.__socket.send(bytes(obj.__dict__))
        else:
            self.__shared_objects.append(SharedObject(obj, uuid=len(self.__shared_objects), on_mutate=self.client_on_mutate))
            self.__socket.send(bytes('ADD OBJ   ', encoding='utf-8'))
            msg = pickle.dumps(obj)
            self.__socket.send(bytes(f"{len(msg):<{HEADERSIZE}}", "utf-8")+msg)

        return self.__shared_objects[len(self.__shared_objects) - 1]

    def add_function(self, func):
        if not self.__isworker:
            source = inspect.getsource(func)
            self.__functions.append({
                "name": func.__name__,
                "code": source,
                "id": len(self.__functions)
            })
            msg = 'ADD FUNC  {}'.format(bytes(repr(self.__functions[len(self.__functions) - 1]), encoding='utf-8'))
            self.__socket.send(bytes('ADD FUNC  ', encoding='utf-8'))
            self.__socket.send(bytes(f"{len(msg):<{HEADERSIZE}}", "utf-8") + bytes(msg, encoding='utf-8'))

    def execute_function(self, func_id, params):
        if not self.__isworker:
            self.__socket.send(bytes('FUNC RUN   '))
            self.__socket.send(bytes(func_id))
            self.__socket.send(bytes(repr(params), encoding='utf-8'))
        else:
            # Send results
            results = self.__functions[id](*params)
            self.__socket.send(bytes('FUNC RET'))
            self.__socket.send(bytes(repr(results), encoding='utf-8'))
            self.__socket.send(bytes('END'))

