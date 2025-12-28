/**
 * LoadingSkeleton Component
 * Animated skeleton loaders for better perceived performance
 */

import React, { useEffect, useRef } from 'react';
import { View, StyleSheet, Animated } from 'react-native';

export interface LoadingSkeletonProps {
  width?: number | string;
  height?: number;
  borderRadius?: number;
  style?: any;
}

const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({
  width = '100%',
  height = 16,
  borderRadius = 4,
  style,
}) => {
  const animatedValue = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    const animation = Animated.loop(
      Animated.sequence([
        Animated.timing(animatedValue, {
          toValue: 1,
          duration: 1000,
          useNativeDriver: true,
        }),
        Animated.timing(animatedValue, {
          toValue: 0,
          duration: 1000,
          useNativeDriver: true,
        }),
      ])
    );

    animation.start();

    return () => animation.stop();
  }, []);

  const opacity = animatedValue.interpolate({
    inputRange: [0, 1],
    outputRange: [0.3, 0.7],
  });

  return (
    <Animated.View
      style={[
        styles.skeleton,
        {
          width,
          height,
          borderRadius,
          opacity,
        },
        style,
      ]}
    />
  );
};

/**
 * Card Skeleton
 * Pre-built skeleton for card layouts
 */
export const CardSkeleton: React.FC = () => (
  <View style={styles.cardContainer}>
    <View style={styles.cardHeader}>
      <LoadingSkeleton width={48} height={48} borderRadius={24} />
      <View style={styles.cardHeaderText}>
        <LoadingSkeleton width="60%" height={16} />
        <LoadingSkeleton width="40%" height={12} style={{ marginTop: 8 }} />
      </View>
    </View>
    <LoadingSkeleton width="100%" height={12} style={{ marginTop: 12 }} />
    <LoadingSkeleton width="80%" height={12} style={{ marginTop: 8 }} />
    <View style={styles.cardFooter}>
      <LoadingSkeleton width={60} height={24} borderRadius={12} />
      <LoadingSkeleton width={80} height={12} />
    </View>
  </View>
);

/**
 * List Item Skeleton
 * Pre-built skeleton for list items
 */
export const ListItemSkeleton: React.FC = () => (
  <View style={styles.listItemContainer}>
    <LoadingSkeleton width={48} height={48} borderRadius={8} />
    <View style={styles.listItemContent}>
      <LoadingSkeleton width="70%" height={16} />
      <LoadingSkeleton width="50%" height={12} style={{ marginTop: 8 }} />
    </View>
    <LoadingSkeleton width={24} height={24} borderRadius={12} />
  </View>
);

/**
 * Avatar Skeleton
 * Pre-built skeleton for avatars
 */
export const AvatarSkeleton: React.FC<{ size?: number }> = ({ size = 48 }) => (
  <LoadingSkeleton width={size} height={size} borderRadius={size / 2} />
);

/**
 * Text Line Skeleton
 * Pre-built skeleton for text lines
 */
export const TextLineSkeleton: React.FC<{ width?: number | string }> = ({ width = '100%' }) => (
  <LoadingSkeleton width={width} height={14} style={{ marginVertical: 4 }} />
);

/**
 * Image Skeleton
 * Pre-built skeleton for images
 */
export const ImageSkeleton: React.FC<{
  width?: number | string;
  height?: number;
  borderRadius?: number;
}> = ({ width = '100%', height = 200, borderRadius = 8 }) => (
  <LoadingSkeleton width={width} height={height} borderRadius={borderRadius} />
);

/**
 * Registration Card Skeleton
 * Specific skeleton for registration cards
 */
