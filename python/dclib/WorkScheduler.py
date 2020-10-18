from dclib.SharedArray import SharedArray
from dclib.SharedFunction import SharedFunction
from dclib.SharedValue import SharedValue

class WorkScheduler:
    def __init__(self, user_id=None, user_key=None, user_pass=None):
        self.workers = []
        '''
            worker {
                "public_key": string,
                "id": string,
                "ip_address": string
            }
        '''

        # Stores the ids for all shared values
        self.value_ids = []

        self.account = {
            "user_id": "",
            "user_key": "",
            "user_pass": ""
        }

        if user_id != None:
            self.account["user_id"] = user_id

        if user_key != None:
            self.account["user_key"] = user_key

        if user_pass != None:
            self.account["user_pass"] = user_pass

    def start_session(self, remote_workers=1):
        # Get account details from server

        # Request workers from server

        # Launch scheduler in new thread to wait for commands
        pass

    def value(self, value):
        return SharedValue(value, self)

    def function(self, func, fragmentation=0):
        return SharedFunction(func, fragmentation, self)

    def array(self, data=[]):

        if len(data) > 0:
            return SharedArray(self, data)

        return SharedArray(self)





