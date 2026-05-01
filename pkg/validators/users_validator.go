package validators

func User(name, email, password, role, cpf string) error {
	if err := Name(name); err != nil {
		return err
	}

	if err := Email(email); err != nil {
		return err
	}

	if err := Password(password); err != nil {
		return err
	}

	if err := Role(role); err != nil {
		return err
	}

	if err := Cpf(cpf); err != nil {
		return err
	}

	return nil
}
