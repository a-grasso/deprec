
# Prerequisites

## Install Necessary Packages
mongosh: https://www.mongodb.com/docs/mongodb-shell/install/

## Environment Variables for Backup / Restore
```bash
MONGO_HOST="localhost"
MONGO_PORT="27017"

MONGO_USER='root'
MONGO_PASSWD='rootpassword'

AUTH_PARAM="--authenticationDatabase admin --username ${MONGO_USER} --password ${MONGO_PASSWD} "

DUMP_DIR="./../dumps/"

CONNECTION="--host ${MONGO_HOST} --port ${MONGO_PORT} ${AUTH_PARAM}"
```

# Drop all Databases
```bash
mongosh ${CONNECTION} ${DATABASE} --eval '
db.adminCommand("listDatabases").databases.
   map(d => d.name).
   filter(n => ["admin", "config", "local"].indexOf(n) == -1 ).
   map(n => db.getSiblingDB(n).dropDatabase());
   '
```

# Dump all Collections of Database

```bash
mongodump ${CONNECTION} --out ${DUMP_DIR}
```

# Import Dumps

```bash
mongorestore ${CONNECTION} ${DUMP_DIR}
```