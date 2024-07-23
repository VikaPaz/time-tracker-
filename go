HOST=localhost
POSTGRES_PORT=5432
USRE=user
PASSWORD=password
DB_NAME=userspackage main

import (
	"github.com/VikaPaz/time_tracker/internal/app"
)

func main() {
	app.Run()
}
