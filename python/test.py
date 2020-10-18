from dclib.WorkScheduler import WorkScheduler

sched = WorkScheduler(user_id='', user_key='', user_pass='')
sched.start_session(remote_workers=10)

def get_sums(array_a, array_b):
    sums = []
    for i in range(0, len(array_a)):
        sums.append(array_a[i] + array_b[i])

# Generate data sets

nums_a = sched.array()
nums_b = sched.array()

for i in range(0, 100000):
    nums_a.append(i)
    nums_b.append(i)


workload = sched.function(get_sums, fragmentation=8)


sums = workload.run((nums_a, nums_b))

for number in sums:
    print(number)

