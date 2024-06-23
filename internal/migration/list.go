package migration

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
)

var Migrations = []*gormigrate.Migration{
	&M20240623164400_user,
	&M20240623164401_room,
	&M20240623164402_room_user,
}
