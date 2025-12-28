/**
 * DocumentCard Component
 * Display card for document list items
 */

import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts, shadows } from '@/lib/utils/theme';
import { formatDate, formatFileSize } from '@/lib/utils/formatting';

export interface Document {
  id: string;
  name: string;
  type: string;
  size: number;
  status: 'PENDING' | 'VERIFIED' | 'REJECTED';
  uploaded_at: string;
  registration_id?: string;
  company_name?: string;
}

export interface DocumentCardProps {
  document: Document;
  onPress: (document: Document) => void;
  onDownload?: (document: Document) => void;
}

const DocumentCard: React.FC<DocumentCardProps> = ({
  document,
  onPress,
  onDownload,
}) => {
  /**
   * Get file icon based on type
   */
  const getFileIcon = (type: string): string => {
    const lowerType = type.toLowerCase();
    if (lowerType.includes('pdf')) return 'file-pdf-box';
    if (lowerType.includes('doc')) return 'file-word-box';
    if (lowerType.includes('xls') || lowerType.includes('sheet')) return 'file-excel-box';
    if (lowerType.includes('image') || lowerType.includes('png') || lowerType.includes('jpg')) return 'file-image';
    return 'file-document';
  };

  /**
   * Get file icon color
   */
  const getFileIconColor = (type: string): string => {
    const lowerType = type.toLowerCase();
    if (lowerType.includes('pdf')) return '#E74C3C';
    if (lowerType.includes('doc')) return '#2E86DE';
    if (lowerType.includes('xls') || lowerType.includes('sheet')) return '#1E824C';
    if (lowerType.includes('image')) return '#9B59B6';
    return colors.textSecondary;
  };

  /**
   * Get status color
   */
  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'VERIFIED':
        return colors.success;
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
      case 'VERIFIED':
        return 'check-circle';
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
      onPress={() => onPress(document)}
      activeOpacity={0.7}
    >
      <View style={styles.content}>
        {/* File Icon */}
        <View
          style={[
            styles.iconContainer,
            { backgroundColor: `${getFileIconColor(document.type)}15` },
          ]}
        >
          <Icon
            name={getFileIcon(document.type)}
            size={32}
            color={getFileIconColor(document.type)}
          />
        </View>

        {/* Document Info */}
        <View style={styles.info}>
          <Text style={styles.name} numberOfLines={1}>
            {document.name}
          </Text>

          {document.company_name && (
            <Text style={styles.companyName} numberOfLines={1}>
              {document.company_name}
            </Text>
          )}

          <View style={styles.metadata}>
            <Text style={styles.metadataText}>{formatFileSize(document.size)}</Text>
            <Text style={styles.metadataDot}>•</Text>
            <Text style={styles.metadataText}>
              {formatDate(new Date(document.uploaded_at), 'short')}
            </Text>
          </View>

          {/* Status Badge */}
          <View
            style={[
              styles.statusBadge,
              { backgroundColor: `${getStatusColor(document.status)}15` },
            ]}
          >
            <Icon
              name={getStatusIcon(document.status)}
              size={12}
              color={getStatusColor(document.status)}
            />
            <Text
              style={[
                styles.statusText,
                { color: getStatusColor(document.status) },
              ]}
            >
              {document.status}
            </Text>
          </View>
        </View>

        {/* Actions */}
        <View style={styles.actions}>
          {onDownload && (
            <TouchableOpacity
              style={styles.actionButton}
              onPress={() => onDownload(document)}
              hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
            >
              <Icon name="download" size={20} color={colors.primary} />
            </TouchableOpacity>
          )}
          <Icon name="chevron-right" size={20} color={colors.textTertiary} />
        </View>
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
  content: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  iconContainer: {
    width: 56,
    height: 56,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: spacing.md,
  },
  info: {
    flex: 1,
  },
  name: {
    fontSize: fonts.base,
    fontWeight: '600',
    color: colors.text,
    marginBottom: 4,
  },
  companyName: {
    fontSize: fonts.sm,
    color: colors.textSecondary,
    marginBottom: 4,
  },
  metadata: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  metadataText: {
    fontSize: fonts.xs,
    color: colors.textTertiary,
  },
  metadataDot: {
    fontSize: fonts.xs,
    color: colors.textTertiary,
    marginHorizontal: spacing.xs,
  },
  statusBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    paddingHorizontal: spacing.sm,
    paddingVertical: 4,
    borderRadius: 10,
  },
  statusText: {
    fontSize: 10,
    fontWeight: '600',
    marginLeft: 4,
  },
  actions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  actionButton: {
    padding: spacing.xs,
  },
});

export default DocumentCard;
