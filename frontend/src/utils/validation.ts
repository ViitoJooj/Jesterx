export interface ValidationResult {
  isValid: boolean;
  error?: string;
}

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const UPPERCASE_REGEX = /[A-Z]/;
const LOWERCASE_REGEX = /[a-z]/;
const DIGIT_REGEX = /[0-9]/;
const SPECIAL_CHAR_REGEX = /[!@#$%^&*()_+\-=[\]{}:,.?]/;

export function validateEmail(email: string): ValidationResult {
  if (!email) {
    return { isValid: false, error: 'Email is required' };
  }

  if (!EMAIL_REGEX.test(email)) {
    return { isValid: false, error: 'Invalid email format' };
  }

  return { isValid: true };
}

export function validatePassword(password: string): ValidationResult {
  if (!password) {
    return { isValid: false, error: 'Password is required' };
  }

  if (password.length < 8) {
    return { isValid: false, error: 'Password must be at least 8 characters' };
  }

  if (!UPPERCASE_REGEX.test(password)) {
    return { isValid: false, error: 'Password must contain at least one uppercase letter' };
  }

  if (!LOWERCASE_REGEX.test(password)) {
    return { isValid: false, error: 'Password must contain at least one lowercase letter' };
  }

  if (!DIGIT_REGEX.test(password)) {
    return { isValid: false, error: 'Password must contain at least one number' };
  }

  if (!SPECIAL_CHAR_REGEX.test(password)) {
    return { isValid: false, error: 'Password must contain at least one special character' };
  }

  return { isValid: true };
}

export function validatePasswordConfirmation(
  password: string,
  confirmPassword: string
): ValidationResult {
  if (!confirmPassword) {
    return { isValid: false, error: 'Please confirm your password' };
  }

  if (password !== confirmPassword) {
    return { isValid: false, error: 'Passwords do not match' };
  }

  return { isValid: true };
}

export function validateRequired(value: string, fieldName: string): ValidationResult {
  if (!value || value.trim() === '') {
    return { isValid: false, error: `${fieldName} is required` };
  }

  return { isValid: true };
}
