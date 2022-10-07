confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

run/api:
	go run ./cmd/api

db/psql:
	psql ${APOD_DB_DSN}

db/migrations/new:
	@echo "Creating migration fiels for ${name}"
	migrate create -seq -ext=.sql -dir=./migrations ${name}

db/migrations/up: confirm
	@echo "Runnig up migrations..."
	migrate -path ./migrations -database ${APOD_DB_DSN} up
