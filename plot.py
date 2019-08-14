from matplotlib import pyplot as plt

fig = plt.figure()
ax = fig.add_subplot(111)

with open("route.txt") as f:
    content = f.readlines()
content = [x.strip() for x in content]

X = [int(i) for i in content[0].split(',')]
Y = [int(i) for i in content[1].split(',')]
P = [int(i) for i in content[2].split(',')]

print X
print Y
print P

plt.plot(X,Y)

for i, txt in enumerate(P):
    ax.annotate(txt, (X[i], Y[i]))

plt.grid()
plt.show()