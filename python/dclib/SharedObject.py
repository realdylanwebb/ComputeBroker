class CleanSetAttrMeta(type):
    def __new__(mcls, name, bases, attrs):
        def __setattr__(self, attr, value):
            object.__setattr__(self, attr, value)

        init_attrs = dict(attrs)
        init_attrs['__setattr__'] = __setattr__

        init_cls = super(CleanSetAttrMeta, mcls).__new__(mcls, name, bases, init_attrs)

        real_cls = super(CleanSetAttrMeta, mcls).__new__(mcls, name, (init_cls,), attrs)
        init_cls.__real_cls = real_cls

        return init_cls

    def __call__(cls, *args, **kwargs):
        real_setattr = cls.__setattr__
        cls.__setattr__ = object.__setattr__
        self = super(CleanSetAttrMeta, cls).__call__(*args, **kwargs)
        cls.__setattr__ = real_setattr
        return self


class SharedObjectContainer(object):
    __metaclass__ = CleanSetAttrMeta
    __on_mutate = None
    __uuid = None

    def __init__(self, on_mutate=None, uuid=None):
        self.__on_mutate = on_mutate
        self.__uuid = uuid
        super(SharedObjectContainer, self).__init__()

    def __setattr__(self, key, value):
        print('SET {} = {}'.format(key, value))

        if self.__on_mutate is not None:
            self.__on_mutate(self, key, value, self.__uuid)

        super(SharedObjectContainer, self).__setattr__(key, value)


class SharedObject(SharedObjectContainer):
    def __init__(self, value, uuid=None, on_mutate=None):
        if isinstance(value, type):
            self.value = value()

        super(SharedObject, self).__init__(on_mutate, uuid)
