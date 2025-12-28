/**
 * Modal Component
 * Reusable modal/dialog component
 */

import React, { ReactNode } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal as RNModal,
  TouchableOpacity,
  ScrollView,
  Dimensions,
  TouchableWithoutFeedback,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import Button from './Button';

const { height } = Dimensions.get('window');

export interface ModalProps {
  visible: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  footer?: ReactNode;
  size?: 'small' | 'medium' | 'large' | 'full';
  showCloseButton?: boolean;
  closeOnBackdropPress?: boolean;
  primaryButton?: {
    title: string;
    onPress: () => void;
    loading?: boolean;
    variant?: 'primary' | 'danger';
  };
  secondaryButton?: {
    title: string;
    onPress: () => void;
  };
}

const Modal: React.FC<ModalProps> = ({
  visible,
  onClose,
  title,
  children,
  footer,
  size = 'medium',
  showCloseButton = true,
  closeOnBackdropPress = true,
  primaryButton,
  secondaryButton,
}) => {
  const getModalHeight = (): number | string => {
    switch (size) {
      case 'small':
        return height * 0.3;
      case 'medium':
        return height * 0.5;
      case 'large':
        return height * 0.7;
      case 'full':
        return height * 0.9;
      default:
        return height * 0.5;
    }
  };

  const handleBackdropPress = () => {
    if (closeOnBackdropPress) {
      onClose();
    }
  };

  return (
    <RNModal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <TouchableWithoutFeedback onPress={handleBackdropPress}>
        <View style={styles.backdrop}>
          <TouchableWithoutFeedback>
            <View style={[styles.modal, { maxHeight: getModalHeight() }]}>
              {/* Header */}
              {(title || showCloseButton) && (
                <View style={styles.header}>
                  {title && <Text style={styles.title}>{title}</Text>}
                  {showCloseButton && (
                    <TouchableOpacity
                      onPress={onClose}
                      style={styles.closeButton}
                    >
                      <Icon name="close" size={24} color="#6b7280" />
                    </TouchableOpacity>
                  )}
                </View>
              )}

              {/* Content */}
              <ScrollView
                style={styles.content}
                contentContainerStyle={styles.contentContainer}
                showsVerticalScrollIndicator={true}
              >
                {children}
              </ScrollView>

              {/* Footer */}
              {(footer || primaryButton || secondaryButton) && (
                <View style={styles.footer}>
                  {footer ? (
                    footer
                  ) : (
                    <View style={styles.footerButtons}>
                      {secondaryButton && (
                        <Button
                          title={secondaryButton.title}
                          onPress={secondaryButton.onPress}
                          variant="outline"
                          style={styles.footerButton}
                        />
                      )}
                      {primaryButton && (
                        <Button
                          title={primaryButton.title}
                          onPress={primaryButton.onPress}
                          variant={primaryButton.variant || 'primary'}
                          loading={primaryButton.loading}
                          style={styles.footerButton}
                        />
                      )}
                    </View>
                  )}
                </View>
              )}
            </View>
          </TouchableWithoutFeedback>
        </View>
      </TouchableWithoutFeedback>
    </RNModal>
  );
};

const styles = StyleSheet.create({
  backdrop: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  modal: {
    backgroundColor: '#FFFFFF',
    borderRadius: 16,
    width: '100%',
    maxWidth: 500,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.3,
    shadowRadius: 16,
    elevation: 10,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 24,
    paddingVertical: 20,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  title: {
    flex: 1,
    fontSize: 20,
    fontWeight: '700',
    color: '#111827',
  },
  closeButton: {
    padding: 4,
    marginLeft: 16,
  },
  content: {
    maxHeight: height * 0.6,
  },
  contentContainer: {
    padding: 24,
  },
  footer: {
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
    paddingHorizontal: 24,
    paddingVertical: 16,
  },
  footerButtons: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  footerButton: {
    marginLeft: 12,
    minWidth: 100,
  },
});

export default Modal;
