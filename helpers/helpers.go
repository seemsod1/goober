package helpers

func CheckPassword(pass string, tocmp string) bool {
	return pass == tocmp
}
