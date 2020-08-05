# docker-compose up
.PHONY: up
up:
	docker-compose up --build

# docker-compose down
# volumeのお掃除
.PHONY: down
down:
	docker-compose down && \
	docker volume prune --force && \
	docker-compose kill && \
	docker network prune --force

