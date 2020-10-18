import subprocess, sys

from dclib.SharedObject import SharedObject


def obj_setattr(self, attr, value):
    print('SET {} = {}'.format(attr, value))
    super(type(self), self).__setattr__(attr, value)

# Called when an object is modified and we are the client
# Send updates to all connected workers
def client_on_mutate(self, attr, value, uuid):
    pass

# Called when an object is modified and we are a worker
# Should send an update to the client
def worker_on_mutate(self, attr, value, uuid):
    pass


class SharedEnvironment:
    def __init__(self, requirements=None, version=3.8, isworker=False):
        self.__requirements = requirements
        self.__version = version
        self.__functions = []
        self.__shared_objects = []
        self.__isworker = isworker

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
            self.__shared_objects.append(SharedObject(obj, uuid=len(self.__shared_objects), on_mutate=worker_on_mutate))
        else:
            self.__shared_objects.append(SharedObject(obj, uuid=len(self.__shared_objects), on_mutate=client_on_mutate))

