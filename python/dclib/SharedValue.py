from enum import Enum
import sys
import hashlib


class Sync(Enum):
    ONDEMAND = 0
    IMMEDIATE = 1

class SharedValue:
    def __init__(self, value, scheduler):
        self.value = value
        self.__size = sys.getsizeof(self.value)
        self.__bytes = bytearray(self.value)
        self.__scheduler = scheduler
        self.__update_workers()

    def __update_workers(self):
        for worker in self.__scheduler.workers:
            sock = worker.socket

            # send update to value

    def checksum(self):
        hash = hashlib.sha1()
        hash.update(self.__bytes)
        return hash.digest()

    def set(self, value):
        self.__value = value
        self.__size = sys.getsizeof(self.value)
        self.__bytes = bytearray(self.value)
        self.__update_workers()

    def get(self):
        return self.__value

    def type(self):
        return type(self.__value)

    def size(self):
        return len(self.__bytes)


