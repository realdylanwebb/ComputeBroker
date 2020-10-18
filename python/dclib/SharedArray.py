class SharedArray:
    def __init__(self, scheduler, data=None):
        self.values = []
        self.__scheduler = scheduler
        self.__update_workers()

    def __update_workers(self):
        for worker in self.__scheduler.workers:
            sock = worker.socket

            # send update to value

    def __len__(self):
        return len(self.values)

    def __getitem__(self, key):
        return self.values[key]

    def append(self, value):
        self.values.append(value)
        self.__update_workers()

    def insert(self, key, value):
        self.values[key] = value
        self.__update_workers()

    def clear(self):
        self.values.clear()
        self.__update_workers()