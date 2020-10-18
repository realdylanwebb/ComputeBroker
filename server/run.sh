#Needed because the sql driver were using has some C code in it :)
export CGO_ENABLED=1

#SQLite database file name, :memory: for an in memory database
export DBNAME=""

#Host for the service to run as
export SERVICEHOST=""

#Path to the encryption key, tokens are encrypted using A128GCM
export AUTHENCKEY=""

#Path to the signing key, tokens are signed using HS512
export AUTHSIGKEY=""

#Time to live for signed tokens in form (int)(h || m || s || d)
export TOKENTTL="12h"

./server