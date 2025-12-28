/**
 * RegistrationsScreen
 * List of company registrations with search and filtering
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
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import {
  RegistrationCard,
  SearchBar,
  LoadingSpinner,
  EmptyState,
  BottomSheet,
  Button,
} from '@/lib/components';
import { Registration } from '@/lib/components/RegistrationCard';
import { useGetRegistrationsQuery } from '@/store/api/registrationApi';
import { colors, spacing, fonts } from '@/lib/utils/theme';

type FilterStatus = 'ALL' | 'DRAFT' | 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'REJECTED';

const RegistrationsScreen: React.FC = () => {
  const navigation = useNavigation();

  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('ALL');
  const [showFilterSheet, setShowFilterSheet] = useState(false);

  // Fetch registrations
  const {
    data: registrationsData,
    isLoading,
    error,
    refetch,
  } = useGetRegistrationsQuery({ page: 1, limit: 100 });

  /**
   * Filter and search registrations
   */
  const filteredRegistrations = useMemo(() => {
    if (!registrationsData?.registrations) return [];

    let filtered = [...registrationsData.registrations];

    // Apply status filter
    if (filterStatus !== 'ALL') {
      filtered = filtered.filter((reg) => reg.status === filterStatus);
    }

    // Apply search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (reg) =>
          reg.company_name.toLowerCase().includes(query) ||
          reg.registration_number.toLowerCase().includes(query) ||
          reg.company_type.toLowerCase().includes(query)
      );
    }

    return filtered;
  }, [registrationsData, filterStatus, searchQuery]);

  /**
   * Get filter counts
   */
  const filterCounts = useMemo(() => {
    if (!registrationsData?.registrations) {
      return {
        ALL: 0,
        DRAFT: 0,
        PENDING: 0,
        IN_PROGRESS: 0,
        COMPLETED: 0,
        REJECTED: 0,
      };
    }

    const counts = {
      ALL: registrationsData.registrations.length,
      DRAFT: 0,
      PENDING: 0,
      IN_PROGRESS: 0,
      COMPLETED: 0,
      REJECTED: 0,
    };

    registrationsData.registrations.forEach((reg) => {
      if (reg.status in counts) {
        counts[reg.status]++;
      }
    });

    return counts;
  }, [registrationsData]);

  /**
   * Handle registration press
   */
  const handleRegistrationPress = (registration: Registration) => {
    // TODO: Navigate to RegistrationDetailsScreen
    console.log('[Registrations] Selected:', registration.id);
  };

  /**
   * Handle new registration
   */
  const handleNewRegistration = () => {
    navigation.navigate('NewRegistration' as never);
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
  const renderItem = ({ item }: { item: Registration }) => (
    <RegistrationCard registration={item} onPress={handleRegistrationPress} />
  );

  /**
   * Render empty state
   */
  const renderEmpty = () => {
    if (searchQuery || filterStatus !== 'ALL') {
      return (
        <EmptyState
          icon="filter-remove"
          title="No registrations found"
          description="Try adjusting your search or filters"
          actionLabel="Clear Filters"
          onActionPress={handleClearFilters}
        />
      );
    }

    return (
      <EmptyState
        icon="file-document-plus"
        title="No registrations yet"
        description="Create your first company registration to get started"
        actionLabel="New Registration"
        onActionPress={handleNewRegistration}
      />
    );
  };

  /**
   * Render loading state
   */
  if (isLoading && !registrationsData) {
    return (
      <View style={styles.centerContainer}>
        <LoadingSpinner size="large" message="Loading registrations..." />
      </View>
    );
  }

  /**
   * Render error state
   */
  if (error) {
    return (
      <View style={styles.centerContainer}>
        <EmptyState
          icon="alert-circle"
          title="Failed to load registrations"
          description="Please check your connection and try again"
          actionLabel="Retry"
          onActionPress={() => refetch()}
        />
      </View>
    );
  }

  const hasActiveFilters = filterStatus !== 'ALL' || searchQuery.length > 0;

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Registrations</Text>
        <TouchableOpacity
          style={styles.newButton}
          onPress={handleNewRegistration}
        >
          <Icon name="plus" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      {/* Search Bar */}
      <View style={styles.searchContainer}>
        <SearchBar
          value={searchQuery}
          onChangeText={setSearchQuery}
          placeholder="Search registrations..."
          showFilter
          filterActive={hasActiveFilters}
          onFilterPress={() => setShowFilterSheet(true)}
        />
      </View>

      {/* Results Count */}
      <View style={styles.resultsContainer}>
        <Text style={styles.resultsText}>
          {filteredRegistrations.length}{' '}
          {filteredRegistrations.length === 1 ? 'registration' : 'registrations'}
        </Text>
        {hasActiveFilters && (
          <TouchableOpacity onPress={handleClearFilters}>
            <Text style={styles.clearText}>Clear All</Text>
          </TouchableOpacity>
        )}
      </View>

      {/* Registrations List */}
      <FlatList
        data={filteredRegistrations}
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
        title="Filter Registrations"
        showHandle
        snapPoints={[0.6]}
      >
        <View style={styles.filterContent}>
          <Text style={styles.filterLabel}>Status</Text>

          <FilterOption
            label="All Registrations"
            count={filterCounts.ALL}
            isSelected={filterStatus === 'ALL'}
            onPress={() => handleApplyFilter('ALL')}
          />

          <FilterOption
            label="Draft"
            count={filterCounts.DRAFT}
            isSelected={filterStatus === 'DRAFT'}
            onPress={() => handleApplyFilter('DRAFT')}
            color={colors.textTertiary}
          />

          <FilterOption
            label="Pending"
            count={filterCounts.PENDING}
            isSelected={filterStatus === 'PENDING'}
            onPress={() => handleApplyFilter('PENDING')}
            color={colors.info}
          />

          <FilterOption
            label="In Progress"
            count={filterCounts.IN_PROGRESS}
            isSelected={filterStatus === 'IN_PROGRESS'}
            onPress={() => handleApplyFilter('IN_PROGRESS')}
            color={colors.warning}
          />

          <FilterOption
            label="Completed"
            count={filterCounts.COMPLETED}
            isSelected={filterStatus === 'COMPLETED'}
            onPress={() => handleApplyFilter('COMPLETED')}
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
    </View>
  );
};

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
  newButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  searchContainer: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
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
});

export default RegistrationsScreen;
