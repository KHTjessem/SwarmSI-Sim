import sqlite3
import numpy as np
from dbPath import findDBPath
from dataclasses import dataclass


def calcGiniList(x:list):
    """
    Generic calculation of Gini coefficient
    using a list of numbers.
    """
    tot = 0
    avg = 0
    for i, xi in enumerate(x):
        avg += xi
        for j, xj in enumerate(x):
            if xi==xj:
                continue
            tot += abs(xi-xj)
    avg = avg/len(x)
    den = 2 * pow(len(x), 2) * avg

    gini = tot/den
    return gini


@dataclass
class roundStats:
    nodes: list
    sumEarnings: int


# def calcGini(rs: roundStats):
#     """
#     Calculates gini coefficient of nodes in 
#     the roundStat node list
#     """
#     numer = 0
#     for i, xi in enumerate(rs.nodes):
#         for j, xj in enumerate(rs.nodes):
#             numer += abs(xi[1]-xj[1])
#     denom = 2 * pow(len(rs.nodes), 2) * (rs.sumEarnings/len(rs.nodes))
#     return numer/denom


def calcGini(rs: roundStats) -> float:
    """
    Calculates gini coefficient of nodes in 
    the roundStat node list
    """
    earnings = np.array([x[1] for x in rs.nodes])
    n = len(earnings)
    numer = np.sum(np.abs(np.subtract.outer(earnings, earnings)))
    denom = 2 * pow(n, 2) * (rs.sumEarnings/n)
    return numer/denom

def getRoundStats(db:sqlite3.Connection, runID,round) -> roundStats:
    q="""
    SELECT nodeID, earnings, stake
    FROM nround 
    WHERE roundID=?
    """
    qsum = "SELECT sum(earnings) FROM nround WHERE roundID=?"

    cur = db.cursor()
    r = cur.execute(q, (round,))
    res = r.fetchall()

    s = cur.execute(qsum, (round,))
    s = s.fetchone()
    return roundStats(res, s[0])


def caclGiniSql(db:sqlite3.Connection, roundID, col:str="earnings") -> float:
    """Calculates earnings gini with sql query
    col: name of column to calculate gini from. 
    """

    q = f"""
    SELECT  1-2 * sum(({col} * (rownum-1) + cast({col} as float)/2 )) / count(*) / sum({col}) 
    AS gini
    FROM
    (
        SELECT nodeID, {col}, row_number() OVER (
            ORDER BY {col} DESC
        ) rownum
        FROM nround WHERE roundID=?
    )
    """
    cur = db.cursor()
    r = cur.execute(q, (roundID,))
    return r.fetchone()[0]


## Gini for stake should be approx 0 (+- for rounding errors) if calculating it on groups
## becase the sum of stake inside a groupd should equal the stake of 
## of the first node.
def calcBucketGini(db: sqlite3.Connection, roundID, col="earnings"):
    q = f"""
        SELECT  1-2 * sum(({col} * (rownum-1) + cast({col} as float)/2 )) / count(*) / sum({col}) 
        AS gini
        FROM
        (
            SELECT nodeID, sum({col}) as {col}, row_number() OVER (
                ORDER BY sum({col}) DESC
            ) rownum		
            FROM nround NATURAL JOIN nodeGroup 
            WHERE roundID=? GROUP BY groupID
        )
    """
    cur = db.cursor()
    r = cur.execute(q, (roundID,))
    return r.fetchone()[0]

if __name__ == "__main__":
    # Should equal 0.3
    x = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    print(calcGiniList(x))

    # some roundID
    rid = 123

    #sql connection:
    con = sqlite3.connect(findDBPath("simRes.db"))
    
    asd = getRoundStats(con, 0, rid)

    res = calcGini(asd)
    print(f'Gini for roundID {rid}:\t\t{res}')

    resQ = caclGiniSql(con, rid)
    print(f'From sql for roundID {rid}:\t{resQ}')