export const RegistrationCardSkeleton: React.FC = () => (
  <View style={styles.registrationCard}>
    <View style={styles.registrationHeader}>
      <LoadingSkeleton width={40} height={40} borderRadius={20} />
      <View style={styles.registrationHeaderText}>
        <LoadingSkeleton width="70%" height={16} />
        <LoadingSkeleton width="50%" height={12} style={{ marginTop: 6 }} />
      </View>
    </View>
    <View style={styles.registrationBody}>
      <LoadingSkeleton width="40%" height={12} />
      <LoadingSkeleton width="60%" height={12} style={{ marginTop: 8 }} />
    </View>
    <View style={styles.registrationFooter}>
      <LoadingSkeleton width={80} height={24} borderRadius={12} />
      <LoadingSkeleton width={60} height={12} />
    </View>
  </View>
);

/**
 * Commission Card Skeleton
 * Specific skeleton for commission cards
 */
export const CommissionCardSkeleton: React.FC = () => (
  <View style={styles.commissionCard}>
    <View style={styles.commissionHeader}>
      <LoadingSkeleton width="60%" height={18} />
      <LoadingSkeleton width={60} height={20} borderRadius={10} />
    </View>
    <LoadingSkeleton width="80%" height={14} style={{ marginTop: 8 }} />
    <View style={styles.commissionFooter}>
      <View style={styles.commissionAmount}>
        <LoadingSkeleton width={100} height={24} />
        <LoadingSkeleton width={60} height={12} style={{ marginTop: 4 }} />
      </View>
    </View>
  </View>
);

/**
 * Document Card Skeleton
 * Specific skeleton for document cards
 */
export const DocumentCardSkeleton: React.FC = () => (
  <View style={styles.documentCard}>
    <LoadingSkeleton width={48} height={48} borderRadius={8} />
    <View style={styles.documentContent}>
      <LoadingSkeleton width="70%" height={16} />
      <LoadingSkeleton width="50%" height={12} style={{ marginTop: 6 }} />
      <LoadingSkeleton width="30%" height={10} style={{ marginTop: 6 }} />
    </View>
    <View style={styles.documentActions}>
      <LoadingSkeleton width={60} height={20} borderRadius={10} />
      <LoadingSkeleton width={32} height={32} borderRadius={16} style={{ marginTop: 8 }} />
    </View>
  </View>
);

/**
 * Profile Header Skeleton
 * Specific skeleton for profile header
 */
export const ProfileHeaderSkeleton: React.FC = () => (
  <View style={styles.profileHeader}>
    <LoadingSkeleton width={128} height={128} borderRadius={64} />
    <LoadingSkeleton width={200} height={24} style={{ marginTop: 16 }} />
    <LoadingSkeleton width={150} height={16} style={{ marginTop: 8 }} />
    <View style={styles.profileActions}>
      <LoadingSkeleton width={120} height={40} borderRadius={8} />
      <LoadingSkeleton width={100} height={40} borderRadius={8} style={{ marginLeft: 12 }} />
    </View>
  </View>
);

const styles = StyleSheet.create({
  skeleton: {
    backgroundColor: '#e5e7eb',
  },
  cardContainer: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
  },
  cardHeader: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  cardHeaderText: {
    flex: 1,
    marginLeft: 12,
  },
  cardFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 16,
  },
  listItemContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#FFFFFF',
    borderRadius: 8,
    padding: 12,
    marginBottom: 8,
  },
  listItemContent: {
    flex: 1,
    marginLeft: 12,
  },
  registrationCard: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
  },
  registrationHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
  },
  registrationHeaderText: {
    flex: 1,
    marginLeft: 12,
  },
  registrationBody: {
    marginBottom: 12,
  },
  registrationFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  commissionCard: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
  },
  commissionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  commissionFooter: {
    marginTop: 12,
  },
  commissionAmount: {
    alignItems: 'flex-start',
  },
  documentCard: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    flexDirection: 'row',
    alignItems: 'center',
  },
  documentContent: {
    flex: 1,
    marginLeft: 12,
  },
  documentActions: {
    alignItems: 'center',
  },
  profileHeader: {
    alignItems: 'center',
    backgroundColor: '#FFFFFF',
    paddingTop: 32,
    paddingBottom: 20,
    paddingHorizontal: 20,
  },
  profileActions: {
    flexDirection: 'row',
    marginTop: 20,
  },
});

export default LoadingSkeleton;
