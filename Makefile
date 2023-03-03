test_model:
	go test github.com/sgoldenf/playlist/internal/model/playlist

test_server:
	go test github.com/sgoldenf/playlist/internal/server

compose_database:
	docker compose up --detach

migrate_up:
	migrate -path db/migrations \
    	-database "postgres://sgoldenf:sgoldenf@localhost:5432/playlist?sslmode=disable" up

remove_database:
	migrate -path db/migrations \
            -database "postgres://sgoldenf:sgoldenf@localhost:5432/playlist?sslmode=disable" down ; \
    docker stop playlist_db ; \
    docker rm playlist_db

gen: go_gen

go_gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./api/*.proto
