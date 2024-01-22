#cmngb

        phonebook=# \d test_table_entries;
                        Table "public.test_table_entries"
        Column  |         Type          | Collation | Nullable | Default 
        ---------+-----------------------+-----------+----------+---------
        id      | character varying(20) |           | not null | 
        contact | jsonb                 |           |          | 
        Indexes:
            "test_table_entries_pkey" PRIMARY KEY, btree (id)



        phonebook=# \d users;
                                            Table "public.users"
        Column   |         Type          | Collation | Nullable |                Default                 
        -----------+-----------------------+-----------+----------+----------------------------------------
        user_id   | integer               |           | not null | nextval('users_user_id_seq'::regclass)
        user_name | character varying(50) |           |          | 
        secret    | character varying(50) |           |          | 
        Indexes:
            "users_pkey" PRIMARY KEY, btree (user_id)
            "users_user_name_key" UNIQUE CONSTRAINT, btree (user_name)
