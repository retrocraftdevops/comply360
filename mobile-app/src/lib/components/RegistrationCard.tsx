/**
 * RegistrationCard Component
 * Display card for registration list items
 */

import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts, shadows } from '@/lib/utils/theme';
import { formatDate, formatCurrency } from '@/lib/utils/formatting';

export interface Registration {
  id: string;
  company_name: string;
  company_type: string;
  registration_number: string;
  status: 'DRAFT' | 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'REJECTED';
  country: string;
  created_at: string;
  updated_at: string;
  amount?: number;
}

export interface RegistrationCardProps {
  registration: Registration;
  onPress: (registration: Registration) => void;
}

const RegistrationCard: React.FC<RegistrationCardProps> = ({
  registration,
  onPress,
}) => {
  /**
   * Get status color
   */
  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'COMPLETED':
        return colors.success;
      case 'IN_PROGRESS':
        return colors.warning;
      case 'PENDING':
        return colors.info;
      case 'REJECTED':
        return colors.error;
      case 'DRAFT':
      default:
        return colors.textTertiary;
    }
  };

  /**
   * Get status icon
   */
  const getStatusIcon = (status: string): string => {
    switch (status) {
      case 'COMPLETED':
        return 'check-circle';
      case 'IN_PROGRESS':
        return 'progress-clock';
      case 'PENDING':
        return 'clock-outline';
      case 'REJECTED':
        return 'close-circle';
      case 'DRAFT':
      default:
        return 'file-document-outline';
    }
  };

  /**
   * Format status text
   */
  const formatStatus = (status: string): string => {
    return status.replace('_', ' ').toLowerCase().replace(/\b\w/g, (l) => l.toUpperCase());
  };

  return (
    <TouchableOpacity
      style={styles.card}
      onPress={() => onPress(registration)}
      activeOpacity={0.7}
    >
      <View style={styles.header}>
        <View style={styles.headerLeft}>
          <Text style={styles.companyName} numberOfLines={1}>
            {registration.company_name}
          </Text>
          <Text style={styles.companyType}>{registration.company_type}</Text>
        </View>

        <View
          style={[
            styles.statusBadge,
            { backgroundColor: `${getStatusColor(registration.status)}15` },
          ]}
        >
          <Icon
            name={getStatusIcon(registration.status)}
            size={14}
            color={getStatusColor(registration.status)}
          />
          <Text
            style={[
              styles.statusText,
              { color: getStatusColor(registration.status) },
            ]}
          >
            {formatStatus(registration.status)}
          </Text>
        </View>
      </View>

      <View style={styles.details}>
        <View style={styles.detailItem}>
          <Icon name="file-document" size={16} color={colors.textSecondary} />
          <Text style={styles.detailText}>{registration.registration_number}</Text>
        </View>

        <View style={styles.detailItem}>
          <Icon name="flag" size={16} color={colors.textSecondary} />
          <Text style={styles.detailText}>{registration.country}</Text>
        </View>

        {registration.amount && (
          <View style={styles.detailItem}>
            <Icon name="currency-usd" size={16} color={colors.textSecondary} />
            <Text style={styles.detailText}>{formatCurrency(registration.amount)}</Text>
          </View>
        )}
      </View>

      <View style={styles.footer}>
        <Text style={styles.dateText}>
          Updated {formatDate(new Date(registration.updated_at), 'short')}
        </Text>
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
    fontSize: fonts.lg,
    fontWeight: '700',
    color: colors.text,
    marginBottom: 4,
  },
  companyType: {
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
  details: {
    borderTopWidth: 1,
    borderTopColor: colors.border,
    paddingTop: spacing.md,
    marginBottom: spacing.md,
  },
  detailItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  detailText: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
    marginLeft: spacing.sm,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  dateText: {
    fontSize: fonts.xs,
    color: colors.textTertiary,
  },
});

export default RegistrationCard;
