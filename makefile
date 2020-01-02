CNT = rb-tracker__mongo-db
CMD = /bin/bash

PHONY: run build

build:
	docker pull andresvidal/rpi3-mongodb3:latest
	@echo Clean pull succesful.

run: build
	docker run -d --name $(CNT) --restart unless-stopped -p 27017:27017 -p 28017:28017 andresvidal/rpi3-mongodb3:latest mongod --auth
	sleep 10
	@echo "Container started..."
	docker exec -it $(CNT) mongo admin --eval "db.createUser({user: 'admin', pwd: 'admin', roles:[{role:'root',db:'admin'}]});"
	sleep 10
	@echo "Admin user generated..."
	docker exec -it $(CNT) mongo rb-tracker --authenticationDatabase admin -u admin -p admin --eval "db.getSiblingDB('rb-tracker');db.createUser({user: 'msteger', pwd: 'msteger', roles:[{role:'dbOwner',db:'rb-tracker'}]});"
	@echo "DB and DB-User setup!"
	docker-compose up -d
	@echo "Dockerized Environment started (as daemon)."

kill:
	docker-compose down
	docker stop $(CNT)
	@echo "Dockerized Environment stopped."

clean: kill
	docker rm $(CNT)
	docker network prune
	docker volume prune
	@echo "Docker network and volume pruned."

exec:
	docker exec -it $(CNT) $(CMD)