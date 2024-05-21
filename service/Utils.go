package service

import "log"
import "golang.org/x/crypto/bcrypt"
import "time"

func hashAndSalt(pwd string) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

/*
returns nil error if success
*/
func validateTime(timeStr, expectedLayout string) error {
	var _, err = time.Parse(expectedLayout, timeStr)
	return err
}

func parseToUnixTimestamp(timeStr, layout string) int64 {
	var tm, _ = time.Parse(layout, timeStr)
	return tm.UnixMilli()
}

func formatUnixTimestampToString(unixTimestamp int64, layout string) string {
	var tm = time.UnixMilli(unixTimestamp)
	return tm.Format(layout)
}

// --- Useful notes ---

// strconv.FormatUint(uint64(user.Id), 10) - convert uint64 to string
// strconv.Itoa(id) - convert int to string
