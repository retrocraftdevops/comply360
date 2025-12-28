/**
 * Button Component
 * Reusable button with multiple variants and states
 */

import React from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
  ViewStyle,
  TextStyle,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';

export interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  icon?: string;
  iconPosition?: 'left' | 'right';
  fullWidth?: boolean;
  style?: ViewStyle;
  textStyle?: TextStyle;
}

const Button: React.FC<ButtonProps> = ({
  title,
  onPress,
  variant = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  icon,
  iconPosition = 'left',
  fullWidth = false,
  style,
  textStyle,
}) => {
  const getButtonStyle = (): ViewStyle[] => {
    const baseStyles: ViewStyle[] = [styles.button];

    // Size styles
    if (size === 'small') {
      baseStyles.push(styles.buttonSmall);
    } else if (size === 'large') {
      baseStyles.push(styles.buttonLarge);
    } else {
      baseStyles.push(styles.buttonMedium);
    }

    // Variant styles
    switch (variant) {
      case 'primary':
        baseStyles.push(styles.buttonPrimary);
        break;
      case 'secondary':
        baseStyles.push(styles.buttonSecondary);
        break;
      case 'outline':
        baseStyles.push(styles.buttonOutline);
        break;
      case 'ghost':
        baseStyles.push(styles.buttonGhost);
        break;
      case 'danger':
        baseStyles.push(styles.buttonDanger);
        break;
    }

    // State styles
    if (disabled || loading) {
      baseStyles.push(styles.buttonDisabled);
    }

    // Full width
    if (fullWidth) {
      baseStyles.push(styles.buttonFullWidth);
    }

    // Custom style
    if (style) {
      baseStyles.push(style);
    }

    return baseStyles;
  };

  const getTextStyle = (): TextStyle[] => {
    const baseStyles: TextStyle[] = [styles.text];

    // Size styles
    if (size === 'small') {
      baseStyles.push(styles.textSmall);
    } else if (size === 'large') {
      baseStyles.push(styles.textLarge);
    } else {
      baseStyles.push(styles.textMedium);
    }

    // Variant styles
    switch (variant) {
      case 'primary':
        baseStyles.push(styles.textPrimary);
        break;
      case 'secondary':
        baseStyles.push(styles.textSecondary);
        break;
      case 'outline':
        baseStyles.push(styles.textOutline);
        break;
      case 'ghost':
        baseStyles.push(styles.textGhost);
        break;
      case 'danger':
        baseStyles.push(styles.textDanger);
        break;
    }

    // State styles
    if (disabled || loading) {
      baseStyles.push(styles.textDisabled);
    }

    // Custom style
    if (textStyle) {
      baseStyles.push(textStyle);
    }

    return baseStyles;
  };

  const getIconColor = (): string => {
    if (disabled || loading) {
      return '#9ca3af';
    }

    switch (variant) {
      case 'primary':
      case 'danger':
        return '#FFFFFF';
      case 'secondary':
        return '#111827';
      case 'outline':
      case 'ghost':
        return '#7c3aed';
      default:
        return '#111827';
    }
  };

  const getIconSize = (): number => {
    switch (size) {
      case 'small':
        return 16;
      case 'large':
        return 24;
      default:
        return 20;
    }
  };

  return (
    <TouchableOpacity
      style={getButtonStyle()}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator
          size="small"
          color={variant === 'primary' || variant === 'danger' ? '#FFFFFF' : '#7c3aed'}
        />
      ) : (
        <>
          {icon && iconPosition === 'left' && (
            <Icon
              name={icon}
              size={getIconSize()}
              color={getIconColor()}
              style={styles.iconLeft}
            />
          )}
          <Text style={getTextStyle()}>{title}</Text>
          {icon && iconPosition === 'right' && (
            <Icon
              name={icon}
              size={getIconSize()}
              color={getIconColor()}
              style={styles.iconRight}
            />
          )}
        </>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 8,
  },
  // Size variants
  buttonSmall: {
    paddingVertical: 8,
    paddingHorizontal: 12,
  },
  buttonMedium: {
    paddingVertical: 12,
    paddingHorizontal: 20,
  },
  buttonLarge: {
    paddingVertical: 16,
    paddingHorizontal: 24,
  },
  // Color variants
  buttonPrimary: {
    backgroundColor: '#7c3aed',
  },
  buttonSecondary: {
    backgroundColor: '#f3f4f6',
  },
  buttonOutline: {
    backgroundColor: 'transparent',
    borderWidth: 2,
    borderColor: '#7c3aed',
  },
  buttonGhost: {
    backgroundColor: 'transparent',
  },
  buttonDanger: {
    backgroundColor: '#ef4444',
  },
  // State variants
  buttonDisabled: {
    opacity: 0.5,
  },
  buttonFullWidth: {
    width: '100%',
  },
  // Text styles
  text: {
    fontWeight: '600',
  },
  textSmall: {
    fontSize: 14,
  },
  textMedium: {
    fontSize: 16,
  },
  textLarge: {
    fontSize: 18,
  },
  textPrimary: {
    color: '#FFFFFF',
  },
  textSecondary: {
    color: '#111827',
  },
  textOutline: {
    color: '#7c3aed',
  },
  textGhost: {
    color: '#7c3aed',
  },
  textDanger: {
    color: '#FFFFFF',
  },
  textDisabled: {
    color: '#9ca3af',
  },
  // Icon styles
  iconLeft: {
    marginRight: 8,
  },
  iconRight: {
    marginLeft: 8,
  },
});

export default Button;
