DO $$
BEGIN
    FOR i IN 1..30 LOOP
        BEGIN
            PERFORM * FROM pg_stat_activity WHERE datname = 'test_db';
            EXIT;
        EXCEPTION
            WHEN OTHERS THEN
                PERFORM pg_sleep(2);
        END;
    END LOOP;
END $$;

DO $$
BEGIN
    FOR i IN 1..30 LOOP
        BEGIN
            IF EXISTS (SELECT 1 FROM pg_publication WHERE pubname = 'pub_for_all_tables') THEN
                EXIT;
            END IF;
            PERFORM pg_sleep(2);
        EXCEPTION
            WHEN OTHERS THEN
                PERFORM pg_sleep(2);
        END;
    END LOOP;
END $$;

CREATE SUBSCRIPTION IF NOT EXISTS sub_connection
CONNECTION 'host=postgres port=5432 user=test_user password=test_password dbname=test_db'
PUBLICATION pub_for_all_tables;