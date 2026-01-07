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
    return { isValid: false, error: "Email é obrigatório" };
  }

  if (!EMAIL_REGEX.test(email)) {
    return { isValid: false, error: "Formato de email inválido" };
  }

  return { isValid: true };
}

export function validatePassword(password: string): ValidationResult {
  if (!password) {
    return { isValid: false, error: "Senha é obrigatória" };
  }

  if (password.length < 8) {
    return { isValid: false, error: "Senha deve ter pelo menos 8 caracteres" };
  }

  if (!UPPERCASE_REGEX.test(password)) {
    return { isValid: false, error: "Senha deve conter pelo menos uma letra maiúscula" };
  }

  if (!LOWERCASE_REGEX.test(password)) {
    return { isValid: false, error: "Senha deve conter pelo menos uma letra minúscula" };
  }

  if (!DIGIT_REGEX.test(password)) {
    return { isValid: false, error: "Senha deve conter pelo menos um número" };
  }

  if (!SPECIAL_CHAR_REGEX.test(password)) {
    return { isValid: false, error: "Senha deve conter pelo menos um caractere especial" };
  }

  return { isValid: true };
}

export function validatePasswordConfirmation(
  password: string,
  confirmPassword: string
): ValidationResult {
  if (!confirmPassword) {
    return { isValid: false, error: "Confirme sua senha" };
  }

  if (password !== confirmPassword) {
    return { isValid: false, error: "As senhas não conferem" };
  }

  return { isValid: true };
}

export function validateRequired(value: string, fieldName: string): ValidationResult {
  if (!value || value.trim() === '') {
    return { isValid: false, error: `${fieldName} é obrigatório` };
  }

  return { isValid: true };
}
