/**
 * CommissionsScreen
 * List of commissions with filtering and payout requests
 */

import React, { useState, useMemo } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  RefreshControl,
  TouchableOpacity,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import {
  CommissionCard,
  SearchBar,
  LoadingSpinner,
  EmptyState,
  BottomSheet,
  Button,
  Card,
} from '@/lib/components';
import { Commission } from '@/lib/components/CommissionCard';
import { useGetCommissionStatsQuery } from '@/store/api/commissionApi';
import { useAppDispatch } from '@/store/store';
import { showToast } from '@/store/slices/uiSlice';
import { colors, spacing, fonts } from '@/lib/utils/theme';
import { formatCurrency } from '@/lib/utils/formatting';

type FilterStatus = 'ALL' | 'PENDING' | 'APPROVED' | 'PAID' | 'REJECTED';

const CommissionsScreen: React.FC = () => {
  const dispatch = useAppDispatch();

  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('ALL');
  const [showFilterSheet, setShowFilterSheet] = useState(false);
  const [showPayoutSheet, setShowPayoutSheet] = useState(false);

  // Fetch commissions
  const {
    data: commissionData,
    isLoading,
    error,
    refetch,
  } = useGetCommissionStatsQuery();

  /**
   * Mock commissions data (since API returns stats, not list)
   */
  const mockCommissions: Commission[] = [
    {
      id: '1',
      registration_id: 'REG001',
      company_name: 'Acme Corporation',
      amount: 2500.0,
      status: 'PAID',
      commission_type: 'Registration Fee',
      created_at: '2025-12-20T10:00:00Z',
      payment_date: '2025-12-25T10:00:00Z',
    },
    {
      id: '2',
      registration_id: 'REG002',
      company_name: 'Tech Solutions Ltd',
      amount: 1800.0,
      status: 'APPROVED',
      commission_type: 'Registration Fee',
      created_at: '2025-12-22T14:30:00Z',
    },
    {
      id: '3',
      registration_id: 'REG003',
      company_name: 'Global Enterprises',
      amount: 3200.0,
      status: 'PENDING',
      commission_type: 'Renewal Fee',
      created_at: '2025-12-26T09:15:00Z',
    },
    {
      id: '4',
      registration_id: 'REG004',
      company_name: 'Innovation Inc',
      amount: 2100.0,
      status: 'PAID',
      commission_type: 'Registration Fee',
      created_at: '2025-12-18T11:20:00Z',
      payment_date: '2025-12-23T11:20:00Z',
    },
    {
      id: '5',
      registration_id: 'REG005',
      company_name: 'Smart Systems',
      amount: 1500.0,
      status: 'PENDING',
      commission_type: 'Amendment Fee',
      created_at: '2025-12-27T16:00:00Z',
    },
  ];

  /**
   * Filter and search commissions
   */
  const filteredCommissions = useMemo(() => {
    let filtered = [...mockCommissions];

    // Apply status filter
    if (filterStatus !== 'ALL') {
      filtered = filtered.filter((comm) => comm.status === filterStatus);
    }

    // Apply search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (comm) =>
          comm.company_name.toLowerCase().includes(query) ||
          comm.commission_type.toLowerCase().includes(query)
      );
    }

    return filtered;
  }, [mockCommissions, filterStatus, searchQuery]);

  /**
   * Calculate totals
   */
  const totals = useMemo(() => {
    return {
      pending: mockCommissions
        .filter((c) => c.status === 'PENDING')
        .reduce((sum, c) => sum + c.amount, 0),
      approved: mockCommissions
        .filter((c) => c.status === 'APPROVED')
        .reduce((sum, c) => sum + c.amount, 0),
      paid: mockCommissions
        .filter((c) => c.status === 'PAID')
        .reduce((sum, c) => sum + c.amount, 0),
      total: mockCommissions.reduce((sum, c) => sum + c.amount, 0),
    };
  }, [mockCommissions]);

  /**
   * Get filter counts
   */
  const filterCounts = useMemo(() => {
    const counts = {
      ALL: mockCommissions.length,
      PENDING: 0,
      APPROVED: 0,
      PAID: 0,
      REJECTED: 0,
    };

    mockCommissions.forEach((comm) => {
      if (comm.status in counts) {
        counts[comm.status]++;
      }
    });

    return counts;
  }, [mockCommissions]);

  /**
   * Handle commission press
   */
  const handleCommissionPress = (commission: Commission) => {
    console.log('[Commissions] Selected:', commission.id);
  };

  /**
   * Handle payout request
   */
  const handlePayoutRequest = () => {
    if (totals.approved === 0) {
      dispatch(
        showToast({
          message: 'No approved commissions available for payout',
          type: 'error',
        })
      );
      return;
    }

    setShowPayoutSheet(true);
  };

  /**
   * Submit payout request
   */
  const handleSubmitPayoutRequest = () => {
    setShowPayoutSheet(false);
    dispatch(
      showToast({
        message: 'Payout request submitted successfully!',
        type: 'success',
      })
    );
  };

  /**
   * Apply filter
   */
  const handleApplyFilter = (status: FilterStatus) => {
    setFilterStatus(status);
    setShowFilterSheet(false);
  };

  /**
   * Clear all filters
   */
  const handleClearFilters = () => {
    setFilterStatus('ALL');
    setSearchQuery('');
    setShowFilterSheet(false);
  };

  /**
   * Render item
   */
  const renderItem = ({ item }: { item: Commission }) => (
    <CommissionCard commission={item} onPress={handleCommissionPress} />
  );

  /**
   * Render empty state
   */
  const renderEmpty = () => {
    if (searchQuery || filterStatus !== 'ALL') {
      return (
        <EmptyState
          icon="filter-remove"
          title="No commissions found"
          description="Try adjusting your search or filters"
          actionLabel="Clear Filters"
          onActionPress={handleClearFilters}
        />
      );
    }

    return (
      <EmptyState
        icon="currency-usd-off"
        title="No commissions yet"
        description="Commissions will appear here once registrations are completed"
      />
    );
  };

  /**
   * Render loading state
   */
  if (isLoading && !commissionData) {
    return (
      <View style={styles.centerContainer}>
        <LoadingSpinner size="large" message="Loading commissions..." />
      </View>
    );
  }

  const hasActiveFilters = filterStatus !== 'ALL' || searchQuery.length > 0;

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Commissions</Text>
        <TouchableOpacity
          style={styles.payoutButton}
          onPress={handlePayoutRequest}
        >
          <Icon name="bank-transfer" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      {/* Summary Cards */}
      <View style={styles.summaryContainer}>
        <View style={styles.summaryRow}>
          <SummaryCard
            label="Pending"
            amount={totals.pending}
            color={colors.warning}
            icon="clock-outline"
          />
          <SummaryCard
            label="Approved"
            amount={totals.approved}
            color={colors.info}
            icon="check-decagram"
          />
        </View>
        <View style={styles.summaryRow}>
          <SummaryCard
            label="Paid"
            amount={totals.paid}
            color={colors.success}
            icon="cash-check"
          />
          <SummaryCard
            label="Total"
            amount={totals.total}
            color={colors.primary}
            icon="currency-usd"
          />
        </View>
      </View>

      {/* Search Bar */}
      <View style={styles.searchContainer}>
        <SearchBar
          value={searchQuery}
          onChangeText={setSearchQuery}
          placeholder="Search commissions..."
          showFilter
          filterActive={hasActiveFilters}
          onFilterPress={() => setShowFilterSheet(true)}
        />
      </View>

      {/* Results Count */}
      <View style={styles.resultsContainer}>
        <Text style={styles.resultsText}>
          {filteredCommissions.length}{' '}
          {filteredCommissions.length === 1 ? 'commission' : 'commissions'}
        </Text>
        {hasActiveFilters && (
          <TouchableOpacity onPress={handleClearFilters}>
            <Text style={styles.clearText}>Clear All</Text>
          </TouchableOpacity>
        )}
      </View>

      {/* Commissions List */}
      <FlatList
        data={filteredCommissions}
        renderItem={renderItem}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.listContent}
        ListEmptyComponent={renderEmpty}
        refreshControl={
          <RefreshControl
            refreshing={isLoading}
            onRefresh={refetch}
            tintColor={colors.primary}
          />
        }
        showsVerticalScrollIndicator={false}
      />

      {/* Filter Bottom Sheet */}
      <BottomSheet
        visible={showFilterSheet}
        onClose={() => setShowFilterSheet(false)}
        title="Filter Commissions"
        showHandle
        snapPoints={[0.5]}
      >
        <View style={styles.filterContent}>
          <Text style={styles.filterLabel}>Status</Text>

          <FilterOption
            label="All Commissions"
            count={filterCounts.ALL}
            isSelected={filterStatus === 'ALL'}
            onPress={() => handleApplyFilter('ALL')}
          />

          <FilterOption
            label="Pending"
            count={filterCounts.PENDING}
            isSelected={filterStatus === 'PENDING'}
            onPress={() => handleApplyFilter('PENDING')}
            color={colors.warning}
          />

          <FilterOption
            label="Approved"
            count={filterCounts.APPROVED}
            isSelected={filterStatus === 'APPROVED'}
            onPress={() => handleApplyFilter('APPROVED')}
            color={colors.info}
          />

          <FilterOption
            label="Paid"
            count={filterCounts.PAID}
            isSelected={filterStatus === 'PAID'}
            onPress={() => handleApplyFilter('PAID')}
            color={colors.success}
          />

          <FilterOption
            label="Rejected"
            count={filterCounts.REJECTED}
            isSelected={filterStatus === 'REJECTED'}
            onPress={() => handleApplyFilter('REJECTED')}
            color={colors.error}
          />

          <View style={styles.filterActions}>
            <Button
              title="Clear Filters"
              onPress={handleClearFilters}
              variant="outline"
            />
          </View>
        </View>
      </BottomSheet>

      {/* Payout Request Bottom Sheet */}
      <BottomSheet
        visible={showPayoutSheet}
        onClose={() => setShowPayoutSheet(false)}
        title="Request Payout"
        showHandle
        snapPoints={[0.5]}
      >
        <View style={styles.payoutContent}>
          <Card variant="filled" padding="large">
            <Text style={styles.payoutLabel}>Available for Payout</Text>
            <Text style={styles.payoutAmount}>{formatCurrency(totals.approved)}</Text>
          </Card>

          <View style={styles.payoutInfo}>
            <Icon name="information" size={24} color={colors.info} />
            <Text style={styles.payoutInfoText}>
              Only approved commissions can be requested for payout. Funds will be
              transferred within 3-5 business days.
            </Text>
          </View>

          <View style={styles.payoutActions}>
            <Button
              title="Cancel"
              onPress={() => setShowPayoutSheet(false)}
              variant="outline"
            />
            <View style={spacing.sm} />
            <Button
              title="Submit Request"
              onPress={handleSubmitPayoutRequest}
              variant="primary"
              icon="bank-transfer"
              disabled={totals.approved === 0}
            />
          </View>
        </View>
      </BottomSheet>
    </View>
  );
};

