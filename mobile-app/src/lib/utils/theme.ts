/**
 * Theme Configuration
 * Centralized colors, fonts, spacing, and other design tokens
 */

export const colors = {
  // Primary
  primary: '#7c3aed',
  primaryLight: '#a78bfa',
  primaryDark: '#5b21b6',
  primaryBg: '#f3e8ff',

  // Secondary
  secondary: '#3b82f6',
  secondaryLight: '#60a5fa',
  secondaryDark: '#2563eb',
  secondaryBg: '#dbeafe',

  // Success
  success: '#10b981',
  successLight: '#34d399',
  successDark: '#059669',
  successBg: '#d1fae5',

  // Warning
  warning: '#f59e0b',
  warningLight: '#fbbf24',
  warningDark: '#d97706',
  warningBg: '#fef3c7',

  // Error
  error: '#ef4444',
  errorLight: '#f87171',
  errorDark: '#dc2626',
  errorBg: '#fee2e2',

  // Info
  info: '#06b6d4',
  infoLight: '#22d3ee',
  infoDark: '#0891b2',
  infoBg: '#cffafe',

  // Neutral
  white: '#FFFFFF',
  black: '#000000',
  gray50: '#f9fafb',
  gray100: '#f3f4f6',
  gray200: '#e5e7eb',
  gray300: '#d1d5db',
  gray400: '#9ca3af',
  gray500: '#6b7280',
  gray600: '#4b5563',
  gray700: '#374151',
  gray800: '#1f2937',
  gray900: '#111827',

  // Background
  background: '#f9fafb',
  backgroundCard: '#FFFFFF',
  backgroundOverlay: 'rgba(0, 0, 0, 0.5)',

  // Border
  border: '#e5e7eb',
  borderLight: '#f3f4f6',
  borderDark: '#d1d5db',

  // Text
  textPrimary: '#111827',
  textSecondary: '#6b7280',
  textTertiary: '#9ca3af',
  textInverse: '#FFFFFF',
  textDisabled: '#d1d5db',
};

export const fonts = {
  // Font families
  regular: 'System',
  medium: 'System',
  semiBold: 'System',
  bold: 'System',

  // Font sizes
  xs: 12,
  sm: 14,
  base: 16,
  lg: 18,
  xl: 20,
  '2xl': 24,
  '3xl': 28,
  '4xl': 32,
  '5xl': 36,
  '6xl': 48,

  // Font weights
  weightRegular: '400' as const,
  weightMedium: '500' as const,
  weightSemiBold: '600' as const,
  weightBold: '700' as const,

  // Line heights
  lineHeightTight: 1.2,
  lineHeightNormal: 1.5,
  lineHeightRelaxed: 1.75,
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  base: 16,
  lg: 20,
  xl: 24,
  '2xl': 32,
  '3xl': 40,
  '4xl': 48,
  '5xl': 64,
  '6xl': 80,
};

export const borderRadius = {
  none: 0,
  sm: 4,
  base: 8,
  md: 12,
  lg: 16,
  xl: 24,
  full: 9999,
};

export const shadows = {
  sm: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  base: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  md: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 8,
    elevation: 4,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.2,
    shadowRadius: 16,
    elevation: 8,
  },
  xl: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 12 },
    shadowOpacity: 0.25,
    shadowRadius: 24,
    elevation: 12,
  },
};

export const opacity = {
  disabled: 0.5,
  hover: 0.8,
  overlay: 0.75,
};

export const zIndex = {
  base: 0,
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
};

export const breakpoints = {
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
};

/**
 * Theme object combining all design tokens
 */
export const theme = {
  colors,
  fonts,
  spacing,
  borderRadius,
  shadows,
  opacity,
  zIndex,
  breakpoints,
};

export type Theme = typeof theme;

export default theme;
