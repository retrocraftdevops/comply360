/**
 * Formatting Utilities
 * Common formatting functions for display
 */

/**
 * Format currency in South African Rand
 */
export const formatCurrency = (
  amount: number,
  options?: {
    showSymbol?: boolean;
    decimals?: number;
    locale?: string;
  }
): string => {
  const {
    showSymbol = true,
    decimals = 2,
    locale = 'en-ZA',
  } = options || {};

  const formatted = amount.toLocaleString(locale, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });

  return showSymbol ? `R ${formatted}` : formatted;
};

/**
 * Format percentage
 */
export const formatPercentage = (
  value: number,
  options?: {
    decimals?: number;
    showSymbol?: boolean;
  }
): string => {
  const { decimals = 1, showSymbol = true } = options || {};

  const formatted = value.toFixed(decimals);
  return showSymbol ? `${formatted}%` : formatted;
};

/**
 * Format date
 */
export const formatDate = (
  date: string | Date,
  format: 'short' | 'medium' | 'long' | 'full' = 'medium'
): string => {
  const dateObj = typeof date === 'string' ? new Date(date) : date;

  const options: Intl.DateTimeFormatOptions = {
    short: { year: 'numeric', month: '2-digit', day: '2-digit' },
    medium: { year: 'numeric', month: 'short', day: 'numeric' },
    long: { year: 'numeric', month: 'long', day: 'numeric' },
    full: { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' },
  }[format];

  return dateObj.toLocaleDateString('en-ZA', options);
};

/**
 * Format time
 */
export const formatTime = (
  date: string | Date,
  options?: {
    includeSeconds?: boolean;
    use24Hour?: boolean;
  }
): string => {
  const { includeSeconds = false, use24Hour = false } = options || {};

  const dateObj = typeof date === 'string' ? new Date(date) : date;

  const formatOptions: Intl.DateTimeFormatOptions = {
    hour: '2-digit',
    minute: '2-digit',
    ...(includeSeconds && { second: '2-digit' }),
    hour12: !use24Hour,
  };

  return dateObj.toLocaleTimeString('en-ZA', formatOptions);
};

/**
 * Format date and time
 */
export const formatDateTime = (
  date: string | Date,
  options?: {
    dateFormat?: 'short' | 'medium' | 'long';
    includeSeconds?: boolean;
  }
): string => {
  const { dateFormat = 'medium', includeSeconds = false } = options || {};

  const dateStr = formatDate(date, dateFormat);
  const timeStr = formatTime(date, { includeSeconds });

  return `${dateStr} ${timeStr}`;
};

/**
 * Format relative time (e.g., "2 hours ago", "in 3 days")
 */
export const formatRelativeTime = (date: string | Date): string => {
  const dateObj = typeof date === 'string' ? new Date(date) : date;
  const now = new Date();
  const diffMs = now.getTime() - dateObj.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);
  const diffWeek = Math.floor(diffDay / 7);
  const diffMonth = Math.floor(diffDay / 30);
  const diffYear = Math.floor(diffDay / 365);

  if (diffSec < 60) {
    return 'just now';
  } else if (diffMin < 60) {
    return `${diffMin} minute${diffMin > 1 ? 's' : ''} ago`;
  } else if (diffHour < 24) {
    return `${diffHour} hour${diffHour > 1 ? 's' : ''} ago`;
  } else if (diffDay < 7) {
    return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`;
  } else if (diffWeek < 4) {
    return `${diffWeek} week${diffWeek > 1 ? 's' : ''} ago`;
  } else if (diffMonth < 12) {
    return `${diffMonth} month${diffMonth > 1 ? 's' : ''} ago`;
  } else {
    return `${diffYear} year${diffYear > 1 ? 's' : ''} ago`;
  }
};

/**
 * Format phone number (South African)
 */
export const formatPhoneNumber = (phone: string): string => {
  const cleaned = phone.replace(/\D/g, '');

  // South African format: +27 XX XXX XXXX
  if (cleaned.startsWith('27')) {
    return `+${cleaned.slice(0, 2)} ${cleaned.slice(2, 4)} ${cleaned.slice(4, 7)} ${cleaned.slice(7)}`;
  }

  // Local format: 0XX XXX XXXX
  if (cleaned.startsWith('0')) {
    return `${cleaned.slice(0, 3)} ${cleaned.slice(3, 6)} ${cleaned.slice(6)}`;
  }

  return phone;
};

/**
 * Format ID number (South African)
 */
export const formatIDNumber = (idNumber: string): string => {
  const cleaned = idNumber.replace(/\D/g, '');

  if (cleaned.length === 13) {
    return `${cleaned.slice(0, 6)}-${cleaned.slice(6, 10)}-${cleaned.slice(10, 11)}-${cleaned.slice(11, 12)}-${cleaned.slice(12)}`;
  }

  return idNumber;
};

/**
 * Format file size
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};

/**
 * Format number with thousand separators
 */
export const formatNumber = (
  number: number,
  options?: {
    decimals?: number;
    locale?: string;
  }
): string => {
  const { decimals = 0, locale = 'en-ZA' } = options || {};

  return number.toLocaleString(locale, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
};

/**
 * Format duration (milliseconds to readable format)
 */
export const formatDuration = (milliseconds: number): string => {
  if (milliseconds < 1000) {
    return `${milliseconds}ms`;
  }

  const seconds = Math.floor(milliseconds / 1000);

  if (seconds < 60) {
    return `${seconds}s`;
  }

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;

  if (minutes < 60) {
    return remainingSeconds > 0
      ? `${minutes}m ${remainingSeconds}s`
      : `${minutes}m`;
  }

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;

  return remainingMinutes > 0
    ? `${hours}h ${remainingMinutes}m`
    : `${hours}h`;
};

/**
 * Truncate text with ellipsis
 */
export const truncateText = (
  text: string,
  maxLength: number,
  options?: {
    ellipsis?: string;
    breakWords?: boolean;
  }
): string => {
  const { ellipsis = '...', breakWords = false } = options || {};

  if (text.length <= maxLength) {
    return text;
  }

  const truncated = text.substring(0, maxLength - ellipsis.length);

  if (breakWords) {
    return truncated + ellipsis;
  }

  // Try to break at last word boundary
  const lastSpace = truncated.lastIndexOf(' ');
  if (lastSpace > 0) {
    return truncated.substring(0, lastSpace) + ellipsis;
  }

  return truncated + ellipsis;
};

/**
 * Capitalize first letter
 */
export const capitalize = (text: string): string => {
  if (!text) return text;
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase();
};

/**
 * Title case (capitalize each word)
 */
export const titleCase = (text: string): string => {
  return text
    .split(' ')
    .map((word) => capitalize(word))
    .join(' ');
};

/**
 * Format initials from name
 */
export const formatInitials = (
  name: string,
  options?: {
    maxInitials?: number;
  }
): string => {
  const { maxInitials = 2 } = options || {};

  const words = name.trim().split(/\s+/);
  const initials = words
    .slice(0, maxInitials)
    .map((word) => word.charAt(0).toUpperCase())
    .join('');

  return initials;
};

/**
 * Format credit card number
 */
export const formatCreditCard = (cardNumber: string): string => {
  const cleaned = cardNumber.replace(/\s/g, '');
  const groups = cleaned.match(/.{1,4}/g) || [];
  return groups.join(' ');
};

/**
 * Mask sensitive data (e.g., card numbers, IDs)
 */
export const maskSensitiveData = (
  data: string,
  options?: {
    visibleStart?: number;
    visibleEnd?: number;
    maskChar?: string;
  }
): string => {
  const { visibleStart = 0, visibleEnd = 4, maskChar = '*' } = options || {};

  if (data.length <= visibleStart + visibleEnd) {
    return data;
  }

  const start = data.substring(0, visibleStart);
  const end = data.substring(data.length - visibleEnd);
  const masked = maskChar.repeat(data.length - visibleStart - visibleEnd);

  return start + masked + end;
};

/**
 * Format status badge text
 */
export const formatStatus = (status: string): string => {
  return status
    .split('_')
    .map((word) => capitalize(word))
    .join(' ');
};

/**
 * Format array to comma-separated string
 */
export const formatList = (
  items: string[],
  options?: {
    conjunction?: 'and' | 'or';
    maxItems?: number;
  }
): string => {
  const { conjunction = 'and', maxItems } = options || {};

  const displayItems = maxItems ? items.slice(0, maxItems) : items;
  const remaining = items.length - displayItems.length;

  if (displayItems.length === 0) {
    return '';
  }

  if (displayItems.length === 1) {
    return displayItems[0];
  }

  const allButLast = displayItems.slice(0, -1).join(', ');
  const last = displayItems[displayItems.length - 1];
  const formatted = `${allButLast} ${conjunction} ${last}`;

  if (remaining > 0) {
    return `${formatted} and ${remaining} more`;
  }

  return formatted;
};
