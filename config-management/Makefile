create-certs:
	docker-compose -f create-certs.yml run --rm create_certs

run:
	docker-compose -f docker-compose.yml up -d es01

set-up:
	docker exec es01 /bin/bash -c "bin/elasticsearch-setup-passwords auto --batch --url https://es01:9200"

run-monstache:
	docker-compose -f docker-compose.yml up -d monstache

run-mongodb:
	docker-compose -f docker-compose.yml up -d mongodb