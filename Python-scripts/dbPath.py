from os import getcwd
from os.path import join

# Find the path to db,
def findDBPath():
    """Attemtps to find the path to the database inside
    SwarmSI-Sim/src/data/
    """
    cwd = getcwd()
    foldName = "Python-scripts"
    parentName = "SwarmSI-Sim"
    if cwd[len(cwd)-len(foldName):len(cwd)] == foldName:
        dbpth = join(cwd[:len(cwd)-len(foldName)], "src", "data", "simRes.db")
        return dbpth
    elif cwd[len(cwd)-len(parentName):len(cwd)] == parentName:
        dbpth = join(cwd, "src", "data", "simRes.db")
        return dbpth
    else:
        # Can change to do a search of sorts, but I think
        # this will cover the most of my cases.
        raise "Not found"


if __name__ == "__main__":
    print(findDBPath())