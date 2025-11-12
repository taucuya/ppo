import matplotlib.pyplot as plt
from matplotlib.figure import Figure

with open("./benchmark/grad.txt", "r") as file:
    lines = file.readlines()

X = []
Y = []
for line in lines:
    x, y = line.strip().split()
    if x[-2:] != 'ms':
        x = float(x[:-1]) * 1000
    else:
        x = float(x[:-2])
    X.append(x)
    Y.append(int(y))

plt.plot(Y, X)
plt.ylabel("Время, мс")
plt.xlabel("Пользователи, шт")
plt.title("График измеряемой величины во времени")
plt.show()

with open("./benchmark/percentiles.txt", "r") as file:
    lines = file.readlines()

p50 = []
p75 = []
p90 = []
p95 = []
p99 = []
y = []

for line in lines:
    p = []
    for i in line.strip().split():
        if i[-1:] != 's':
            i = int(i)
        elif i[-2:] != 'ms':
            i = float(i[:-1]) * 1000
        else:
            i = float(i[:-2])
        p.append(i)
    p50.append(p[0])
    p75.append(p[1])
    p90.append(p[2])
    p95.append(p[3])
    p99.append(p[4])
    y.append(p[5])

fig, axes = plt.subplots(nrows=2, ncols=3, figsize=(10, 10))

axes[0, 0].hist(p50, color="red")
axes[0, 0].set_title("p50")
axes[0, 0].set_xlabel("Время, мс")
axes[0, 1].hist(p75, color="red")
axes[0, 1].set_title("p75")
axes[0, 1].set_xlabel("Время, мс")
axes[0, 2].hist(p90, color="red")
axes[0, 2].set_title("p90")
axes[0, 2].set_xlabel("Время, мс")
axes[1, 0].hist(p95, color="red")
axes[1, 0].set_title("p95")
axes[1, 0].set_xlabel("Время, мс")
axes[1, 1].hist(p99, color="red")
axes[1, 1].set_title("p99")
axes[1, 1].set_xlabel("Время, мс")
axes[1, 2].hist(X, color="red")
axes[1, 2].set_title("p100")
axes[1, 2].set_xlabel("Время, мс")
plt.tight_layout()
plt.show()

