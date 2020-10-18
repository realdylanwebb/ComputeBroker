class SharedFunction:
    # When fragmentation is 0, the data split is automatic based on the number of workers
    def __init__(self, func, scheduler, fragmentation=0):
        self.func = func
        self.fragmentation = fragmentation
        self.scheduler = scheduler

    def execute(self, parameters):
        # send execute command to workers, giving the value ids for the parameters
        pass
