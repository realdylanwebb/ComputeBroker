from dclib.SharedObject import SharedObject

a = SharedObject(int)

a.value = 5

for i in range(0, 10):
    a.value += i
    print(a)
