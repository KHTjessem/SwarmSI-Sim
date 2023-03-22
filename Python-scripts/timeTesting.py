import time
import sqlite3
from dbPath import findDBPath
import matplotlib.pyplot as plt
import calcGiniCoefficient as cgc


con = sqlite3.connect(findDBPath("simRes.db"))

# lets get gini coefficinet of 16 nodes static network 
# (runID 1 for me) for 2000 rounds

startRound = 350667
endRound = 358667

gini1 = []

start = time.time()
for i in range(startRound, endRound):
    tmp = cgc.getRoundStats(con, 0, i)
    gini1.append(cgc.calcGini(tmp))
end = time.time()
runTime1 = end-start


gini2 = []
start = time.time()
for i in range(startRound, endRound):
    gini2.append(cgc.caclGiniSql(con, i))
end = time.time()
runTime2 = end-start

print(f"Using roundstats runtime: {runTime1}\nUsing sql runtime: {runTime2}")

print("Checking if gini's are matching")
for i in range(len(gini1)):
    if round(gini1[i], 5) != round(gini2[i], 5):
        print(f"There is a potential missmatch: {gini1[i]} != {gini2[i]}")
print("Done checking")

plt.subplot(2, 1, 1)
plt.plot(gini1)

plt.subplot(2, 1, 2)
plt.plot(gini2)
plt.show()




