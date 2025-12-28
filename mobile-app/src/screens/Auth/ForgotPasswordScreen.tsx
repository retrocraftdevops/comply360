/**
 * Forgot Password Screen
 * Allows users to request a password reset via email
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  TouchableOpacity,
  ScrollView,
} from 'react-native';
import { TextInput, Button } from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

import { AuthService } from '../../services/auth';
import { AuthStackParamList } from '../../navigation/AuthNavigator';

type ForgotPasswordNavigationProp = StackNavigationProp<AuthStackParamList, 'ForgotPassword'>;

const ForgotPasswordScreen = () => {
  const navigation = useNavigation<ForgotPasswordNavigationProp>();

  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleResetPassword = async () => {
    // Reset state
    setError('');

    // Validate email
    if (!email) {
      setError('Please enter your email address');
      return;
    }

    if (!validateEmail(email)) {
      setError('Please enter a valid email address');
      return;
    }

    setIsLoading(true);

    try {
      await AuthService.requestPasswordReset(email);
      setSuccess(true);
    } catch (err: any) {
      console.error('Password reset error:', err);
      setError(err.message || 'Failed to send reset email. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleBackToLogin = () => {
    navigation.navigate('Login');
  };

  if (success) {
    return (
      <View style={styles.container}>
        <ScrollView contentContainerStyle={styles.successContainer}>
          {/* Success Icon */}
          <View style={styles.successIconContainer}>
            <Icon name="email-check" size={80} color="#10b981" />
          </View>

          {/* Success Message */}
          <Text style={styles.successTitle}>Check Your Email</Text>
          <Text style={styles.successMessage}>
            We've sent password reset instructions to:
          </Text>
          <Text style={styles.emailText}>{email}</Text>

          {/* Instructions */}
          <View style={styles.instructionsContainer}>
            <Text style={styles.instructionsTitle}>Next Steps:</Text>

            <View style={styles.instructionItem}>
              <View style={styles.stepNumber}>
                <Text style={styles.stepNumberText}>1</Text>
              </View>
              <Text style={styles.instructionText}>Check your email inbox</Text>
            </View>

            <View style={styles.instructionItem}>
              <View style={styles.stepNumber}>
                <Text style={styles.stepNumberText}>2</Text>
              </View>
              <Text style={styles.instructionText}>
                Click the password reset link (valid for 1 hour)
              </Text>
            </View>

            <View style={styles.instructionItem}>
              <View style={styles.stepNumber}>
                <Text style={styles.stepNumberText}>3</Text>
              </View>
              <Text style={styles.instructionText}>Create a new password</Text>
            </View>

            <View style={styles.instructionItem}>
              <View style={styles.stepNumber}>
                <Text style={styles.stepNumberText}>4</Text>
              </View>
              <Text style={styles.instructionText}>Sign in with your new password</Text>
            </View>
          </View>

          {/* Help Text */}
          <View style={styles.helpContainer}>
            <Icon name="information" size={20} color="#6b7280" />
            <Text style={styles.helpText}>
              Didn't receive the email? Check your spam folder or try again in a few minutes.
            </Text>
          </View>

          {/* Actions */}
          <View style={styles.actionsContainer}>
            <Button
              mode="contained"
              onPress={handleBackToLogin}
              style={styles.loginButton}
              contentStyle={styles.buttonContent}
              labelStyle={styles.buttonLabel}
            >
              Back to Login
            </Button>

            <TouchableOpacity
              onPress={() => {
                setSuccess(false);
                setEmail('');
              }}
              style={styles.resendButton}
            >
              <Text style={styles.resendButtonText}>Send Again</Text>
            </TouchableOpacity>
          </View>
        </ScrollView>
      </View>
    );
  }

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView contentContainerStyle={styles.contentContainer}>
        {/* Header with Back Button */}
        <TouchableOpacity onPress={handleBackToLogin} style={styles.backButton}>
          <Icon name="arrow-left" size={24} color="#111827" />
          <Text style={styles.backButtonText}>Back to Login</Text>
        </TouchableOpacity>

        {/* Icon */}
        <View style={styles.iconContainer}>
          <View style={styles.iconCircle}>
            <Icon name="lock-reset" size={60} color="#7c3aed" />
          </View>
        </View>

        {/* Title and Description */}
        <Text style={styles.title}>Forgot Password?</Text>
        <Text style={styles.description}>
          No worries! Enter your email address and we'll send you instructions to reset your
          password.
        </Text>

        {/* Email Input */}
        <View style={styles.formContainer}>
          <TextInput
            label="Email Address"
            value={email}
            onChangeText={(text) => {
              setEmail(text);
              setError('');
            }}
            mode="outlined"
            autoCapitalize="none"
            keyboardType="email-address"
            autoComplete="email"
            textContentType="emailAddress"
            left={<TextInput.Icon icon="email-outline" />}
            style={styles.input}
            error={!!error}
            disabled={isLoading}
          />

          {error && (
            <View style={styles.errorContainer}>
              <Icon name="alert-circle" size={16} color="#ef4444" />
              <Text style={styles.errorText}>{error}</Text>
            </View>
          )}

          <Button
            mode="contained"
            onPress={handleResetPassword}
            loading={isLoading}
            disabled={isLoading}
            style={styles.submitButton}
            contentStyle={styles.buttonContent}
            labelStyle={styles.buttonLabel}
          >
            {isLoading ? 'Sending...' : 'Send Reset Instructions'}
          </Button>
        </View>

        {/* Security Notice */}
        <View style={styles.securityContainer}>
          <Icon name="shield-check" size={20} color="#10b981" />
          <Text style={styles.securityText}>
            For your security, the password reset link will expire after 1 hour.
          </Text>
        </View>

        {/* Alternative Options */}
        <View style={styles.alternativesContainer}>
          <Text style={styles.alternativesTitle}>Need Help?</Text>

          <View style={styles.alternativeItem}>
            <Icon name="email" size={20} color="#6b7280" />
            <View style={styles.alternativeTextContainer}>
              <Text style={styles.alternativeLabel}>Email Support</Text>
              <Text style={styles.alternativeValue}>support@comply360.com</Text>
            </View>
          </View>

          <View style={styles.alternativeItem}>
            <Icon name="phone" size={20} color="#6b7280" />
            <View style={styles.alternativeTextContainer}>
              <Text style={styles.alternativeLabel}>Phone Support</Text>
              <Text style={styles.alternativeValue}>+27 11 123 4567</Text>
            </View>
          </View>

          <View style={styles.alternativeItem}>
            <Icon name="clock-outline" size={20} color="#6b7280" />
            <View style={styles.alternativeTextContainer}>
              <Text style={styles.alternativeLabel}>Support Hours</Text>
              <Text style={styles.alternativeValue}>Mon-Fri: 8:00 AM - 6:00 PM SAST</Text>
            </View>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
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
  backButton: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 24,
  },
  backButtonText: {
    fontSize: 16,
    color: '#111827',
    marginLeft: 8,
    fontWeight: '500',
  },
  iconContainer: {
    alignItems: 'center',
    marginBottom: 24,
  },
  iconCircle: {
    width: 120,
    height: 120,
    borderRadius: 60,
    backgroundColor: '#f3e8ff',
    justifyContent: 'center',
    alignItems: 'center',
  },
  title: {
    fontSize: 28,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 12,
    textAlign: 'center',
  },
  description: {
    fontSize: 16,
    color: '#6b7280',
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 32,
    paddingHorizontal: 8,
  },
  formContainer: {
    marginBottom: 24,
  },
  input: {
    marginBottom: 16,
    backgroundColor: '#FFFFFF',
  },
  errorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
    paddingHorizontal: 12,
    paddingVertical: 8,
    backgroundColor: '#fee2e2',
    borderRadius: 8,
  },
  errorText: {
    marginLeft: 8,
    color: '#ef4444',
    fontSize: 14,
  },
  submitButton: {
    backgroundColor: '#7c3aed',
  },
  buttonContent: {
    paddingVertical: 8,
  },
  buttonLabel: {
    fontSize: 16,
    fontWeight: '600',
  },
  securityContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f0fdf4',
    padding: 16,
    borderRadius: 12,
    marginBottom: 24,
  },
  securityText: {
    flex: 1,
    fontSize: 13,
    color: '#047857',
    lineHeight: 18,
    marginLeft: 12,
  },
  alternativesContainer: {
    borderTopWidth: 1,
    borderTopColor: '#e5e7eb',
    paddingTop: 24,
  },
  alternativesTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 16,
  },
  alternativeItem: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  alternativeTextContainer: {
    marginLeft: 12,
  },
  alternativeLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 2,
  },
  alternativeValue: {
    fontSize: 14,
    color: '#6b7280',
  },
  // Success screen styles
  successContainer: {
    padding: 24,
    alignItems: 'center',
  },
  successIconContainer: {
    width: 140,
    height: 140,
    borderRadius: 70,
    backgroundColor: '#f0fdf4',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 24,
    marginTop: 40,
  },
  successTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 12,
    textAlign: 'center',
  },
  successMessage: {
    fontSize: 16,
    color: '#6b7280',
    textAlign: 'center',
    marginBottom: 8,
  },
  emailText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#7c3aed',
    textAlign: 'center',
    marginBottom: 32,
  },
  instructionsContainer: {
    width: '100%',
    marginBottom: 24,
  },
  instructionsTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111827',
    marginBottom: 16,
  },
  instructionItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: 16,
  },
  stepNumber: {
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: '#7c3aed',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  stepNumberText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  instructionText: {
    flex: 1,
    fontSize: 15,
    color: '#374151',
    lineHeight: 22,
    paddingTop: 3,
  },
  helpContainer: {
    flexDirection: 'row',
    backgroundColor: '#f3f4f6',
    padding: 16,
    borderRadius: 12,
    marginBottom: 24,
    width: '100%',
  },
  helpText: {
    flex: 1,
    fontSize: 13,
    color: '#6b7280',
    lineHeight: 18,
    marginLeft: 12,
  },
  actionsContainer: {
    width: '100%',
  },
  loginButton: {
    marginBottom: 12,
    backgroundColor: '#7c3aed',
  },
  resendButton: {
    paddingVertical: 12,
    alignItems: 'center',
  },
  resendButtonText: {
    fontSize: 16,
    color: '#7c3aed',
    fontWeight: '600',
  },
});

export default ForgotPasswordScreen;
