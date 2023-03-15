# A python script to test sqlite

from os import getcwd
from os.path import join as pathJoin
from os.path import exists as fileExists

import sqlite3 as db
from sqlite3 import Error

dbpath = pathJoin(getcwd(), "src/data/simRes.db")
if not fileExists(dbpath):
    print("No database file was found, creating a new one")

conn = None
try:
    conn = db.connect(dbpath)
except Error as e:
    print(e)


with open('src/sql-script/simSchema.sql', 'r') as f:
    sql_script = f.read()
c = conn.cursor()
c.executescript(sql_script)
conn.commit()


for i in range(0, 16384):
    q = "INSERT INTO node(nodeID) VALUES(?)"
    c.execute(q, (i,))

conn.commit()
print("Eyy")