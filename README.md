#cmngb

        phonebook=# \d entries;
                            Table "public.entries"
        Column  |         Type          | Collation | Nullable | Default 
        ---------+-----------------------+-----------+----------+---------
        id      | character varying(20) |           | not null | 
        contact | jsonb                 |           |          | 
        Indexes:
            "entries_pkey" PRIMARY KEY, btree (id)


See docker image: 
    docker run -it --name cmngb donetsmaksim/cmngb
It's a simple service without postgresql.

Postgresql and separate users is coming...