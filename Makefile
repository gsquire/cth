.PHONY: build
build:
	go build

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test -v ./decode

.PHONY: vet
vet:
	go vet ./...

# SQLite migration rules.
DB := scan-results.sqlite

.PHONY: down
down:
	rm -f ${DB}

.PHONY: up
up:
	touch ${DB}
	sqlite3 ${DB} < create_table.sql
