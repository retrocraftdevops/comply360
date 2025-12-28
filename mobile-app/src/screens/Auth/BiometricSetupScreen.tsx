/**
 * Biometric Setup Screen
 * Allows users to enable fingerprint/Face ID authentication
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { Button } from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

import { useAppDispatch, useAppSelector } from '../../store/store';
import { enableBiometric } from '../../store/slices/authSlice';
import { BiometricService } from '../../services/biometrics';
import { AuthStackParamList } from '../../navigation/AuthNavigator';

type BiometricSetupNavigationProp = StackNavigationProp<AuthStackParamList, 'BiometricSetup'>;

const BiometricSetupScreen = () => {
  const navigation = useNavigation<BiometricSetupNavigationProp>();
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);

  const [biometricType, setBiometricType] = useState<string>('');
  const [isEnrolling, setIsEnrolling] = useState(false);
  const [setupComplete, setSetupComplete] = useState(false);

  useEffect(() => {
    checkBiometricType();
  }, []);

  const checkBiometricType = async () => {
    const type = await BiometricService.getBiometricType();
    setBiometricType(type || 'Biometric');
  };

  const handleEnableBiometric = async () => {
    setIsEnrolling(true);

    try {
      // Check if biometrics are available
      const available = await BiometricService.isBiometricAvailable();

      if (!available) {
        Alert.alert(
          'Biometrics Not Available',
          `${biometricType} authentication is not available on this device. Please ensure it's enabled in your device settings.`,
          [{ text: 'OK' }]
        );
        setIsEnrolling(false);
        return;
      }

      // Prompt user to authenticate with biometrics
      const authenticated = await BiometricService.authenticate(
        `Enable ${biometricType}`,
        `Authenticate to enable ${biometricType} login for Comply360`
      );

      if (authenticated) {
        // Store credentials securely
        // Note: In production, you'd get email/password from previous login
        // For now, we'll just enable the biometric flag
        const stored = await BiometricService.storeCredentials(
          user?.email || '',
          '' // Password should come from previous login context
        );

        if (stored) {
          // Update Redux state
          dispatch(enableBiometric());
          setSetupComplete(true);

          // Show success message
          Alert.alert(
            'Success!',
            `${biometricType} authentication has been enabled. You can now use ${biometricType} to sign in.`,
            [
              {
                text: 'Continue',
                onPress: () => navigation.navigate('Login'),
              },
            ]
          );
        } else {
          throw new Error('Failed to store credentials securely');
        }
      }
    } catch (error: any) {
      console.error('Biometric setup error:', error);
      Alert.alert(
        'Setup Failed',
        error.message || 'Failed to enable biometric authentication. Please try again.',
        [{ text: 'OK' }]
      );
    } finally {
      setIsEnrolling(false);
    }
  };

  const handleSkip = () => {
    Alert.alert(
      'Skip Biometric Setup',
      'You can enable biometric authentication later from your profile settings.',
      [
        { text: 'Go Back', style: 'cancel' },
        {
          text: 'Skip',
          style: 'destructive',
          onPress: () => navigation.navigate('Login'),
        },
      ]
    );
  };

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.contentContainer}>
      {/* Header */}
      <View style={styles.header}>
        <View style={styles.iconContainer}>
          <Icon
            name={biometricType.toLowerCase().includes('face') ? 'face-recognition' : 'fingerprint'}
            size={80}
            color="#7c3aed"
          />
        </View>
        <Text style={styles.title}>
          {setupComplete ? 'Setup Complete!' : `Enable ${biometricType}`}
        </Text>
        <Text style={styles.subtitle}>
          {setupComplete
            ? `${biometricType} authentication is now enabled`
            : `Use ${biometricType} for quick and secure access to your account`}
        </Text>
      </View>

      {/* Benefits Section */}
      <View style={styles.benefitsContainer}>
        <Text style={styles.benefitsTitle}>Benefits</Text>

        <View style={styles.benefitItem}>
          <View style={styles.benefitIconContainer}>
            <Icon name="lightning-bolt" size={24} color="#7c3aed" />
          </View>
          <View style={styles.benefitTextContainer}>
            <Text style={styles.benefitTitle}>Faster Login</Text>
            <Text style={styles.benefitDescription}>
              Sign in instantly without typing your password
            </Text>
          </View>
        </View>

        <View style={styles.benefitItem}>
          <View style={styles.benefitIconContainer}>
            <Icon name="shield-check" size={24} color="#7c3aed" />
          </View>
          <View style={styles.benefitTextContainer}>
            <Text style={styles.benefitTitle}>Enhanced Security</Text>
            <Text style={styles.benefitDescription}>
              Your biometric data never leaves your device
            </Text>
          </View>
        </View>

        <View style={styles.benefitItem}>
          <View style={styles.benefitIconContainer}>
            <Icon name="lock" size={24} color="#7c3aed" />
          </View>
          <View style={styles.benefitTextContainer}>
            <Text style={styles.benefitTitle}>Secure Storage</Text>
            <Text style={styles.benefitDescription}>
              Credentials stored in device secure enclave
            </Text>
          </View>
        </View>

        <View style={styles.benefitItem}>
          <View style={styles.benefitIconContainer}>
            <Icon name="cog" size={24} color="#7c3aed" />
          </View>
          <View style={styles.benefitTextContainer}>
            <Text style={styles.benefitTitle}>Easy Management</Text>
            <Text style={styles.benefitDescription}>
              Disable anytime from your profile settings
            </Text>
          </View>
        </View>
      </View>

      {/* Privacy Notice */}
      <View style={styles.privacyContainer}>
        <Icon name="information" size={20} color="#6b7280" />
        <Text style={styles.privacyText}>
          Your biometric data is processed locally on your device and is never sent to our servers.
          We store an encrypted token in your device's secure storage.
        </Text>
      </View>

      {/* Action Buttons */}
      <View style={styles.actionsContainer}>
        {!setupComplete && (
          <>
            <Button
              mode="contained"
              onPress={handleEnableBiometric}
              loading={isEnrolling}
              disabled={isEnrolling}
              style={styles.enableButton}
              contentStyle={styles.buttonContent}
              labelStyle={styles.buttonLabel}
            >
              {isEnrolling ? 'Setting Up...' : `Enable ${biometricType}`}
            </Button>

            <TouchableOpacity
              onPress={handleSkip}
              disabled={isEnrolling}
              style={styles.skipButton}
            >
              <Text style={styles.skipButtonText}>Skip for Now</Text>
            </TouchableOpacity>
          </>
        )}

        {setupComplete && (
          <Button
            mode="contained"
            onPress={() => navigation.navigate('Login')}
            style={styles.enableButton}
            contentStyle={styles.buttonContent}
            labelStyle={styles.buttonLabel}
          >
            Continue to Login
          </Button>
        )}
      </View>

      {/* Requirements Section */}
      <View style={styles.requirementsContainer}>
        <Text style={styles.requirementsTitle}>Requirements</Text>
        <View style={styles.requirementItem}>
          <Icon name="check-circle" size={16} color="#10b981" />
          <Text style={styles.requirementText}>
            {biometricType} must be enabled on your device
          </Text>
        </View>
        <View style={styles.requirementItem}>
          <Icon name="check-circle" size={16} color="#10b981" />
          <Text style={styles.requirementText}>
            At least one {biometricType.toLowerCase()} enrolled
          </Text>
        </View>
        <View style={styles.requirementItem}>
          <Icon name="check-circle" size={16} color="#10b981" />
          <Text style={styles.requirementText}>Device must have secure lock screen</Text>
        </View>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  contentContainer: {
    padding: 24,
  },
  header: {
    alignItems: 'center',
    marginBottom: 32,
  },
  iconContainer: {
    width: 140,
    height: 140,
    borderRadius: 70,
    backgroundColor: '#f3e8ff',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 8,
    textAlign: 'center',
  },
  subtitle: {
    fontSize: 16,
    color: '#6b7280',
    textAlign: 'center',
    lineHeight: 22,
    paddingHorizontal: 16,
  },
  benefitsContainer: {
    marginBottom: 24,
  },
  benefitsTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 16,
  },
  benefitItem: {
    flexDirection: 'row',
    marginBottom: 20,
  },
  benefitIconContainer: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: '#f3e8ff',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  benefitTextContainer: {
    flex: 1,
    justifyContent: 'center',
  },
  benefitTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 2,
  },
  benefitDescription: {
    fontSize: 14,
    color: '#6b7280',
    lineHeight: 20,
  },
  privacyContainer: {
    flexDirection: 'row',
    backgroundColor: '#f3f4f6',
    padding: 16,
    borderRadius: 12,
    marginBottom: 24,
  },
  privacyText: {
    flex: 1,
    fontSize: 13,
    color: '#6b7280',
    lineHeight: 18,
    marginLeft: 12,
  },
  actionsContainer: {
    marginBottom: 24,
  },
  enableButton: {
    marginBottom: 12,
    backgroundColor: '#7c3aed',
  },
  buttonContent: {
    paddingVertical: 8,
  },
  buttonLabel: {
    fontSize: 16,
    fontWeight: '600',
  },
  skipButton: {
    paddingVertical: 12,
    alignItems: 'center',
  },
  skipButtonText: {
    fontSize: 16,
    color: '#6b7280',
    fontWeight: '600',
  },
  requirementsContainer: {
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
    paddingTop: 20,
  },
  requirementsTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 12,
  },
  requirementItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  requirementText: {
    fontSize: 14,
    color: '#6b7280',
    marginLeft: 8,
  },
});

export default BiometricSetupScreen;
