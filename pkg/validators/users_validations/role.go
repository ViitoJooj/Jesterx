package validators

import "errors"

func Role(role string) error {
	switch role {
	case "super_admin":
		return nil
	case "admin":
		return nil
	case "owner":
		return nil
	case "manager":
		return nil
	case "editor":
		return nil
	case "designer":
		return nil
	case "developer":
		return nil
	case "support":
		return nil
	case "stock":
		return nil
	case "shipping":
		return nil
	case "product_manager":
		return nil
	case "order_manager":
		return nil
	case "teacher":
		return nil
	case "student":
		return nil
	case "event_owner":
		return nil
	case "event_manager":
		return nil
	case "checkin":
		return nil
	case "marketer":
		return nil
	case "finance":
		return nil
	}

	return errors.New("This role not exists.")
}
