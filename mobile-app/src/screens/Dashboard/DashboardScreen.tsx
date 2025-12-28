/**
 * Dashboard Screen
 * Main dashboard with real-time stats, recent activity, and quick actions
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
  TouchableOpacity,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useAppSelector } from '../../store/store';
import { useGetRegistrationStatsQuery } from '../../store/api/registrationApi';
import { useGetCommissionStatsQuery } from '../../store/api/commissionApi';
import { useGetDocumentStatsQuery } from '../../store/api/documentApi';
import { LoadingSpinner, Card } from '../../lib/components';

interface StatCard {
  id: string;
  title: string;
  value: string;
  change: string;
  changeType: 'increase' | 'decrease' | 'neutral';
  icon: string;
  color: string;
  loading?: boolean;
}

interface QuickAction {
  id: string;
  title: string;
  icon: string;
  color: string;
  route: string;
}

const DashboardScreen = () => {
  const { user } = useAppSelector((state) => state.auth);

  // Fetch stats from API
  const {
    data: registrationStats,
    isLoading: isLoadingRegistrations,
    refetch: refetchRegistrations,
  } = useGetRegistrationStatsQuery();

  const {
    data: commissionStats,
    isLoading: isLoadingCommissions,
    refetch: refetchCommissions,
  } = useGetCommissionStatsQuery();

  const {
    data: documentStats,
    isLoading: isLoadingDocuments,
    refetch: refetchDocuments,
  } = useGetDocumentStatsQuery();

  const [refreshing, setRefreshing] = useState(false);

  // Build stats from API data
  const stats: StatCard[] = [
    {
      id: '1',
      title: 'Active Registrations',
      value: isLoadingRegistrations
        ? '...'
        : `${registrationStats?.in_progress || 0}`,
      change: `${registrationStats?.total || 0} total`,
      changeType: 'neutral',
      icon: 'file-document',
      color: '#7c3aed',
      loading: isLoadingRegistrations,
    },
    {
      id: '2',
      title: 'Total Commissions',
      value: isLoadingCommissions
        ? '...'
        : `R ${(commissionStats?.total_earned || 0).toLocaleString()}`,
      change: `R ${(commissionStats?.this_month_earned || 0).toLocaleString()} this month`,
      changeType: commissionStats?.this_month_earned
        ? commissionStats.this_month_earned > 0
          ? 'increase'
          : 'neutral'
        : 'neutral',
      icon: 'cash-multiple',
      color: '#10b981',
      loading: isLoadingCommissions,
    },
    {
      id: '3',
      title: 'Pending Documents',
      value: isLoadingDocuments ? '...' : `${documentStats?.pending || 0}`,
      change: documentStats?.pending
        ? documentStats.pending > 5
          ? `${documentStats.pending - 5} over target`
          : 'within target'
        : 'no pending',
      changeType: documentStats?.pending
        ? documentStats.pending > 5
          ? 'neutral'
          : 'neutral'
        : 'neutral',
      icon: 'folder-open',
      color: '#f59e0b',
      loading: isLoadingDocuments,
    },
    {
      id: '4',
      title: 'Completed',
      value: isLoadingRegistrations
        ? '...'
        : `${registrationStats?.completed || 0}`,
      change: 'registrations done',
      changeType: 'increase',
      icon: 'check-circle',
      color: '#3b82f6',
      loading: isLoadingRegistrations,
    },
  ];

  const quickActions: QuickAction[] = [
    {
      id: '1',
      title: 'New Registration',
      icon: 'plus-circle',
      color: '#7c3aed',
      route: 'Registrations',
    },
    {
      id: '2',
      title: 'Upload Document',
      icon: 'upload',
      color: '#3b82f6',
      route: 'Documents',
    },
    {
      id: '3',
      title: 'View Commissions',
      icon: 'currency-usd',
      color: '#10b981',
      route: 'Commissions',
    },
    {
      id: '4',
      title: 'Profile',
      icon: 'account',
      color: '#6b7280',
      route: 'Profile',
    },
  ];

  const onRefresh = async () => {
    setRefreshing(true);
    await Promise.all([
      refetchRegistrations(),
      refetchCommissions(),
      refetchDocuments(),
    ]);
    setRefreshing(false);
  };

  const handleQuickAction = (route: string) => {
    console.log('[Dashboard] Quick action clicked:', route);
    // TODO: Navigate to route using navigation
  };

  const getGreeting = (): string => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good morning';
    if (hour < 18) return 'Good afternoon';
    return 'Good evening';
  };

  const getUserName = (): string => {
    if (user?.first_name) {
      return user.first_name;
    }
    return 'there';
  };

  const formatCurrency = (amount: number): string => {
    return `R ${amount.toLocaleString('en-ZA', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    })}`;
  };

  // Show full-screen loading on first load
  const isFirstLoad =
    isLoadingRegistrations && isLoadingCommissions && isLoadingDocuments;

  if (isFirstLoad) {
    return <LoadingSpinner fullScreen message="Loading dashboard..." />;
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.contentContainer}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={onRefresh}
          colors={['#7c3aed']}
          tintColor="#7c3aed"
        />
      }
    >
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>{getGreeting()},</Text>
          <Text style={styles.userName}>{getUserName()}</Text>
        </View>
        <TouchableOpacity style={styles.notificationButton}>
          <Icon name="bell-outline" size={24} color="#111827" />
          <View style={styles.notificationBadge}>
            <Text style={styles.notificationBadgeText}>
              {commissionStats?.count_pending || 0}
            </Text>
          </View>
        </TouchableOpacity>
      </View>

      {/* Stats Grid */}
      <View style={styles.statsGrid}>
        {stats.map((stat) => (
          <Card key={stat.id} style={styles.statCard} padding="medium">
            <View style={[styles.statIconContainer, { backgroundColor: `${stat.color}15` }]}>
              <Icon name={stat.icon} size={24} color={stat.color} />
            </View>
            <Text style={styles.statValue}>{stat.value}</Text>
            <Text style={styles.statTitle}>{stat.title}</Text>
            <View style={styles.statChangeContainer}>
              <Text
                style={[
                  styles.statChange,
                  stat.changeType === 'increase' && styles.statChangePositive,
                  stat.changeType === 'decrease' && styles.statChangeNegative,
                ]}
              >
                {stat.change}
              </Text>
            </View>
          </Card>
        ))}
      </View>

      {/* Quick Actions */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Quick Actions</Text>
        <View style={styles.quickActionsGrid}>
          {quickActions.map((action) => (
            <Card
              key={action.id}
              style={styles.quickActionCard}
              padding="medium"
              onPress={() => handleQuickAction(action.route)}
            >
              <View style={[styles.quickActionIcon, { backgroundColor: `${action.color}15` }]}>
                <Icon name={action.icon} size={28} color={action.color} />
              </View>
              <Text style={styles.quickActionTitle}>{action.title}</Text>
            </Card>
          ))}
        </View>
      </View>

      {/* Commission Summary */}
      {commissionStats && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Commission Summary</Text>
          <Card padding="large">
            <View style={styles.summaryRow}>
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Earned</Text>
                <Text style={[styles.summaryValue, { color: '#10b981' }]}>
                  {formatCurrency(commissionStats.total_earned)}
                </Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Paid</Text>
                <Text style={[styles.summaryValue, { color: '#3b82f6' }]}>
                  {formatCurrency(commissionStats.total_paid)}
                </Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Pending</Text>
                <Text style={[styles.summaryValue, { color: '#f59e0b' }]}>
                  {formatCurrency(commissionStats.total_pending)}
                </Text>
              </View>
            </View>
          </Card>
        </View>
      )}

      {/* Document Summary */}
      {documentStats && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Document Status</Text>
          <Card padding="large">
            <View style={styles.summaryRow}>
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Total</Text>
                <Text style={styles.summaryValue}>{documentStats.total}</Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Verified</Text>
                <Text style={[styles.summaryValue, { color: '#10b981' }]}>
                  {documentStats.verified}
                </Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Pending</Text>
                <Text style={[styles.summaryValue, { color: '#f59e0b' }]}>
                  {documentStats.pending}
                </Text>
              </View>
            </View>
          </Card>
        </View>
      )}

      {/* Registration Summary */}
      {registrationStats && (
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Registration Status</Text>
          <Card padding="large">
            <View style={styles.summaryRow}>
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Total</Text>
                <Text style={styles.summaryValue}>{registrationStats.total}</Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>Completed</Text>
                <Text style={[styles.summaryValue, { color: '#10b981' }]}>
                  {registrationStats.completed}
                </Text>
              </View>
              <View style={styles.summaryDivider} />
              <View style={styles.summaryItem}>
                <Text style={styles.summaryLabel}>In Progress</Text>
                <Text style={[styles.summaryValue, { color: '#3b82f6' }]}>
                  {registrationStats.in_progress}
                </Text>
              </View>
            </View>
          </Card>
        </View>
      )}
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f9fafb',
  },
  contentContainer: {
    padding: 16,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 24,
  },
  greeting: {
    fontSize: 16,
    color: '#6b7280',
    marginBottom: 4,
  },
  userName: {
    fontSize: 24,
    fontWeight: '700',
    color: '#111827',
  },
  notificationButton: {
    position: 'relative',
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#FFFFFF',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  notificationBadge: {
    position: 'absolute',
    top: 8,
    right: 8,
    minWidth: 18,
    height: 18,
    borderRadius: 9,
    backgroundColor: '#ef4444',
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 4,
  },
  notificationBadgeText: {
    fontSize: 10,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  statsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginHorizontal: -6,
    marginBottom: 24,
  },
  statCard: {
    width: '50%',
    padding: 6,
  },
  statIconContainer: {
    width: 48,
    height: 48,
    borderRadius: 24,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 12,
  },
  statValue: {
    fontSize: 24,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 4,
  },
  statTitle: {
    fontSize: 13,
    color: '#6b7280',
    marginBottom: 6,
  },
  statChangeContainer: {
    marginTop: 4,
  },
  statChange: {
    fontSize: 12,
    color: '#6b7280',
  },
  statChangePositive: {
    color: '#10b981',
  },
  statChangeNegative: {
    color: '#ef4444',
  },
  section: {
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 12,
  },
  quickActionsGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginHorizontal: -6,
  },
  quickActionCard: {
    width: '50%',
    padding: 6,
    alignItems: 'center',
  },
  quickActionIcon: {
    width: 64,
    height: 64,
    borderRadius: 32,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 12,
  },
  quickActionTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: '#111827',
    textAlign: 'center',
  },
  summaryRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  summaryItem: {
    flex: 1,
    alignItems: 'center',
  },
  summaryLabel: {
    fontSize: 13,
    color: '#6b7280',
    marginBottom: 6,
  },
  summaryValue: {
    fontSize: 20,
    fontWeight: '700',
    color: '#111827',
  },
  summaryDivider: {
    width: 1,
    height: 40,
    backgroundColor: '#e5e7eb',
  },
});

export default DashboardScreen;
