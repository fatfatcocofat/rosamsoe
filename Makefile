jwtkey:
	node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"

watch:
	air

docs:
	swag init

dev:
	go run .