/**
 * Dark Theme Color Scheme
 * Complete dark mode colors for the Comply360 app
 */

export const darkColors = {
  // Primary brand colors
  primary: '#9f7aea',
  primaryLight: '#b794f6',
  primaryDark: '#805ad5',

  // Backgrounds
  background: '#1a1a1a',
  backgroundSecondary: '#242424',
  surface: '#2d2d2d',
  surfaceElevated: '#383838',

  // Text colors
  text: '#ffffff',
  textSecondary: '#b3b3b3',
  textTertiary: '#808080',
  textDisabled: '#4d4d4d',

  // Border colors
  border: '#404040',
  borderLight: '#333333',
  borderDark: '#4d4d4d',

  // Status colors
  success: '#10b981',
  successLight: '#34d399',
  successDark: '#059669',
  successBackground: '#064e3b',

  warning: '#f59e0b',
  warningLight: '#fbbf24',
  warningDark: '#d97706',
  warningBackground: '#78350f',

  error: '#ef4444',
  errorLight: '#f87171',
  errorDark: '#dc2626',
  errorBackground: '#7f1d1d',

  info: '#3b82f6',
  infoLight: '#60a5fa',
  infoDark: '#2563eb',
  infoBackground: '#1e3a8a',

  // Semantic colors
  pending: '#f59e0b',
  approved: '#10b981',
  rejected: '#ef4444',
  completed: '#3b82f6',
  inProgress: '#8b5cf6',
  draft: '#6b7280',

  // Chart colors (dark mode optimized)
  chart1: '#9f7aea',
  chart2: '#10b981',
  chart3: '#3b82f6',
  chart4: '#f59e0b',
  chart5: '#ef4444',
  chart6: '#8b5cf6',
  chart7: '#06b6d4',
  chart8: '#ec4899',

  // Component-specific
  cardBackground: '#2d2d2d',
  inputBackground: '#383838',
  inputBorder: '#4d4d4d',
  inputFocusBorder: '#9f7aea',

  // Shadows (dark mode shadows are lighter)
  shadowLight: 'rgba(255, 255, 255, 0.05)',
  shadowMedium: 'rgba(255, 255, 255, 0.1)',
  shadowDark: 'rgba(255, 255, 255, 0.15)',

  // Overlays
  overlay: 'rgba(0, 0, 0, 0.8)',
  overlayLight: 'rgba(0, 0, 0, 0.6)',

  // Icon colors
  iconPrimary: '#ffffff',
  iconSecondary: '#b3b3b3',
  iconTertiary: '#808080',
};

export const darkSpacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
};

export const darkFonts = {
  regular: 'System',
  medium: 'System',
  bold: 'System',
  sizes: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 18,
    xl: 20,
    xxl: 24,
    xxxl: 32,
  },
  weights: {
    regular: '400',
    medium: '500',
    semibold: '600',
    bold: '700',
  },
};

export const darkTheme = {
  colors: darkColors,
  spacing: darkSpacing,
  fonts: darkFonts,
  isDark: true,
};

export default darkTheme;
