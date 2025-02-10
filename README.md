# cth

Below are the step-by-step instructions I used to test my solution.

You need a C compiler for the SQLite module. See [here](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#go-sqlite3)
for more information.

```sh
podman compose --file docker-compose.yml up --detach

go mod tidy

# Optional: Run the unit tests.
make test

make

# The SQLite migrations.
make down && make up

PUBSUB_EMULATOR_HOST=localhost:8085 ./cth

podman compose down

# Optional: Clean up the binary and database.
make clean && make down
```

## Some Thoughts
I reached for SQLite and relied on its capability to enforce a unique constraint for our 3-tuple.
The schema can be seen in create_table.sql. Our Go program creates a client that then listens on the
subscription until it is interrupted. We have a separate goroutine that handles the database
insertions since we only want one writer. The inserts are wrapped in a transaction which may be a
bit heavy-handed but easily adjusted. Perhaps we could batch records in the future? I did make the
channel that submits new records buffered so I could control the throughput a little more.

Google's Pub-Sub model is a bit different than previous message queues I have used. The idea of a
subscription is different than Kafka's consumer groups for example.

It was a fun exercise!