/**
 * Summary Card Component
 */
interface SummaryCardProps {
  label: string;
  amount: number;
  color: string;
  icon: string;
}

const SummaryCard: React.FC<SummaryCardProps> = ({ label, amount, color, icon }) => (
  <View style={[styles.summaryCard, { borderLeftColor: color }]}>
    <View style={styles.summaryCardHeader}>
      <Text style={styles.summaryCardLabel}>{label}</Text>
      <Icon name={icon} size={20} color={color} />
    </View>
    <Text style={styles.summaryCardAmount}>{formatCurrency(amount)}</Text>
  </View>
);

/**
 * Filter Option Component
 */
interface FilterOptionProps {
  label: string;
  count: number;
  isSelected: boolean;
  onPress: () => void;
  color?: string;
}

const FilterOption: React.FC<FilterOptionProps> = ({
  label,
  count,
  isSelected,
  onPress,
  color = colors.text,
}) => (
  <TouchableOpacity
    style={[styles.filterOption, isSelected && styles.filterOptionSelected]}
    onPress={onPress}
  >
    <View style={styles.filterOptionLeft}>
      <View style={[styles.filterIndicator, { backgroundColor: color }]} />
      <Text style={[styles.filterOptionLabel, isSelected && styles.filterOptionLabelSelected]}>
        {label}
      </Text>
    </View>

    <View style={styles.filterOptionRight}>
      <Text style={[styles.filterOptionCount, isSelected && styles.filterOptionCountSelected]}>
        {count}
      </Text>
      {isSelected && <Icon name="check" size={20} color={colors.primary} />}
    </View>
  </TouchableOpacity>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    paddingBottom: spacing.md,
    backgroundColor: '#FFFFFF',
  },
  title: {
    fontSize: fonts['2xl'],
    fontWeight: '700',
    color: colors.text,
  },
  payoutButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  summaryContainer: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  summaryRow: {
    flexDirection: 'row',
    gap: spacing.md,
    marginBottom: spacing.md,
  },
  summaryCard: {
    flex: 1,
    backgroundColor: colors.background,
    borderRadius: 8,
    padding: spacing.md,
    borderLeftWidth: 4,
  },
  summaryCardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  summaryCardLabel: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
  },
  summaryCardAmount: {
    fontSize: fonts.lg,
    fontWeight: '700',
    color: colors.text,
  },
  searchContainer: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    backgroundColor: '#FFFFFF',
  },
  resultsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
  },
  resultsText: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
  },
  clearText: {
    fontSize: fonts.sm,
    color: colors.primary,
    fontWeight: '600',
  },
  listContent: {
    padding: spacing.lg,
  },
  filterContent: {
    paddingBottom: spacing.lg,
  },
  filterLabel: {
    fontSize: fonts.sm,
    fontWeight: '600',
    color: colors.textSecondary,
    marginBottom: spacing.md,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  filterOption: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.md,
    borderRadius: 8,
    marginBottom: spacing.xs,
  },
  filterOptionSelected: {
    backgroundColor: `${colors.primary}10`,
  },
  filterOptionLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  filterIndicator: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: spacing.md,
  },
  filterOptionLabel: {
    fontSize: fonts.base,
    color: colors.text,
  },
  filterOptionLabelSelected: {
    fontWeight: '600',
    color: colors.primary,
  },
  filterOptionRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  filterOptionCount: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
    fontWeight: '600',
  },
  filterOptionCountSelected: {
    color: colors.primary,
  },
  filterActions: {
    marginTop: spacing.xl,
  },
  payoutContent: {
    paddingBottom: spacing.lg,
  },
  payoutLabel: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
  },
  payoutAmount: {
    fontSize: fonts['3xl'],
    fontWeight: '700',
    color: colors.primary,
  },
  payoutInfo: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginTop: spacing.lg,
    padding: spacing.md,
    backgroundColor: `${colors.info}10`,
    borderRadius: 8,
  },
  payoutInfoText: {
    flex: 1,
    marginLeft: spacing.md,
    fontSize: fonts.sm,
    color: colors.text,
    lineHeight: 20,
  },
  payoutActions: {
    marginTop: spacing.xl,
  },
});

export default CommissionsScreen;
