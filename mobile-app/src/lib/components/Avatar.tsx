/**
 * Avatar Component
 * User avatar/profile picture display
 */

import React from 'react';
import { View, Text, Image, StyleSheet, TouchableOpacity } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts } from '@/lib/utils/theme';

export interface AvatarProps {
  source?: { uri: string } | number;
  name?: string;
  size?: 'small' | 'medium' | 'large' | 'xlarge';
  onPress?: () => void;
  showBadge?: boolean;
  badgeIcon?: string;
  editable?: boolean;
}

const Avatar: React.FC<AvatarProps> = ({
  source,
  name,
  size = 'medium',
  onPress,
  showBadge = false,
  badgeIcon = 'check',
  editable = false,
}) => {
  /**
   * Get avatar size
   */
  const getSize = (): number => {
    switch (size) {
      case 'small':
        return 40;
      case 'medium':
        return 64;
      case 'large':
        return 96;
      case 'xlarge':
        return 128;
      default:
        return 64;
    }
  };

  /**
   * Get font size for initials
   */
  const getFontSize = (): number => {
    switch (size) {
      case 'small':
        return 16;
      case 'medium':
        return 24;
      case 'large':
        return 36;
      case 'xlarge':
        return 48;
      default:
        return 24;
    }
  };

  /**
   * Get initials from name
   */
  const getInitials = (fullName?: string): string => {
    if (!fullName) return '?';

    const names = fullName.trim().split(' ');
    if (names.length === 1) {
      return names[0].charAt(0).toUpperCase();
    }

    return (names[0].charAt(0) + names[names.length - 1].charAt(0)).toUpperCase();
  };

  /**
   * Get background color based on name
   */
  const getBackgroundColor = (fullName?: string): string => {
    if (!fullName) return colors.textTertiary;

    const colors_list = [
      '#E74C3C', '#E67E22', '#F39C12', '#F1C40F',
      '#2ECC71', '#1ABC9C', '#3498DB', '#9B59B6',
      '#34495E', '#16A085', '#27AE60', '#2980B9',
      '#8E44AD', '#C0392B', '#D35400', '#2C3E50',
    ];

    const charCode = fullName.charCodeAt(0);
    return colors_list[charCode % colors_list.length];
  };

  const avatarSize = getSize();
  const fontSize = getFontSize();
  const initials = getInitials(name);
  const backgroundColor = source ? 'transparent' : getBackgroundColor(name);

  const AvatarContent = (
    <View
      style={[
        styles.container,
        {
          width: avatarSize,
          height: avatarSize,
          borderRadius: avatarSize / 2,
          backgroundColor,
        },
      ]}
    >
      {source ? (
        <Image
          source={source}
          style={[
            styles.image,
            {
              width: avatarSize,
              height: avatarSize,
              borderRadius: avatarSize / 2,
            },
          ]}
        />
      ) : (
        <Text style={[styles.initials, { fontSize }]}>{initials}</Text>
      )}

      {/* Badge */}
      {showBadge && (
        <View style={[styles.badge, { width: avatarSize * 0.3, height: avatarSize * 0.3 }]}>
          <Icon name={badgeIcon} size={avatarSize * 0.18} color="#FFFFFF" />
        </View>
      )}

      {/* Edit Icon */}
      {editable && (
        <View style={[styles.editBadge, { width: avatarSize * 0.3, height: avatarSize * 0.3 }]}>
          <Icon name="camera" size={avatarSize * 0.18} color="#FFFFFF" />
        </View>
      )}
    </View>
  );

  if (onPress || editable) {
    return (
      <TouchableOpacity onPress={onPress} activeOpacity={0.7}>
        {AvatarContent}
      </TouchableOpacity>
    );
  }

  return AvatarContent;
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
  },
  image: {
    resizeMode: 'cover',
  },
  initials: {
    color: '#FFFFFF',
    fontWeight: '700',
  },
  badge: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: colors.success,
    borderRadius: 100,
    borderWidth: 2,
    borderColor: '#FFFFFF',
    alignItems: 'center',
    justifyContent: 'center',
  },
  editBadge: {
    position: 'absolute',
    bottom: 0,
    right: 0,
    backgroundColor: colors.primary,
    borderRadius: 100,
    borderWidth: 2,
    borderColor: '#FFFFFF',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

export default Avatar;
