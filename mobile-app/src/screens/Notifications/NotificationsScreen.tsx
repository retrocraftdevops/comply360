/**
 * NotificationsScreen
 * Display and manage user notifications
 */

import React, { useState, useMemo } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  RefreshControl,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import NotificationCard, { Notification } from '@/lib/components/NotificationCard';
import { EmptyState, BottomSheet } from '@/lib/components';
import { useAppDispatch } from '@/store/store';
import { showToast } from '@/store/slices/uiSlice';
import { colors } from '@/lib/utils/theme';

const NotificationsScreen: React.FC = () => {
  const dispatch = useAppDispatch();

  const [filterType, setFilterType] = useState<'ALL' | 'UNREAD' | 'READ'>('ALL');
  const [showFilterSheet, setShowFilterSheet] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const [notifications, setNotifications] = useState<Notification[]>([
    {
      id: '1',
      title: 'Registration Approved',
      message: 'Your company registration for ABC Corp has been approved by CIPC.',
      type: 'success',
      read: false,
      timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
    },
    {
      id: '2',
      title: 'Commission Payment',
      message: 'Your commission payment of R 2,500 has been processed and will reflect in 2-3 business days.',
      type: 'success',
      read: false,
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
    },
    {
      id: '3',
      title: 'Document Verification Required',
      message: 'Additional documents are required for XYZ Ltd registration. Please upload missing documents.',
      type: 'warning',
      read: true,
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
    },
    {
      id: '4',
      title: 'Registration Rejected',
      message: 'The registration for DEF Company was rejected. Reason: Invalid tax number provided.',
      type: 'error',
      read: true,
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    },
    {
      id: '5',
      title: 'Welcome to Comply360',
      message: 'Thank you for joining Comply360! Start by creating your first company registration.',
      type: 'info',
      read: true,
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3).toISOString(),
    },
  ]);

  const filteredNotifications = useMemo(() => {
    let filtered = [...notifications];

    if (filterType === 'UNREAD') {
      filtered = filtered.filter((n) => !n.read);
    } else if (filterType === 'READ') {
      filtered = filtered.filter((n) => n.read);
    }

    return filtered;
  }, [notifications, filterType]);

  const unreadCount = useMemo(() => {
    return notifications.filter((n) => !n.read).length;
  }, [notifications]);

  const handleRefresh = async () => {
    setRefreshing(true);
    await new Promise((resolve) => setTimeout(resolve, 1500));
    setRefreshing(false);
    dispatch(showToast({ message: 'Notifications refreshed', type: 'success' }));
  };

  const handleNotificationPress = (notification: Notification) => {
    dispatch(showToast({ message: `Opening: ${notification.title}`, type: 'info' }));
  };

  const handleMarkAsRead = (id: string) => {
    setNotifications((prev) =>
      prev.map((n) => (n.id === id ? { ...n, read: true } : n))
    );
  };

  const handleMarkAllAsRead = () => {
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    dispatch(showToast({ message: 'All notifications marked as read', type: 'success' }));
  };

  const handleClearAll = () => {
    setNotifications([]);
    dispatch(showToast({ message: 'All notifications cleared', type: 'success' }));
  };

  const renderNotification = ({ item }: { item: Notification }) => (
    <NotificationCard
      notification={item}
      onPress={handleNotificationPress}
      onMarkAsRead={handleMarkAsRead}
    />
  );

  const renderHeader = () => (
    <View style={styles.header}>
      <View style={styles.headerTop}>
        <Text style={styles.title}>Notifications</Text>
        {unreadCount > 0 && (
          <View style={styles.badge}>
            <Text style={styles.badgeText}>{unreadCount}</Text>
          </View>
        )}
      </View>

      <View style={styles.actions}>
        <TouchableOpacity
          style={styles.filterButton}
          onPress={() => setShowFilterSheet(true)}
        >
          <Icon name="filter-variant" size={20} color={colors.primary} />
          <Text style={styles.filterButtonText}>{filterType}</Text>
          {filterType !== 'ALL' && <View style={styles.filterActiveBadge} />}
        </TouchableOpacity>

        {unreadCount > 0 && (
          <TouchableOpacity style={styles.actionButton} onPress={handleMarkAllAsRead}>
            <Icon name="check-all" size={20} color={colors.textSecondary} />
            <Text style={styles.actionButtonText}>Mark All Read</Text>
          </TouchableOpacity>
        )}

        {notifications.length > 0 && (
          <TouchableOpacity style={styles.actionButton} onPress={handleClearAll}>
            <Icon name="delete-sweep" size={20} color={colors.error} />
            <Text style={[styles.actionButtonText, { color: colors.error }]}>Clear All</Text>
          </TouchableOpacity>
        )}
      </View>
    </View>
  );

  return (
    <View style={styles.container}>
      <FlatList
        data={filteredNotifications}
        renderItem={renderNotification}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContainer}
        ListHeaderComponent={renderHeader}
        ListEmptyComponent={
          <EmptyState
            icon="bell-outline"
            title={filterType === 'UNREAD' ? 'No Unread Notifications' : 'No Notifications'}
            message={
              filterType === 'UNREAD'
                ? "You're all caught up!"
                : "You haven't received any notifications yet."
            }
          />
        }
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
        }
        showsVerticalScrollIndicator={false}
      />

      <BottomSheet
        visible={showFilterSheet}
        onClose={() => setShowFilterSheet(false)}
        title="Filter Notifications"
      >
        <View style={styles.filterOptions}>
          <FilterOption
            icon="filter-outline"
            label="All Notifications"
            count={notifications.length}
            selected={filterType === 'ALL'}
            onPress={() => {
              setFilterType('ALL');
              setShowFilterSheet(false);
            }}
          />
          <FilterOption
            icon="email-mark-as-unread"
            label="Unread"
            count={unreadCount}
            selected={filterType === 'UNREAD'}
            onPress={() => {
              setFilterType('UNREAD');
              setShowFilterSheet(false);
            }}
          />
          <FilterOption
            icon="email-open"
            label="Read"
            count={notifications.length - unreadCount}
            selected={filterType === 'READ'}
            onPress={() => {
              setFilterType('READ');
              setShowFilterSheet(false);
            }}
          />
        </View>
      </BottomSheet>
    </View>
  );
};

