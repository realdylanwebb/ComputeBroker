import subprocess, sys

from dclib.SharedObject import SharedObject


class SharedEnvironment:
    def __init__(self, requirements=None, version=3.8):
        self.__requirements = requirements
        self.__version = version
        self.__functions = []
        self.__shared_objects = []

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

        if not isinstance(obj, type):
            self.__shared_objects.append(SharedObject(obj))
            self.__shared_objects[len(self.__shared_objects) - 1].id = len(self.__shared_objects)
        else:
            self.__shared_objects.append(SharedObject(obj()))
            self.__shared_objects[len(self.__shared_objects) - 1].id = len(self.__shared_objects)
