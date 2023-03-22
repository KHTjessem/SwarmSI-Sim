import time
import sqlite3
from dbPath import findDBPath
import matplotlib.pyplot as plt
import calcGiniCoefficient as cgc


con = sqlite3.connect(findDBPath("simRes2.db"))

# lets get gini coefficinet of 16 nodes static network 
# (runID 1 for me) for 2000 rounds

RUNID=9
# How many rounds to calculate gini for
giniRoundLimit = 25000

q="""
SELECT min(roundID) as minR,
max(roundID) as maxR
FROM rounds
WHERE runID=?
"""
cur = con.cursor()
r = cur.execute(q, (RUNID, ))
res = r.fetchone()

startRound, endRound = res


if endRound - startRound > giniRoundLimit:
    endRound = startRound + giniRoundLimit




giniEarnings = []
giniStake = []
start = time.time()
for i in range(startRound, endRound):
    giniEarnings.append(cgc.caclGiniSql(con, i))
    giniStake.append(cgc.caclGiniSql(con, i, "stake"))

end = time.time()
runTime = end-start

print(f"runtime: {runTime}")


fig, axs = plt.subplots(2, 1)
fig.suptitle(f"gini of round: {startRound}-{endRound}, run {RUNID}")

axs[0].plot(giniEarnings)
axs[0].set_title("Earnings gini")

axs[1].plot(giniStake)
axs[1].set_title("Stake gini")

plt.show()
