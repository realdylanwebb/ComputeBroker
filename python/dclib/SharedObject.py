class AssignableSetattr(type):
    def __new__(mcls, name, bases, attrs):
        def __setattr__(self, attr, value):
            object.__setattr__(self, attr, value)

        init_attrs = dict(attrs)
        init_attrs['__setattr__'] = __setattr__

        init_cls = super(AssignableSetattr, mcls).__new__(mcls, name, bases, init_attrs)

        real_cls = super(AssignableSetattr, mcls).__new__(mcls, name, (init_cls,), attrs)
        init_cls.__real_cls = real_cls

        return init_cls

    def __call__(cls, *args, **kwargs):
        self = super(AssignableSetattr, cls).__call__(*args, **kwargs)
        real_cls = cls.__real_cls
        self.__class__ = real_cls
        return self

def synchronize(_self, _obj):
    # Do shit here
    return _obj

def synchronize_attribute(_self, _attr, _value):
    # Do shit here
    pass

class SharedObject(object):
    __metaclass__ = AssignableSetattr

    def __init__(self, obj):
        self.obj = obj
        self.id = None
        for key, value in synchronize(self, obj).items():
            setattr(self, key, value)

    def __setattr__(self, attr, value):
        synchronize_attribute(self, attr, value)
        super(SharedObject, self).__setattr__(attr, value)

