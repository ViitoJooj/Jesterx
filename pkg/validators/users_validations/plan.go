package users_validations

func Plan(plan string) error {

	switch plan {
	case "free":
		return nil
	case "starter":
		return nil
	case "pro":
		return nil
	case "business":
		return nil
	case "enterprise":
		return nil
	case "tester":
		return nil
	}

	return nil
}
