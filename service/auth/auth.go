package auth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/easy-oj/repos/common/database"
)

func Authenticate(username, password, uuid string) (bool, error) {
	var uid, repoUid int
	var dbPassword string
	row := database.DB.QueryRow("SELECT id, password FROM tb_user WHERE username = ?", username)
	if err := row.Scan(&uid, &dbPassword); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if fmt.Sprintf("%x", sha256.Sum256([]byte(password))) != dbPassword {
		return false, nil
	}
	row = database.DB.QueryRow("SELECT uid FROM tb_repo WHERE uuid = ?", uuid)
	if err := row.Scan(&repoUid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return uid == repoUid, nil
}
