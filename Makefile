rm:
	sudo docker compose stop \
  	&& 	docker compose rm \

up:
	sudo docker compose up --force-recreate