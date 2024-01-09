#cmngb

phonebook=# \d entries;
                      Table "public.entries"
 Column  |         Type          | Collation | Nullable | Default 
---------+-----------------------+-----------+----------+---------
 id      | character varying(20) |           | not null | 
 contact | jsonb                 |           |          | 
Indexes:
    "entries_pkey" PRIMARY KEY, btree (id)
