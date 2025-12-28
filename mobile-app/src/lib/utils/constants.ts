/**
 * Application Constants
 * Centralized constants used throughout the app
 */

/**
 * API Configuration
 */
export const API = {
  BASE_URL: 'http://localhost:8080/api/v1',
  TIMEOUT: 30000, // 30 seconds
  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000, // 1 second
} as const;

/**
 * Storage Keys
 */
export const STORAGE_KEYS = {
  AUTH_TOKEN: '@comply360:auth_token',
  REFRESH_TOKEN: '@comply360:refresh_token',
  USER_DATA: '@comply360:user_data',
  BIOMETRIC_ENABLED: '@comply360:biometric_enabled',
  THEME: '@comply360:theme',
  LANGUAGE: '@comply360:language',
  ONBOARDING_COMPLETED: '@comply360:onboarding_completed',
} as const;

/**
 * Keychain Service Identifiers
 */
export const KEYCHAIN = {
  SERVICE_NAME: 'com.comply360.auth',
  BIOMETRIC_SERVICE: 'com.comply360.biometric',
} as const;

/**
 * Registration Types
 */
export const REGISTRATION_TYPES = {
  COMPANY: 'COMPANY',
  CLOSE_CORPORATION: 'CLOSE_CORPORATION',
  TRUST: 'TRUST',
  NPO: 'NPO',
} as const;

export type RegistrationType = typeof REGISTRATION_TYPES[keyof typeof REGISTRATION_TYPES];

/**
 * Registration Status
 */
export const REGISTRATION_STATUS = {
  DRAFT: 'DRAFT',
  PENDING: 'PENDING',
  IN_PROGRESS: 'IN_PROGRESS',
  COMPLETED: 'COMPLETED',
  REJECTED: 'REJECTED',
} as const;

export type RegistrationStatus = typeof REGISTRATION_STATUS[keyof typeof REGISTRATION_STATUS];

/**
 * Document Categories
 */
export const DOCUMENT_CATEGORIES = {
  CIPC: 'CIPC Documents',
  SARS: 'SARS Documents',
  BANKING: 'Banking Documents',
  IDENTITY: 'Identity Documents',
  PROOF_OF_ADDRESS: 'Proof of Address',
  SHAREHOLDING: 'Shareholding Documents',
  MEMORANDUM: 'Memorandum & Articles',
  RESOLUTION: 'Resolutions',
  OTHER: 'Other',
} as const;

/**
 * Document Status
 */
export const DOCUMENT_STATUS = {
  PENDING: 'PENDING',
  VERIFIED: 'VERIFIED',
  REJECTED: 'REJECTED',
} as const;

export type DocumentStatus = typeof DOCUMENT_STATUS[keyof typeof DOCUMENT_STATUS];

/**
 * Commission Status
 */
export const COMMISSION_STATUS = {
  PENDING: 'PENDING',
  APPROVED: 'APPROVED',
  PAID: 'PAID',
  DISPUTED: 'DISPUTED',
} as const;

export type CommissionStatus = typeof COMMISSION_STATUS[keyof typeof COMMISSION_STATUS];

/**
 * User Roles
 */
export const USER_ROLES = {
  ADMIN: 'admin',
  AGENT: 'agent',
  USER: 'user',
} as const;

export type UserRole = typeof USER_ROLES[keyof typeof USER_ROLES];

/**
 * Notification Types
 */
export const NOTIFICATION_TYPES = {
  INFO: 'info',
  SUCCESS: 'success',
  WARNING: 'warning',
  ERROR: 'error',
} as const;

export type NotificationType = typeof NOTIFICATION_TYPES[keyof typeof NOTIFICATION_TYPES];

/**
 * SADC Countries
 */
export const SADC_COUNTRIES = [
  { code: 'ZA', name: 'South Africa' },
  { code: 'BW', name: 'Botswana' },
  { code: 'ZW', name: 'Zimbabwe' },
  { code: 'NA', name: 'Namibia' },
  { code: 'MZ', name: 'Mozambique' },
  { code: 'ZM', name: 'Zambia' },
  { code: 'MW', name: 'Malawi' },
  { code: 'LS', name: 'Lesotho' },
  { code: 'SZ', name: 'Eswatini' },
  { code: 'AO', name: 'Angola' },
  { code: 'CD', name: 'Democratic Republic of Congo' },
  { code: 'MG', name: 'Madagascar' },
  { code: 'MU', name: 'Mauritius' },
  { code: 'SC', name: 'Seychelles' },
  { code: 'TZ', name: 'Tanzania' },
  { code: 'CD', name: 'Comoros' },
] as const;

/**
 * File Upload Limits
 */
export const FILE_UPLOAD = {
  MAX_SIZE_MB: 10,
  MAX_SIZE_BYTES: 10 * 1024 * 1024,
  ALLOWED_IMAGE_TYPES: ['jpg', 'jpeg', 'png', 'gif', 'webp'],
  ALLOWED_DOCUMENT_TYPES: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt'],
  ALLOWED_ALL_TYPES: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt'],
} as const;

/**
 * Pagination Defaults
 */
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 100,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
} as const;

