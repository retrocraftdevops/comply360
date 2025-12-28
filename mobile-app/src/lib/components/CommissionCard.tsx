/**
 * CommissionCard Component
 * Display card for commission list items
 */

import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts, shadows } from '@/lib/utils/theme';
import { formatDate, formatCurrency } from '@/lib/utils/formatting';

export interface Commission {
  id: string;
  registration_id: string;
  company_name: string;
  amount: number;
  status: 'PENDING' | 'APPROVED' | 'PAID' | 'REJECTED';
  commission_type: string;
  created_at: string;
  payment_date?: string;
}

export interface CommissionCardProps {
  commission: Commission;
  onPress: (commission: Commission) => void;
}

const CommissionCard: React.FC<CommissionCardProps> = ({
  commission,
  onPress,
}) => {
  /**
   * Get status color
   */
  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'PAID':
        return colors.success;
      case 'APPROVED':
        return colors.info;
      case 'PENDING':
        return colors.warning;
      case 'REJECTED':
        return colors.error;
      default:
        return colors.textTertiary;
    }
  };

  /**
   * Get status icon
   */
  const getStatusIcon = (status: string): string => {
    switch (status) {
      case 'PAID':
        return 'check-circle';
      case 'APPROVED':
        return 'check-decagram';
      case 'PENDING':
        return 'clock-outline';
      case 'REJECTED':
        return 'close-circle';
      default:
        return 'help-circle';
    }
  };

  return (
    <TouchableOpacity
      style={styles.card}
      onPress={() => onPress(commission)}
      activeOpacity={0.7}
    >
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.companyName} numberOfLines={1}>
            {commission.company_name}
          </Text>
          <Text style={styles.commissionType}>{commission.commission_type}</Text>
        </View>

        <View
          style={[
            styles.statusBadge,
            { backgroundColor: `${getStatusColor(commission.status)}15` },
          ]}
        >
          <Icon
            name={getStatusIcon(commission.status)}
            size={14}
            color={getStatusColor(commission.status)}
          />
          <Text
            style={[
              styles.statusText,
              { color: getStatusColor(commission.status) },
            ]}
          >
            {commission.status}
          </Text>
        </View>
      </View>

      <View style={styles.amountContainer}>
        <Text style={styles.amountLabel}>Commission Amount</Text>
        <Text style={styles.amount}>{formatCurrency(commission.amount)}</Text>
      </View>

      <View style={styles.footer}>
        <View style={styles.dateInfo}>
          <Icon name="calendar" size={14} color={colors.textSecondary} />
          <Text style={styles.dateText}>
            {commission.payment_date
              ? `Paid ${formatDate(new Date(commission.payment_date), 'short')}`
              : `Created ${formatDate(new Date(commission.created_at), 'short')}`}
          </Text>
        </View>
        <Icon name="chevron-right" size={20} color={colors.textTertiary} />
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: spacing.md,
    marginBottom: spacing.md,
    ...shadows.base,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: spacing.md,
  },
  headerLeft: {
    flex: 1,
    marginRight: spacing.sm,
  },
  companyName: {
    fontSize: fonts.base,
    fontWeight: '600',
    color: colors.text,
    marginBottom: 4,
  },
  commissionType: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: 12,
  },
  statusText: {
    fontSize: fonts.xs,
    fontWeight: '600',
    marginLeft: 4,
  },
  amountContainer: {
    borderTopWidth: 1,
    borderTopColor: colors.border,
    paddingTop: spacing.md,
    marginBottom: spacing.md,
  },
  amountLabel: {
    fontSize: fonts.xs,
    color: colors.textSecondary,
    marginBottom: 4,
  },
  amount: {
    fontSize: fonts['2xl'],
    fontWeight: '700',
    color: colors.primary,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  dateInfo: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  dateText: {
    fontSize: fonts.xs,
    color: colors.textSecondary,
    marginLeft: spacing.xs,
  },
});

export default CommissionCard;
