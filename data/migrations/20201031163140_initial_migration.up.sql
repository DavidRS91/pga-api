CREATE TABLE IF NOT EXISTS players( 
    id SERIAL, 
    name TEXT NOT NULL, 
    score INTEGER, 
    is_cut BOOLEAN,
    CONSTRAINT players_pkey PRIMARY KEY (id)
    );