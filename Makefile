#  variables block
COMPOSE_FILE="./docker-compose.yaml"

# docker-compose commands
run:
	docker-compose -f ${COMPOSE_FILE} up -d

ls:
	docker-compose -f ${COMPOSE_FILE} ps

stop:
	docker-compose -f ${COMPOSE_FILE} down

remove: stop
	docker prune -f

remove-all: remove
	docker volume prune -f

# docker commands
logs:
	docker logs ${SERVICE}

# application sprcific commands
run-app:
	docker-compose -f ${COMPOSE_FILE} up -d mf_app

migrate:
	docker-compose -f ${COMPOSE_FILE} up -d flyway

db-connect:
	docker-compose -f ${COMPOSE_FILE} exec mf_db psql -h localhost -U postgres -d mf -W
