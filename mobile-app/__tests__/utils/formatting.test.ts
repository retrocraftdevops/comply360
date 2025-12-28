/**
 * Formatting Utilities Tests
 */

import {
  formatCurrency,
  formatDate,
  formatRelativeTime,
  formatPhoneNumber,
  formatFileSize,
  formatPercentage,
  capitalizeFirstLetter,
  truncateText,
} from '@/lib/utils/formatting';

describe('formatCurrency', () => {
  it('should format currency with symbol', () => {
    expect(formatCurrency(1234.56)).toBe('R 1,234.56');
    expect(formatCurrency(1000)).toBe('R 1,000.00');
    expect(formatCurrency(0)).toBe('R 0.00');
  });

  it('should format currency without symbol', () => {
    expect(formatCurrency(1234.56, { showSymbol: false })).toBe('1,234.56');
  });

  it('should format with custom decimals', () => {
    expect(formatCurrency(1234.567, { decimals: 0 })).toBe('R 1,235');
    expect(formatCurrency(1234.567, { decimals: 3 })).toBe('R 1,234.567');
  });

  it('should handle negative values', () => {
    expect(formatCurrency(-1234.56)).toBe('R -1,234.56');
  });
});

describe('formatDate', () => {
  const testDate = new Date('2025-12-28T14:30:00Z');

  it('should format date in short format', () => {
    const result = formatDate(testDate, 'short');
    expect(result).toContain('Dec');
    expect(result).toContain('28');
  });

  it('should format date in medium format', () => {
    const result = formatDate(testDate, 'medium');
    expect(result).toContain('Dec');
    expect(result).toContain('28');
    expect(result).toContain('2025');
  });

  it('should format date in long format', () => {
    const result = formatDate(testDate, 'long');
    expect(result).toContain('December');
    expect(result).toContain('28');
  });

  it('should handle string dates', () => {
    const result = formatDate('2025-12-28', 'short');
    expect(result).toBeTruthy();
  });
});

describe('formatRelativeTime', () => {
  it('should format recent times', () => {
    const now = new Date();
    const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000);
    const result = formatRelativeTime(fiveMinutesAgo);
    expect(result).toContain('minute');
  });

  it('should format hours ago', () => {
    const now = new Date();
    const twoHoursAgo = new Date(now.getTime() - 2 * 60 * 60 * 1000);
    const result = formatRelativeTime(twoHoursAgo);
    expect(result).toContain('hour');
  });

  it('should format days ago', () => {
    const now = new Date();
    const threeDaysAgo = new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000);
    const result = formatRelativeTime(threeDaysAgo);
    expect(result).toContain('day');
  });

  it('should handle string dates', () => {
    const result = formatRelativeTime('2025-12-28');
    expect(result).toBeTruthy();
  });
});

describe('formatPhoneNumber', () => {
  it('should format South African phone numbers', () => {
    expect(formatPhoneNumber('0123456789')).toBe('012 345 6789');
    expect(formatPhoneNumber('+27123456789')).toBe('+27 12 345 6789');
  });

  it('should handle phone numbers with spaces', () => {
    expect(formatPhoneNumber('012 345 6789')).toBe('012 345 6789');
  });

  it('should return original if invalid', () => {
    expect(formatPhoneNumber('123')).toBe('123');
  });
});

describe('formatFileSize', () => {
  it('should format bytes', () => {
    expect(formatFileSize(500)).toBe('500 B');
  });

  it('should format kilobytes', () => {
    expect(formatFileSize(1536)).toBe('1.5 KB');
    expect(formatFileSize(2048)).toBe('2.0 KB');
  });

  it('should format megabytes', () => {
    expect(formatFileSize(1572864)).toBe('1.5 MB'); // 1.5 * 1024 * 1024
    expect(formatFileSize(5242880)).toBe('5.0 MB'); // 5 * 1024 * 1024
  });

  it('should format gigabytes', () => {
    expect(formatFileSize(1610612736)).toBe('1.5 GB'); // 1.5 * 1024 * 1024 * 1024
  });

  it('should handle zero', () => {
    expect(formatFileSize(0)).toBe('0 B');
  });
});

describe('formatPercentage', () => {
  it('should format percentages', () => {
    expect(formatPercentage(0.5)).toBe('50%');
    expect(formatPercentage(0.75)).toBe('75%');
    expect(formatPercentage(1)).toBe('100%');
  });

  it('should handle decimals', () => {
    expect(formatPercentage(0.333, 1)).toBe('33.3%');
    expect(formatPercentage(0.6667, 2)).toBe('66.67%');
  });

  it('should handle zero', () => {
    expect(formatPercentage(0)).toBe('0%');
  });
});

describe('capitalizeFirstLetter', () => {
  it('should capitalize first letter', () => {
    expect(capitalizeFirstLetter('hello')).toBe('Hello');
    expect(capitalizeFirstLetter('world')).toBe('World');
  });

  it('should handle already capitalized', () => {
    expect(capitalizeFirstLetter('Hello')).toBe('Hello');
  });

  it('should handle single letter', () => {
    expect(capitalizeFirstLetter('a')).toBe('A');
  });

  it('should handle empty string', () => {
    expect(capitalizeFirstLetter('')).toBe('');
  });
});

describe('truncateText', () => {
  it('should truncate long text', () => {
    expect(truncateText('This is a very long text', 10)).toBe('This is a...');
  });

  it('should not truncate short text', () => {
    expect(truncateText('Short', 10)).toBe('Short');
  });

  it('should handle custom ellipsis', () => {
    expect(truncateText('Long text here', 9, '…')).toBe('Long text…');
  });

  it('should handle exact length', () => {
    expect(truncateText('Exact', 5)).toBe('Exact');
  });
});
