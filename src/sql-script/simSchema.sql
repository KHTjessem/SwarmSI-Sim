CREATE TABLE IF NOT EXISTS  run (
	runID integer PRIMARY KEY AUTOINCREMENT,
	nodeCount integer,
	runDesc varchar
);

CREATE TABLE IF NOT EXISTS  rounds (
	roundID integer PRIMARY KEY AUTOINCREMENT,
	runID integer,
	round integer,
	roundPayout integer,
	totalPayout integer,
    FOREIGN KEY(runID) REFERENCES run(runID)
);

CREATE TABLE IF NOT EXISTS  node (
	nodeID integer
);

CREATE TABLE IF NOT EXISTS  nround (
	nodeID integer,
	roundID integer,
	earnings integer,
	stake integer,
    FOREIGN KEY(nodeID) REFERENCES node(nodeID),
    FOREIGN KEY(roundID) REFERENCES rounds(roundID)
);


-- INDEX
CREATE INDEX idx_roundID_nround ON nround(roundID)

CREATE INDEX idx_round_runID_roundID on rounds(round,runID,roundID)