/**
 * Date Formats
 */
export const DATE_FORMATS = {
  SHORT: 'YYYY-MM-DD',
  MEDIUM: 'MMM DD, YYYY',
  LONG: 'MMMM DD, YYYY',
  FULL: 'dddd, MMMM DD, YYYY',
  TIME: 'HH:mm',
  TIME_12H: 'h:mm A',
  DATETIME: 'YYYY-MM-DD HH:mm',
  DATETIME_FULL: 'MMMM DD, YYYY HH:mm:ss',
} as const;

/**
 * Currency
 */
export const CURRENCY = {
  CODE: 'ZAR',
  SYMBOL: 'R',
  NAME: 'South African Rand',
  DECIMALS: 2,
} as const;

/**
 * Validation Rules
 */
export const VALIDATION = {
  PASSWORD_MIN_LENGTH: 8,
  PASSWORD_MAX_LENGTH: 128,
  USERNAME_MIN_LENGTH: 3,
  USERNAME_MAX_LENGTH: 50,
  COMPANY_NAME_MIN_LENGTH: 2,
  COMPANY_NAME_MAX_LENGTH: 100,
  PHONE_MIN_LENGTH: 10,
  PHONE_MAX_LENGTH: 15,
} as const;

/**
 * Cache Durations (in milliseconds)
 */
export const CACHE_DURATION = {
  SHORT: 5 * 60 * 1000, // 5 minutes
  MEDIUM: 15 * 60 * 1000, // 15 minutes
  LONG: 60 * 60 * 1000, // 1 hour
  DAY: 24 * 60 * 60 * 1000, // 24 hours
} as const;

/**
 * Animation Durations
 */
export const ANIMATION = {
  FAST: 150,
  NORMAL: 300,
  SLOW: 500,
} as const;

/**
 * Debounce/Throttle Delays
 */
export const DELAYS = {
  SEARCH: 500,
  AUTO_SAVE: 2000,
  REFRESH: 1000,
} as const;

/**
 * Error Messages
 */
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNKNOWN_ERROR: 'An unexpected error occurred. Please try again.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  SESSION_EXPIRED: 'Your session has expired. Please log in again.',
  INVALID_CREDENTIALS: 'Invalid email or password.',
  SERVER_ERROR: 'Server error. Please try again later.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  NOT_FOUND: 'The requested resource was not found.',
  TIMEOUT: 'Request timed out. Please try again.',
} as const;

/**
 * Success Messages
 */
export const SUCCESS_MESSAGES = {
  LOGIN: 'Successfully logged in.',
  LOGOUT: 'Successfully logged out.',
  REGISTRATION_CREATED: 'Registration created successfully.',
  REGISTRATION_UPDATED: 'Registration updated successfully.',
  DOCUMENT_UPLOADED: 'Document uploaded successfully.',
  PROFILE_UPDATED: 'Profile updated successfully.',
  PASSWORD_CHANGED: 'Password changed successfully.',
  BIOMETRIC_ENABLED: 'Biometric authentication enabled.',
  BIOMETRIC_DISABLED: 'Biometric authentication disabled.',
} as const;

/**
 * App Information
 */
export const APP_INFO = {
  NAME: 'Comply360',
  FULL_NAME: 'Comply360 Mobile',
  DESCRIPTION: 'SADC Corporate Gateway Platform',
  VERSION: '1.0.0',
  BUILD_NUMBER: '1',
  COMPANY: 'Comply360',
  SUPPORT_EMAIL: 'support@comply360.com',
  SUPPORT_PHONE: '+27 11 123 4567',
  WEBSITE: 'https://comply360.com',
  TERMS_URL: 'https://comply360.com/terms',
  PRIVACY_URL: 'https://comply360.com/privacy',
} as const;

/**
 * Support Contact
 */
export const SUPPORT = {
  EMAIL: 'support@comply360.com',
  PHONE: '+27 11 123 4567',
  HOURS: 'Mon-Fri: 8:00 AM - 6:00 PM SAST',
  WEBSITE: 'https://comply360.com/support',
} as const;

/**
 * Feature Flags
 */
export const FEATURES = {
  BIOMETRIC_AUTH: true,
  PUSH_NOTIFICATIONS: true,
  OFFLINE_MODE: true,
  DOCUMENT_SCANNER: true,
  DARK_MODE: true,
  ANALYTICS: true,
} as const;

/**
 * Platform-specific Constants
 */
export const PLATFORM = {
  IS_IOS: require('react-native').Platform.OS === 'ios',
  IS_ANDROID: require('react-native').Platform.OS === 'android',
  OS: require('react-native').Platform.OS,
} as const;

/**
 * Regex Patterns
 */
export const REGEX = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE_ZA: /^(\+27|0)?[0-9]{9,10}$/,
  ID_NUMBER_ZA: /^\d{13}$/,
  REGISTRATION_NUMBER: /^\d{4}\/\d{6}\/\d{2}$/,
  TAX_NUMBER_ZA: /^\d{10}$/,
  VAT_NUMBER_ZA: /^4\d{9}$/,
  URL: /^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$/,
} as const;

/**
 * HTTP Status Codes
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
} as const;