interface FilterOptionProps {
  icon: string;
  label: string;
  count: number;
  selected: boolean;
  onPress: () => void;
}

const FilterOption: React.FC<FilterOptionProps> = ({
  icon,
  label,
  count,
  selected,
  onPress,
}) => (
  <TouchableOpacity
    style={[styles.filterOption, selected && styles.filterOptionSelected]}
    onPress={onPress}
  >
    <View style={styles.filterOptionLeft}>
      <Icon
        name={icon}
        size={24}
        color={selected ? colors.primary : '#6b7280'}
      />
      <Text style={[styles.filterOptionLabel, selected && styles.filterOptionLabelSelected]}>
        {label}
      </Text>
    </View>
    <View style={styles.filterOptionRight}>
      <View style={[styles.countBadge, selected && styles.countBadgeSelected]}>
        <Text style={[styles.countBadgeText, selected && styles.countBadgeTextSelected]}>
          {count}
        </Text>
      </View>
      {selected && <Icon name="check" size={20} color={colors.primary} />}
    </View>
  </TouchableOpacity>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  listContainer: {
    padding: 20,
  },
  header: {
    marginBottom: 20,
  },
  headerTop: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: '#111827',
  },
  badge: {
    backgroundColor: colors.primary,
    borderRadius: 12,
    paddingHorizontal: 10,
    paddingVertical: 4,
    marginLeft: 12,
  },
  badgeText: {
    fontSize: 14,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  actions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  filterButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    paddingHorizontal: 12,
    backgroundColor: `${colors.primary}10`,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.primary,
  },
  filterButtonText: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.primary,
    marginLeft: 6,
  },
  filterActiveBadge: {
    width: 6,
    height: 6,
    borderRadius: 3,
    backgroundColor: colors.primary,
    marginLeft: 6,
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    paddingHorizontal: 12,
    backgroundColor: '#FFFFFF',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e5e7eb',
  },
  actionButtonText: {
    fontSize: 14,
    fontWeight: '500',
    color: colors.textSecondary,
    marginLeft: 6,
  },
  filterOptions: {
    padding: 4,
  },
  filterOption: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 16,
    paddingHorizontal: 16,
    borderRadius: 12,
    marginBottom: 8,
    backgroundColor: '#f9fafb',
  },
  filterOptionSelected: {
    backgroundColor: `${colors.primary}10`,
  },
  filterOptionLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  filterOptionLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#6b7280',
    marginLeft: 12,
  },
  filterOptionLabelSelected: {
    color: colors.primary,
  },
  filterOptionRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  countBadge: {
    backgroundColor: '#e5e7eb',
    borderRadius: 12,
    paddingHorizontal: 10,
    paddingVertical: 4,
    minWidth: 32,
    alignItems: 'center',
  },
  countBadgeSelected: {
    backgroundColor: colors.primary,
  },
  countBadgeText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#6b7280',
  },
  countBadgeTextSelected: {
    color: '#FFFFFF',
  },
});

export default NotificationsScreen;
