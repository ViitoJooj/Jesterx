import {
  validateEmail,
  validatePassword,
  validatePasswordConfirmation,
  validateRequired,
} from './validation';

describe('Validation Utilities', () => {
  describe('validateEmail', () => {
    test('accepts valid email', () => {
      const result = validateEmail('test@example.com');
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    test('rejects empty email', () => {
      const result = validateEmail('');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Email é obrigatório");
    });

    test('rejects invalid email format', () => {
      const result = validateEmail('invalid-email');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Formato de email inválido");
    });
  });

  describe('validatePassword', () => {
    test('accepts valid password', () => {
      const result = validatePassword('Password123!');
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    test('rejects empty password', () => {
      const result = validatePassword('');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha é obrigatória");
    });

    test('rejects password shorter than 8 characters', () => {
      const result = validatePassword('Pass1!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha deve ter pelo menos 8 caracteres");
    });

    test('rejects password without uppercase', () => {
      const result = validatePassword('password123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha deve conter pelo menos uma letra maiúscula");
    });

    test('rejects password without lowercase', () => {
      const result = validatePassword('PASSWORD123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha deve conter pelo menos uma letra minúscula");
    });

    test('rejects password without number', () => {
      const result = validatePassword('Password!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha deve conter pelo menos um número");
    });

    test('rejects password without special character', () => {
      const result = validatePassword('Password123');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Senha deve conter pelo menos um caractere especial");
    });
  });

  describe('validatePasswordConfirmation', () => {
    test('accepts matching passwords', () => {
      const result = validatePasswordConfirmation('Password123!', 'Password123!');
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    test('rejects empty confirmation', () => {
      const result = validatePasswordConfirmation('Password123!', '');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("Confirme sua senha");
    });

    test('rejects mismatched passwords', () => {
      const result = validatePasswordConfirmation('Password123!', 'Different123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe("As senhas não conferem");
    });
  });

  describe('validateRequired', () => {
    test('accepts non-empty value', () => {
      const result = validateRequired('João', 'Nome');
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    test('rejects empty value', () => {
      const result = validateRequired('', 'Nome');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Nome é obrigatório');
    });

    test('rejects whitespace-only value', () => {
      const result = validateRequired('   ', 'Sobrenome');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Sobrenome é obrigatório');
    });
  });
});
