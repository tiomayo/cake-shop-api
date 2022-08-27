migrate_up:
	migrate -path scripts/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/cake-shop" up
migrate_down:
	migrate -path scripts/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/cake-shop" down