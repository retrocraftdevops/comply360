/**
 * SyncIndicator Component
 * Displays sync status and network connectivity
 */

import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Animated } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import networkStatus from '@/services/networkStatus';
import offlineQueue from '@/services/offlineQueue';
import { colors } from '@/lib/utils/theme';

export interface SyncIndicatorProps {
  onPress?: () => void;
  compact?: boolean;
}

const SyncIndicator: React.FC<SyncIndicatorProps> = ({ onPress, compact = false }) => {
  const [isOnline, setIsOnline] = useState(true);
  const [queueStatus, setQueueStatus] = useState({ total: 0, pending: 0, failed: 0 });
  const [isVisible, setIsVisible] = useState(false);
  const [fadeAnim] = useState(new Animated.Value(0));

  useEffect(() => {
    // Initialize
    setIsOnline(networkStatus.getIsOnline());
    updateQueueStatus();

    // Subscribe to network changes
    const unsubscribe = networkStatus.subscribe((online) => {
      setIsOnline(online);
      if (online) {
        // Auto-sync when coming online
        offlineQueue.process().then(updateQueueStatus);
      }
    });

    // Update queue status periodically
    const interval = setInterval(updateQueueStatus, 5000);

    return () => {
      unsubscribe();
      clearInterval(interval);
    };
  }, []);

  useEffect(() => {
    // Show/hide indicator based on status
    const shouldShow = !isOnline || queueStatus.pending > 0 || queueStatus.failed > 0;

    if (shouldShow !== isVisible) {
      setIsVisible(shouldShow);
      Animated.timing(fadeAnim, {
        toValue: shouldShow ? 1 : 0,
        duration: 300,
        useNativeDriver: true,
      }).start();
    }
  }, [isOnline, queueStatus]);

  const updateQueueStatus = () => {
    setQueueStatus(offlineQueue.getStatus());
  };

  const getStatusText = (): string => {
    if (!isOnline) {
      return 'Offline';
    }
    if (queueStatus.pending > 0) {
      return `Syncing ${queueStatus.pending} item${queueStatus.pending > 1 ? 's' : ''}`;
    }
    if (queueStatus.failed > 0) {
      return `${queueStatus.failed} failed`;
    }
    return 'Synced';
  };

  const getStatusIcon = (): string => {
    if (!isOnline) {
      return 'cloud-off-outline';
    }
    if (queueStatus.pending > 0) {
      return 'cloud-sync';
    }
    if (queueStatus.failed > 0) {
      return 'cloud-alert';
    }
    return 'cloud-check';
  };

  const getStatusColor = (): string => {
    if (!isOnline) {
      return colors.warning;
    }
    if (queueStatus.pending > 0) {
      return colors.info;
    }
    if (queueStatus.failed > 0) {
      return colors.error;
    }
    return colors.success;
  };

  const handlePress = () => {
    if (onPress) {
      onPress();
    } else if (!isOnline && queueStatus.pending > 0) {
      // Default behavior: retry sync
      offlineQueue.process().then(updateQueueStatus);
    } else if (queueStatus.failed > 0) {
      // Retry failed items
      offlineQueue.retryFailed().then(updateQueueStatus);
    }
  };

  if (!isVisible) {
    return null;
  }

  const statusColor = getStatusColor();

  if (compact) {
    return (
      <Animated.View style={[styles.compactContainer, { opacity: fadeAnim }]}>
        <TouchableOpacity
          style={[styles.compactButton, { backgroundColor: `${statusColor}15` }]}
          onPress={handlePress}
          activeOpacity={0.7}
        >
          <Icon name={getStatusIcon()} size={16} color={statusColor} />
        </TouchableOpacity>
      </Animated.View>
    );
  }

  return (
    <Animated.View style={[styles.container, { opacity: fadeAnim }]}>
      <TouchableOpacity
        style={[styles.button, { backgroundColor: `${statusColor}15`, borderColor: statusColor }]}
        onPress={handlePress}
        activeOpacity={0.7}
      >
        <Icon name={getStatusIcon()} size={20} color={statusColor} />
        <Text style={[styles.text, { color: statusColor }]}>{getStatusText()}</Text>
        {queueStatus.pending > 0 && (
          <View style={[styles.badge, { backgroundColor: statusColor }]}>
            <Text style={styles.badgeText}>{queueStatus.pending}</Text>
          </View>
        )}
      </TouchableOpacity>
    </Animated.View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginVertical: 8,
  },
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 10,
    paddingHorizontal: 16,
    borderRadius: 8,
    borderWidth: 1,
  },
  text: {
    fontSize: 14,
    fontWeight: '600',
    marginLeft: 8,
  },
  badge: {
    marginLeft: 8,
    borderRadius: 10,
    paddingHorizontal: 8,
    paddingVertical: 2,
  },
  badgeText: {
    fontSize: 12,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  compactContainer: {
    position: 'absolute',
    top: 8,
    right: 8,
    zIndex: 1000,
  },
  compactButton: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
});

export default SyncIndicator;
