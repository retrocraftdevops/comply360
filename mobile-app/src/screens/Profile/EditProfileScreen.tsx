/**
 * EditProfileScreen
 * User profile editing with form validation
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  Alert,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { Avatar, FormInput, Button, Card } from '@/lib/components';
import { useAppDispatch, useAppSelector } from '@/store/store';
import { updateUser } from '@/store/slices/authSlice';
import { showToast } from '@/store/slices/uiSlice';
import { validateEmail, validatePhoneNumber } from '@/lib/utils/validation';
import { colors, spacing } from '@/lib/utils/theme';
import { useNavigation } from '@react-navigation/native';

const EditProfileScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.auth.user);

  const [name, setName] = useState(user?.name || '');
  const [email, setEmail] = useState(user?.email || '');
  const [phone, setPhone] = useState(user?.phone || '');
  const [company, setCompany] = useState(user?.company || '');
  const [isLoading, setIsLoading] = useState(false);

  const [errors, setErrors] = useState<{
    name?: string;
    email?: string;
    phone?: string;
    company?: string;
  }>({});

  const handleChangeAvatar = () => {
    Alert.alert(
      'Change Avatar',
      'Choose an option',
      [
        {
          text: 'Take Photo',
          onPress: () => dispatch(showToast({ message: 'Camera coming soon!', type: 'info' })),
        },
        {
          text: 'Choose from Gallery',
          onPress: () => dispatch(showToast({ message: 'Gallery coming soon!', type: 'info' })),
        },
        {
          text: 'Cancel',
          style: 'cancel',
        },
      ],
      { cancelable: true }
    );
  };

  const validateForm = (): boolean => {
    const newErrors: typeof errors = {};

    // Name validation
    if (!name.trim()) {
      newErrors.name = 'Name is required';
    } else if (name.trim().length < 2) {
      newErrors.name = 'Name must be at least 2 characters';
    }

    // Email validation
    if (!email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!validateEmail(email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    // Phone validation (optional)
    if (phone.trim() && !validatePhoneNumber(phone)) {
      newErrors.phone = 'Please enter a valid phone number';
    }

    // Company validation (optional)
    if (company.trim() && company.trim().length < 2) {
      newErrors.company = 'Company name must be at least 2 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSave = async () => {
    if (!validateForm()) {
      dispatch(showToast({ message: 'Please fix the errors in the form', type: 'error' }));
      return;
    }

    setIsLoading(true);

    try {
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1500));

      // Update user in Redux store
      dispatch(
        updateUser({
          name: name.trim(),
          email: email.trim(),
          phone: phone.trim() || undefined,
          company: company.trim() || undefined,
        })
      );

      dispatch(showToast({ message: 'Profile updated successfully!', type: 'success' }));
      navigation.goBack();
    } catch (error) {
      console.error('Failed to update profile:', error);
      dispatch(showToast({ message: 'Failed to update profile', type: 'error' }));
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    navigation.goBack();
  };

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
    >
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Avatar Section */}
        <View style={styles.avatarSection}>
          <Avatar name={name || user?.name || 'User'} size="xlarge" editable />
          <TouchableOpacity style={styles.changeAvatarButton} onPress={handleChangeAvatar}>
            <Icon name="camera" size={20} color={colors.primary} />
            <Text style={styles.changeAvatarText}>Change Avatar</Text>
          </TouchableOpacity>
        </View>

        {/* Form Section */}
        <View style={styles.formSection}>
          <Text style={styles.sectionTitle}>Personal Information</Text>
          <Card variant="outlined" padding="none">
            <View style={styles.formContainer}>
              <FormInput
                label="Full Name"
                value={name}
                onChangeText={setName}
                placeholder="Enter your full name"
                icon="account"
                required
                error={errors.name}
              />

              <FormInput
                label="Email Address"
                value={email}
                onChangeText={setEmail}
                placeholder="your.email@example.com"
                keyboardType="email-address"
                autoCapitalize="none"
                icon="email"
                required
                error={errors.email}
              />

              <FormInput
                label="Phone Number"
                value={phone}
                onChangeText={setPhone}
                placeholder="0XX XXX XXXX"
                keyboardType="phone-pad"
                icon="phone"
                helperText="Optional - South African format"
                error={errors.phone}
              />

              <FormInput
                label="Company"
                value={company}
                onChangeText={setCompany}
                placeholder="Your company name"
                icon="domain"
                helperText="Optional"
                error={errors.company}
              />
            </View>
          </Card>
        </View>

        {/* Account Information */}
        <View style={styles.formSection}>
          <Text style={styles.sectionTitle}>Account Information</Text>
          <Card variant="outlined" padding="medium">
            <View style={styles.infoRow}>
              <Icon name="shield-account" size={20} color="#6b7280" />
              <View style={styles.infoContent}>
                <Text style={styles.infoLabel}>Role</Text>
                <Text style={styles.infoValue}>Agent</Text>
              </View>
            </View>
            <View style={[styles.infoRow, { marginTop: 12 }]}>
              <Icon name="calendar-clock" size={20} color="#6b7280" />
              <View style={styles.infoContent}>
                <Text style={styles.infoLabel}>Member Since</Text>
                <Text style={styles.infoValue}>December 2025</Text>
              </View>
            </View>
          </Card>
        </View>

        {/* Action Buttons */}
        <View style={styles.buttonSection}>
          <Button
            title="Save Changes"
            onPress={handleSave}
            variant="primary"
            icon="content-save"
            loading={isLoading}
            fullWidth
          />
          <View style={{ height: 12 }} />
          <Button
            title="Cancel"
            onPress={handleCancel}
            variant="outline"
            icon="close"
            disabled={isLoading}
            fullWidth
          />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  avatarSection: {
    backgroundColor: '#FFFFFF',
    alignItems: 'center',
    paddingTop: 24,
    paddingBottom: 20,
    paddingHorizontal: 20,
  },
  changeAvatarButton: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 12,
    paddingVertical: 8,
    paddingHorizontal: 16,
    backgroundColor: `${colors.primary}10`,
    borderRadius: 20,
  },
  changeAvatarText: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.primary,
    marginLeft: 6,
  },
  formSection: {
    paddingHorizontal: 20,
    paddingTop: 24,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 12,
  },
  formContainer: {
    padding: 16,
  },
  infoRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  infoContent: {
    marginLeft: 12,
    flex: 1,
  },
  infoLabel: {
    fontSize: 12,
    color: '#6b7280',
    marginBottom: 2,
  },
  infoValue: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
  },
  buttonSection: {
    paddingHorizontal: 20,
    paddingTop: 24,
    paddingBottom: 32,
  },
});

export default EditProfileScreen;
