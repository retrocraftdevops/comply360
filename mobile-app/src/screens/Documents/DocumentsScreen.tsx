/**
 * DocumentsScreen
 * List of documents with filtering and upload functionality
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
  DocumentCard,
  SearchBar,
  LoadingSpinner,
  EmptyState,
  BottomSheet,
  Button,
  Card,
} from '@/lib/components';
import { Document } from '@/lib/components/DocumentCard';
import { useGetDocumentStatsQuery } from '@/store/api/documentApi';
import { useAppDispatch } from '@/store/store';
import { showToast } from '@/store/slices/uiSlice';
import { colors, spacing, fonts } from '@/lib/utils/theme';

type FilterStatus = 'ALL' | 'PENDING' | 'VERIFIED' | 'REJECTED';
type FilterType = 'ALL' | 'PDF' | 'IMAGE' | 'DOCUMENT' | 'SPREADSHEET';

const DocumentsScreen: React.FC = () => {
  const dispatch = useAppDispatch();

  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('ALL');
  const [filterType, setFilterType] = useState<FilterType>('ALL');
  const [showFilterSheet, setShowFilterSheet] = useState(false);
  const [showUploadSheet, setShowUploadSheet] = useState(false);

  // Fetch documents
  const {
    data: documentData,
    isLoading,
    error,
    refetch,
  } = useGetDocumentStatsQuery();

  /**
   * Mock documents data
   */
  const mockDocuments: Document[] = [
    {
      id: '1',
      name: 'Company Registration Certificate.pdf',
      type: 'application/pdf',
      size: 2458624, // 2.4 MB
      status: 'VERIFIED',
      uploaded_at: '2025-12-20T10:00:00Z',
      registration_id: 'REG001',
      company_name: 'Acme Corporation',
    },
    {
      id: '2',
      name: 'Tax Clearance Certificate.pdf',
      type: 'application/pdf',
      size: 1856432, // 1.8 MB
      status: 'VERIFIED',
      uploaded_at: '2025-12-22T14:30:00Z',
      registration_id: 'REG001',
      company_name: 'Acme Corporation',
    },
    {
      id: '3',
      name: 'Director ID Copy.jpg',
      type: 'image/jpeg',
      size: 3245678, // 3.2 MB
      status: 'PENDING',
      uploaded_at: '2025-12-26T09:15:00Z',
      registration_id: 'REG002',
      company_name: 'Tech Solutions Ltd',
    },
    {
      id: '4',
      name: 'Business Plan 2025.docx',
      type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      size: 5678432, // 5.7 MB
      status: 'VERIFIED',
      uploaded_at: '2025-12-18T11:20:00Z',
      registration_id: 'REG003',
      company_name: 'Global Enterprises',
    },
    {
      id: '5',
      name: 'Financial Statements.xlsx',
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      size: 4234567, // 4.2 MB
      status: 'PENDING',
      uploaded_at: '2025-12-27T16:00:00Z',
      registration_id: 'REG004',
      company_name: 'Innovation Inc',
    },
    {
      id: '6',
      name: 'Proof of Address.jpg',
      type: 'image/jpeg',
      size: 2145678, // 2.1 MB
      status: 'REJECTED',
      uploaded_at: '2025-12-25T12:00:00Z',
      registration_id: 'REG005',
      company_name: 'Smart Systems',
    },
  ];

  /**
   * Get document type category
   */
  const getDocumentTypeCategory = (type: string): FilterType => {
    const lowerType = type.toLowerCase();
    if (lowerType.includes('pdf')) return 'PDF';
    if (lowerType.includes('image') || lowerType.includes('jpeg') || lowerType.includes('png')) return 'IMAGE';
    if (lowerType.includes('word') || lowerType.includes('doc')) return 'DOCUMENT';
    if (lowerType.includes('sheet') || lowerType.includes('excel') || lowerType.includes('xls')) return 'SPREADSHEET';
    return 'ALL';
  };

  /**
   * Filter and search documents
   */
  const filteredDocuments = useMemo(() => {
    let filtered = [...mockDocuments];

    // Apply status filter
    if (filterStatus !== 'ALL') {
      filtered = filtered.filter((doc) => doc.status === filterStatus);
    }

    // Apply type filter
    if (filterType !== 'ALL') {
      filtered = filtered.filter((doc) => getDocumentTypeCategory(doc.type) === filterType);
    }

    // Apply search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (doc) =>
          doc.name.toLowerCase().includes(query) ||
          (doc.company_name && doc.company_name.toLowerCase().includes(query))
      );
    }

    return filtered;
  }, [mockDocuments, filterStatus, filterType, searchQuery]);

  /**
   * Get filter counts
   */
  const filterCounts = useMemo(() => {
    const statusCounts = {
      ALL: mockDocuments.length,
      PENDING: 0,
      VERIFIED: 0,
      REJECTED: 0,
    };

    const typeCounts = {
      ALL: mockDocuments.length,
      PDF: 0,
      IMAGE: 0,
      DOCUMENT: 0,
      SPREADSHEET: 0,
    };

    mockDocuments.forEach((doc) => {
      // Status counts
      if (doc.status in statusCounts) {
        statusCounts[doc.status]++;
      }

      // Type counts
      const typeCategory = getDocumentTypeCategory(doc.type);
      if (typeCategory !== 'ALL' && typeCategory in typeCounts) {
        typeCounts[typeCategory]++;
      }
    });

    return { status: statusCounts, type: typeCounts };
  }, [mockDocuments]);

  /**
   * Handle document press
   */
  const handleDocumentPress = (document: Document) => {
    console.log('[Documents] Selected:', document.id);
    dispatch(
      showToast({
        message: `Opening ${document.name}`,
        type: 'info',
      })
    );
  };

  /**
   * Handle document download
   */
  const handleDocumentDownload = (document: Document) => {
    console.log('[Documents] Download:', document.id);
    dispatch(
      showToast({
        message: `Downloading ${document.name}`,
        type: 'success',
      })
    );
  };

  /**
   * Handle upload document
   */
  const handleUploadDocument = () => {
    setShowUploadSheet(true);
  };

  /**
   * Handle upload from camera
   */
  const handleUploadFromCamera = () => {
    setShowUploadSheet(false);
    dispatch(
      showToast({
        message: 'Camera feature coming in next update!',
        type: 'info',
      })
    );
  };

  /**
   * Handle upload from gallery
   */
  const handleUploadFromGallery = () => {
    setShowUploadSheet(false);
    dispatch(
      showToast({
        message: 'Gallery picker coming in next update!',
        type: 'info',
      })
    );
  };

  /**
   * Handle upload from files
   */
  const handleUploadFromFiles = () => {
    setShowUploadSheet(false);
    dispatch(
      showToast({
        message: 'File picker coming in next update!',
        type: 'info',
      })
    );
  };

  /**
   * Clear all filters
   */
  const handleClearFilters = () => {
    setFilterStatus('ALL');
    setFilterType('ALL');
    setSearchQuery('');
    setShowFilterSheet(false);
  };

  /**
   * Render item
   */
  const renderItem = ({ item }: { item: Document }) => (
    <DocumentCard
      document={item}
      onPress={handleDocumentPress}
      onDownload={handleDocumentDownload}
    />
  );

  /**
   * Render empty state
   */
  const renderEmpty = () => {
    if (searchQuery || filterStatus !== 'ALL' || filterType !== 'ALL') {
      return (
        <EmptyState
          icon="filter-remove"
          title="No documents found"
          description="Try adjusting your search or filters"
          actionLabel="Clear Filters"
          onActionPress={handleClearFilters}
        />
      );
    }

    return (
      <EmptyState
        icon="file-upload"
        title="No documents yet"
        description="Upload documents to get started"
        actionLabel="Upload Document"
        onActionPress={handleUploadDocument}
      />
    );
  };

  /**
   * Render loading state
   */
  if (isLoading && !documentData) {
    return (
      <View style={styles.centerContainer}>
        <LoadingSpinner size="large" message="Loading documents..." />
      </View>
    );
  }

  const hasActiveFilters =
    filterStatus !== 'ALL' || filterType !== 'ALL' || searchQuery.length > 0;

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Documents</Text>
        <TouchableOpacity
          style={styles.uploadButton}
          onPress={handleUploadDocument}
        >
          <Icon name="upload" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      {/* Search Bar */}
      <View style={styles.searchContainer}>
        <SearchBar
          value={searchQuery}
          onChangeText={setSearchQuery}
          placeholder="Search documents..."
          showFilter
          filterActive={hasActiveFilters}
          onFilterPress={() => setShowFilterSheet(true)}
        />
      </View>

      {/* Results Count */}
      <View style={styles.resultsContainer}>
        <Text style={styles.resultsText}>
          {filteredDocuments.length}{' '}
          {filteredDocuments.length === 1 ? 'document' : 'documents'}
        </Text>
        {hasActiveFilters && (
          <TouchableOpacity onPress={handleClearFilters}>
            <Text style={styles.clearText}>Clear All</Text>
          </TouchableOpacity>
        )}
      </View>

      {/* Documents List */}
      <FlatList
        data={filteredDocuments}
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
        title="Filter Documents"
        showHandle
        snapPoints={[0.7]}
      >
        <View style={styles.filterContent}>
          {/* Status Filters */}
          <Text style={styles.filterLabel}>Status</Text>
          <FilterOption
            label="All Documents"
            count={filterCounts.status.ALL}
            isSelected={filterStatus === 'ALL'}
            onPress={() => setFilterStatus('ALL')}
          />
          <FilterOption
            label="Pending"
            count={filterCounts.status.PENDING}
            isSelected={filterStatus === 'PENDING'}
            onPress={() => setFilterStatus('PENDING')}
            color={colors.warning}
          />
          <FilterOption
            label="Verified"
            count={filterCounts.status.VERIFIED}
            isSelected={filterStatus === 'VERIFIED'}
            onPress={() => setFilterStatus('VERIFIED')}
            color={colors.success}
          />
          <FilterOption
            label="Rejected"
            count={filterCounts.status.REJECTED}
            isSelected={filterStatus === 'REJECTED'}
            onPress={() => setFilterStatus('REJECTED')}
            color={colors.error}
          />

          {/* Type Filters */}
          <Text style={[styles.filterLabel, styles.filterLabelSpaced]}>Document Type</Text>
          <FilterOption
            label="All Types"
            count={filterCounts.type.ALL}
            isSelected={filterType === 'ALL'}
            onPress={() => setFilterType('ALL')}
          />
          <FilterOption
            label="PDF"
            count={filterCounts.type.PDF}
            isSelected={filterType === 'PDF'}
            onPress={() => setFilterType('PDF')}
            color="#E74C3C"
          />
          <FilterOption
            label="Images"
            count={filterCounts.type.IMAGE}
            isSelected={filterType === 'IMAGE'}
            onPress={() => setFilterType('IMAGE')}
            color="#9B59B6"
          />
          <FilterOption
            label="Documents"
            count={filterCounts.type.DOCUMENT}
            isSelected={filterType === 'DOCUMENT'}
            onPress={() => setFilterType('DOCUMENT')}
            color="#2E86DE"
          />
          <FilterOption
            label="Spreadsheets"
            count={filterCounts.type.SPREADSHEET}
            isSelected={filterType === 'SPREADSHEET'}
            onPress={() => setFilterType('SPREADSHEET')}
            color="#1E824C"
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

      {/* Upload Bottom Sheet */}
      <BottomSheet
        visible={showUploadSheet}
        onClose={() => setShowUploadSheet(false)}
        title="Upload Document"
        showHandle
        snapPoints={[0.4]}
      >
        <View style={styles.uploadContent}>
          <UploadOption
            icon="camera"
            label="Take Photo"
            description="Use camera to capture document"
            onPress={handleUploadFromCamera}
          />
          <UploadOption
            icon="image"
            label="Choose from Gallery"
            description="Select from photo library"
            onPress={handleUploadFromGallery}
          />
          <UploadOption
            icon="file-document"
            label="Browse Files"
            description="Select file from device"
            onPress={handleUploadFromFiles}
          />
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

