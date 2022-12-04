
# Prerequisites

```bash
MONGO_HOST="localhost"
MONGO_PORT="27017"

MONGO_USER='root'
MONGO_PASSWD='rootpassword'

AUTH_PARAM="--authenticationDatabase admin --username ${MONGO_USER} --password ${MONGO_PASSWD} "

DUMP_DIR="./../dumps/"
```

# Dump all Collections of Database

```bash
mongodump --host ${MONGO_HOST} --port ${MONGO_PORT} ${AUTH_PARAM} --out ${DUMP_DIR}
```

# Import Collection Dump of Database

```bash
mongorestore --uri="mongodb://localhost:27017" -u root -p rootpassword ./../dumps/
```