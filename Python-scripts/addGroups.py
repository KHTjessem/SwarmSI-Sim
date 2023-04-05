import sqlite3
from dbPath import findDBPath

# // Keeps track of bukects internally. Bucket B_i has i nodes.
# // B_1 has 1 nodes, B_2 has 2 etc.
# func (bss bucketSumStake) GetStake(nodeID int) int {
# 	// st := bss.stake / bss.bucket

# 	st := bss.stake / *bss.bucket

# 	*bss.nc++
# 	if *bss.nc == *bss.bucket {
# 		*bss.bucket++
# 		*bss.nc = 0
# 	}

# 	return st
# }


con = sqlite3.connect(findDBPath("simRes2.db"))
q = "INSERT INTO nodeGroup(nodeID, groupID) VALUES(?,?)"


bucket = 1
nc = 0

for i in range(2080):
    nc += 1
    cur = con.cursor()
    r = cur.execute(q, (i, bucket))
    if nc == bucket:
        bucket += 1
        nc = 0

con.commit()

print("Done adding node groups.")


