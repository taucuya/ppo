#!/bin/bash
set -e

echo "Waiting for primary database to be ready..."

until pg_isready -h postgres -U test_user -d test_db; do
  echo "Waiting for primary database..."
  sleep 2
done

echo "Primary database is ready, setting up subscription..."

psql -v ON_ERROR_STOP=1 -U test_user -d test_db <<-EOSQL
    CREATE SUBSCRIPTION sub_connection
    CONNECTION 'host=postgres port=5432 user=test_user password=test_password dbname=test_db'
    PUBLICATION pub_for_all_tables;
EOSQL

echo "Replication setup completed!"