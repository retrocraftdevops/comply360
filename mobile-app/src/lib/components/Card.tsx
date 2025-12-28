/**
 * Card Component
 * Reusable card container with customizable styles
 */

import React, { ReactNode } from 'react';
import {
  View,
  StyleSheet,
  TouchableOpacity,
  ViewStyle,
} from 'react-native';

export interface CardProps {
  children: ReactNode;
  onPress?: () => void;
  variant?: 'default' | 'elevated' | 'outlined' | 'flat';
  padding?: 'none' | 'small' | 'medium' | 'large';
  style?: ViewStyle;
}

const Card: React.FC<CardProps> = ({
  children,
  onPress,
  variant = 'default',
  padding = 'medium',
  style,
}) => {
  const getCardStyle = (): ViewStyle[] => {
    const baseStyles: ViewStyle[] = [styles.card];

    // Variant styles
    switch (variant) {
      case 'elevated':
        baseStyles.push(styles.cardElevated);
        break;
      case 'outlined':
        baseStyles.push(styles.cardOutlined);
        break;
      case 'flat':
        baseStyles.push(styles.cardFlat);
        break;
      default:
        baseStyles.push(styles.cardDefault);
    }

    // Padding styles
    switch (padding) {
      case 'none':
        baseStyles.push(styles.paddingNone);
        break;
      case 'small':
        baseStyles.push(styles.paddingSmall);
        break;
      case 'large':
        baseStyles.push(styles.paddingLarge);
        break;
      default:
        baseStyles.push(styles.paddingMedium);
    }

    // Custom style
    if (style) {
      baseStyles.push(style);
    }

    return baseStyles;
  };

  if (onPress) {
    return (
      <TouchableOpacity
        style={getCardStyle()}
        onPress={onPress}
        activeOpacity={0.8}
      >
        {children}
      </TouchableOpacity>
    );
  }

  return <View style={getCardStyle()}>{children}</View>;
};

const styles = StyleSheet.create({
  card: {
    borderRadius: 12,
    overflow: 'hidden',
  },
  // Variant styles
  cardDefault: {
    backgroundColor: '#FFFFFF',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 3,
    elevation: 2,
  },
  cardElevated: {
    backgroundColor: '#FFFFFF',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.15,
    shadowRadius: 8,
    elevation: 5,
  },
  cardOutlined: {
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: '#e5e7eb',
  },
  cardFlat: {
    backgroundColor: '#f9fafb',
  },
  // Padding styles
  paddingNone: {
    padding: 0,
  },
  paddingSmall: {
    padding: 12,
  },
  paddingMedium: {
    padding: 16,
  },
  paddingLarge: {
    padding: 24,
  },
});

export default Card;
