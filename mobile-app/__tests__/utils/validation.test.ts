/**
 * Validation Utilities Tests
 */

import {
  validateEmail,
  validatePassword,
  validatePhoneNumber,
  validateIDNumber,
  validateTaxNumber,
  validateVATNumber,
  validateRequired,
  validateURL,
  validateMinLength,
  validateMaxLength,
} from '@/lib/utils/validation';

describe('validateEmail', () => {
  it('should validate correct email addresses', () => {
    expect(validateEmail('user@example.com')).toBe(true);
    expect(validateEmail('test.user@company.co.za')).toBe(true);
    expect(validateEmail('admin+filter@domain.org')).toBe(true);
  });

  it('should reject invalid email addresses', () => {
    expect(validateEmail('invalid')).toBe(false);
    expect(validateEmail('user@')).toBe(false);
    expect(validateEmail('@domain.com')).toBe(false);
    expect(validateEmail('user @domain.com')).toBe(false);
    expect(validateEmail('')).toBe(false);
  });
});

describe('validatePassword', () => {
  it('should validate strong passwords', () => {
    const result = validatePassword('MyP@ssw0rd123');
    expect(result.isValid).toBe(true);
    expect(result.errors).toHaveLength(0);
  });

  it('should reject weak passwords', () => {
    const result = validatePassword('weak');
    expect(result.isValid).toBe(false);
    expect(result.errors.length).toBeGreaterThan(0);
  });

  it('should require minimum length', () => {
    const result = validatePassword('Short1!');
    expect(result.isValid).toBe(false);
    expect(result.errors).toContain('Password must be at least 8 characters long');
  });

  it('should require uppercase letter', () => {
    const result = validatePassword('password123!');
    expect(result.isValid).toBe(false);
    expect(result.errors).toContain('Password must contain at least one uppercase letter');
  });

  it('should require lowercase letter', () => {
    const result = validatePassword('PASSWORD123!');
    expect(result.isValid).toBe(false);
    expect(result.errors).toContain('Password must contain at least one lowercase letter');
  });

  it('should require number', () => {
    const result = validatePassword('Password!');
    expect(result.isValid).toBe(false);
    expect(result.errors).toContain('Password must contain at least one number');
  });

  it('should require special character', () => {
    const result = validatePassword('Password123');
    expect(result.isValid).toBe(false);
    expect(result.errors).toContain('Password must contain at least one special character');
  });
});

describe('validatePhoneNumber', () => {
  it('should validate South African phone numbers', () => {
    expect(validatePhoneNumber('+27123456789')).toBe(true);
    expect(validatePhoneNumber('0123456789')).toBe(true);
    expect(validatePhoneNumber('27123456789')).toBe(true);
  });

  it('should reject invalid phone numbers', () => {
    expect(validatePhoneNumber('12345')).toBe(false);
    expect(validatePhoneNumber('abc')).toBe(false);
    expect(validatePhoneNumber('')).toBe(false);
  });
});

describe('validateIDNumber', () => {
  it('should validate correct South African ID numbers', () => {
    expect(validateIDNumber('9001015009087')).toBe(true); // Valid 13-digit ID
  });

  it('should reject invalid ID numbers', () => {
    expect(validateIDNumber('123')).toBe(false); // Too short
    expect(validateIDNumber('12345678901234')).toBe(false); // Too long
    expect(validateIDNumber('abc1234567890')).toBe(false); // Contains letters
    expect(validateIDNumber('')).toBe(false); // Empty
  });

  it('should reject IDs with invalid dates', () => {
    expect(validateIDNumber('0013015009087')).toBe(false); // Invalid month (13)
    expect(validateIDNumber('9000325009087')).toBe(false); // Invalid day (32)
  });
});

describe('validateTaxNumber', () => {
  it('should validate correct tax numbers', () => {
    expect(validateTaxNumber('0123456789')).toBe(true); // 10 digits
  });

  it('should reject invalid tax numbers', () => {
    expect(validateTaxNumber('123')).toBe(false); // Too short
    expect(validateTaxNumber('12345678901')).toBe(false); // Too long
    expect(validateTaxNumber('abc1234567')).toBe(false); // Contains letters
    expect(validateTaxNumber('')).toBe(false); // Empty
  });
});

describe('validateVATNumber', () => {
  it('should validate correct VAT numbers', () => {
    expect(validateVATNumber('4123456789')).toBe(true); // Starts with 4, 10 digits
  });

  it('should reject invalid VAT numbers', () => {
    expect(validateVATNumber('5123456789')).toBe(false); // Doesn't start with 4
    expect(validateVATNumber('412345678')).toBe(false); // Too short
    expect(validateVATNumber('41234567890')).toBe(false); // Too long
    expect(validateVATNumber('')).toBe(false); // Empty
  });
});

describe('validateRequired', () => {
  it('should validate non-empty values', () => {
    expect(validateRequired('value')).toBe(true);
    expect(validateRequired('  text  ')).toBe(true);
  });

  it('should reject empty values', () => {
    expect(validateRequired('')).toBe(false);
    expect(validateRequired('   ')).toBe(false);
  });
});

describe('validateURL', () => {
  it('should validate correct URLs', () => {
    expect(validateURL('https://example.com')).toBe(true);
    expect(validateURL('http://test.co.za')).toBe(true);
    expect(validateURL('https://sub.domain.com/path?query=value')).toBe(true);
  });

  it('should reject invalid URLs', () => {
    expect(validateURL('not a url')).toBe(false);
    expect(validateURL('example.com')).toBe(false); // Missing protocol
    expect(validateURL('')).toBe(false);
  });
});

describe('validateMinLength', () => {
  it('should validate strings meeting minimum length', () => {
    expect(validateMinLength('hello', 3)).toBe(true);
    expect(validateMinLength('test', 4)).toBe(true);
  });

  it('should reject strings below minimum length', () => {
    expect(validateMinLength('hi', 3)).toBe(false);
    expect(validateMinLength('', 1)).toBe(false);
  });
});

describe('validateMaxLength', () => {
  it('should validate strings within maximum length', () => {
    expect(validateMaxLength('hello', 10)).toBe(true);
    expect(validateMaxLength('test', 4)).toBe(true);
  });

  it('should reject strings exceeding maximum length', () => {
    expect(validateMaxLength('hello', 3)).toBe(false);
    expect(validateMaxLength('toolongstring', 5)).toBe(false);
  });
});
