CREATE TABLE IF NOT EXISTS players( 
    id SERIAL, 
    name TEXT NOT NULL, 
    score INTEGER, 
    is_cut BOOLEAN,
    tsn_id INTEGER,
    CONSTRAINT players_pkey PRIMARY KEY (id)
    );