/**
 * Upload Option Component
 */
interface UploadOptionProps {
  icon: string;
  label: string;
  description: string;
  onPress: () => void;
}

const UploadOption: React.FC<UploadOptionProps> = ({
  icon,
  label,
  description,
  onPress,
}) => (
  <TouchableOpacity style={styles.uploadOption} onPress={onPress}>
    <View style={styles.uploadOptionIcon}>
      <Icon name={icon} size={24} color={colors.primary} />
    </View>
    <View style={styles.uploadOptionText}>
      <Text style={styles.uploadOptionLabel}>{label}</Text>
      <Text style={styles.uploadOptionDescription}>{description}</Text>
    </View>
    <Icon name="chevron-right" size={20} color={colors.textTertiary} />
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
  uploadButton: {
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
  filterLabelSpaced: {
    marginTop: spacing.xl,
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
  uploadContent: {
    paddingBottom: spacing.lg,
  },
  uploadOption: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: spacing.md,
    borderRadius: 8,
    marginBottom: spacing.sm,
    backgroundColor: colors.background,
  },
  uploadOptionIcon: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: `${colors.primary}15`,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: spacing.md,
  },
  uploadOptionText: {
    flex: 1,
  },
  uploadOptionLabel: {
    fontSize: fonts.base,
    fontWeight: '600',
    color: colors.text,
    marginBottom: 2,
  },
  uploadOptionDescription: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
  },
});

export default DocumentsScreen;
