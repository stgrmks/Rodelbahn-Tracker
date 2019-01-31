CNT = rb-tracker__mongodb
CMD = /bin/bash

PHONY: run build

build:
	docker pull mongo:latest
	@echo Clean pull succesful.

run: build
	docker-compose up -d
	@echo "Dockerized Environment started (as daemon)."

kill:
	docker-compose down
	@echo "Dockerized Environment stopped."

clean: kill
	docker network prune
	docker volume prune
	@echo "Docker network and volume pruned."

exec:
	docker exec -it $(CNT) $(CMD)
