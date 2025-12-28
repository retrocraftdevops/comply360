/**
 * Validation Utilities
 * Common validation functions for forms and inputs
 */

/**
 * Email validation
 * Uses RFC 5322 compliant regex
 */
export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email.trim());
};

/**
 * Password validation
 * At least 8 characters, 1 uppercase, 1 lowercase, 1 number
 */
export const validatePassword = (password: string): {
  isValid: boolean;
  errors: string[];
} => {
  const errors: string[] = [];

  if (password.length < 8) {
    errors.push('Password must be at least 8 characters');
  }

  if (!/[A-Z]/.test(password)) {
    errors.push('Password must contain at least one uppercase letter');
  }

  if (!/[a-z]/.test(password)) {
    errors.push('Password must contain at least one lowercase letter');
  }

  if (!/[0-9]/.test(password)) {
    errors.push('Password must contain at least one number');
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

/**
 * Phone number validation (South African format)
 * Accepts: +27, 0, or plain 10-digit numbers
 */
export const validatePhoneNumber = (phone: string): boolean => {
  const phoneRegex = /^(\+27|0)?[0-9]{9,10}$/;
  return phoneRegex.test(phone.replace(/\s/g, ''));
};

/**
 * ID Number validation (South African)
 * Basic format check: YYMMDD-SSSS-C-A-Z
 */
export const validateIDNumber = (idNumber: string): boolean => {
  // Remove spaces and dashes
  const cleaned = idNumber.replace(/[\s-]/g, '');

  // Must be 13 digits
  if (!/^\d{13}$/.test(cleaned)) {
    return false;
  }

  // Validate date portion (YYMMDD)
  const year = parseInt(cleaned.substring(0, 2), 10);
  const month = parseInt(cleaned.substring(2, 4), 10);
  const day = parseInt(cleaned.substring(4, 6), 10);

  if (month < 1 || month > 12) {
    return false;
  }

  if (day < 1 || day > 31) {
    return false;
  }

  return true;
};

/**
 * Registration number validation (CIPC format)
 * Format: YYYY/NNNNNN/07 or similar
 */
export const validateRegistrationNumber = (regNumber: string): boolean => {
  const regNumberRegex = /^\d{4}\/\d{6}\/\d{2}$/;
  return regNumberRegex.test(regNumber);
};

/**
 * Tax number validation (South African)
 * Format: 10 digits
 */
export const validateTaxNumber = (taxNumber: string): boolean => {
  const cleaned = taxNumber.replace(/\s/g, '');
  return /^\d{10}$/.test(cleaned);
};

/**
 * VAT number validation (South African)
 * Format: 4XXXXXXXXX (starts with 4, then 9 digits)
 */
export const validateVATNumber = (vatNumber: string): boolean => {
  const cleaned = vatNumber.replace(/\s/g, '');
  return /^4\d{9}$/.test(cleaned);
};

/**
 * Company name validation
 * At least 2 characters, alphanumeric and some special chars
 */
export const validateCompanyName = (name: string): boolean => {
  if (name.trim().length < 2) {
    return false;
  }

  // Allow letters, numbers, spaces, hyphens, apostrophes, parentheses, ampersands
  const nameRegex = /^[a-zA-Z0-9\s\-'()&.]+$/;
  return nameRegex.test(name);
};

/**
 * Amount validation
 * Must be a positive number
 */
export const validateAmount = (amount: string | number): boolean => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  return !isNaN(num) && num > 0;
};

/**
 * Percentage validation
 * Must be between 0 and 100
 */
export const validatePercentage = (percentage: string | number): boolean => {
  const num = typeof percentage === 'string' ? parseFloat(percentage) : percentage;
  return !isNaN(num) && num >= 0 && num <= 100;
};

/**
 * Required field validation
 */
export const validateRequired = (value: any): boolean => {
  if (value === null || value === undefined) {
    return false;
  }

  if (typeof value === 'string') {
    return value.trim().length > 0;
  }

  if (Array.isArray(value)) {
    return value.length > 0;
  }

  return true;
};

/**
 * Min length validation
 */
export const validateMinLength = (value: string, minLength: number): boolean => {
  return value.trim().length >= minLength;
};

/**
 * Max length validation
 */
export const validateMaxLength = (value: string, maxLength: number): boolean => {
  return value.trim().length <= maxLength;
};

/**
 * URL validation
 */
export const validateURL = (url: string): boolean => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

/**
 * Date validation
 * Checks if date is valid and in the past or future
 */
export const validateDate = (
  date: string | Date,
  options?: {
    allowPast?: boolean;
    allowFuture?: boolean;
    minDate?: Date;
    maxDate?: Date;
  }
): boolean => {
  const dateObj = typeof date === 'string' ? new Date(date) : date;

  // Check if valid date
  if (isNaN(dateObj.getTime())) {
    return false;
  }

  const now = new Date();

  // Check past/future restrictions
  if (options?.allowPast === false && dateObj < now) {
    return false;
  }

  if (options?.allowFuture === false && dateObj > now) {
    return false;
  }

  // Check min/max date
  if (options?.minDate && dateObj < options.minDate) {
    return false;
  }

  if (options?.maxDate && dateObj > options.maxDate) {
    return false;
  }

  return true;
};

/**
 * Credit card validation (basic Luhn algorithm)
 */
export const validateCreditCard = (cardNumber: string): boolean => {
  const cleaned = cardNumber.replace(/\s/g, '');

  if (!/^\d{13,19}$/.test(cleaned)) {
    return false;
  }

  let sum = 0;
  let isEven = false;

  for (let i = cleaned.length - 1; i >= 0; i--) {
    let digit = parseInt(cleaned.charAt(i), 10);

    if (isEven) {
      digit *= 2;
      if (digit > 9) {
        digit -= 9;
      }
    }

    sum += digit;
    isEven = !isEven;
  }

  return sum % 10 === 0;
};

/**
 * File size validation
 * Checks if file size is within limits
 */
export const validateFileSize = (
  sizeInBytes: number,
  maxSizeInMB: number
): boolean => {
  const maxSizeInBytes = maxSizeInMB * 1024 * 1024;
  return sizeInBytes <= maxSizeInBytes;
};

/**
 * File type validation
 * Checks if file extension is in allowed list
 */
export const validateFileType = (
  fileName: string,
  allowedExtensions: string[]
): boolean => {
  const extension = fileName.split('.').pop()?.toLowerCase();
  return extension
    ? allowedExtensions.map((ext) => ext.toLowerCase()).includes(extension)
    : false;
};

/**
 * Combined form validation
 * Validates multiple fields and returns errors
 */
export const validateForm = <T extends Record<string, any>>(
  data: T,
  rules: Partial<Record<keyof T, (value: any) => string | null>>
): { isValid: boolean; errors: Partial<Record<keyof T, string>> } => {
  const errors: Partial<Record<keyof T, string>> = {};

  for (const field in rules) {
    const validator = rules[field];
    if (validator) {
      const error = validator(data[field]);
      if (error) {
        errors[field] = error;
      }
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
};
