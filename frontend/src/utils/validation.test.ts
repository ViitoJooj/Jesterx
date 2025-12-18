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
      expect(result.error).toBe('Email is required');
    });

    test('rejects invalid email format', () => {
      const result = validateEmail('invalid-email');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Invalid email format');
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
      expect(result.error).toBe('Password is required');
    });

    test('rejects password shorter than 8 characters', () => {
      const result = validatePassword('Pass1!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Password must be at least 8 characters');
    });

    test('rejects password without uppercase', () => {
      const result = validatePassword('password123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Password must contain at least one uppercase letter');
    });

    test('rejects password without lowercase', () => {
      const result = validatePassword('PASSWORD123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Password must contain at least one lowercase letter');
    });

    test('rejects password without number', () => {
      const result = validatePassword('Password!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Password must contain at least one number');
    });

    test('rejects password without special character', () => {
      const result = validatePassword('Password123');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Password must contain at least one special character');
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
      expect(result.error).toBe('Please confirm your password');
    });

    test('rejects mismatched passwords', () => {
      const result = validatePasswordConfirmation('Password123!', 'Different123!');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Passwords do not match');
    });
  });

  describe('validateRequired', () => {
    test('accepts non-empty value', () => {
      const result = validateRequired('John', 'First Name');
      expect(result.isValid).toBe(true);
      expect(result.error).toBeUndefined();
    });

    test('rejects empty value', () => {
      const result = validateRequired('', 'First Name');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('First Name is required');
    });

    test('rejects whitespace-only value', () => {
      const result = validateRequired('   ', 'Last Name');
      expect(result.isValid).toBe(false);
      expect(result.error).toBe('Last Name is required');
    });
  });
});
