#!/bin/bash
set -e
echo "#### \n\n #### \n\n #### \n\n CREATE DATEBASE pmdbstore\n #### \n\n #### \n\n #### \n\n"
export PGPASSWORD=postgres123
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "pmdbstore" <<-EOSQL
  CREATE DATABASE pmdbstore;
  GRANT ALL PRIVILEGES ON DATABASE pmdbstore TO "postgres";
EOSQL
