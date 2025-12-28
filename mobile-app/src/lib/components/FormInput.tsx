/**
 * FormInput Component
 * Validated input field with label, error, and helper text
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  TouchableOpacity,
  TextInputProps,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts } from '@/lib/utils/theme';

export interface FormInputProps extends TextInputProps {
  label: string;
  value: string;
  onChangeText: (text: string) => void;
  error?: string;
  helperText?: string;
  required?: boolean;
  disabled?: boolean;
  multiline?: boolean;
  numberOfLines?: number;
  icon?: string;
  rightIcon?: string;
  onRightIconPress?: () => void;
  validate?: (value: string) => string | null;
  validateOnBlur?: boolean;
  secureTextEntry?: boolean;
}

const FormInput: React.FC<FormInputProps> = ({
  label,
  value,
  onChangeText,
  error: externalError,
  helperText,
  required = false,
  disabled = false,
  multiline = false,
  numberOfLines = 1,
  icon,
  rightIcon,
  onRightIconPress,
  validate,
  validateOnBlur = true,
  secureTextEntry: initialSecureTextEntry = false,
  ...textInputProps
}) => {
  const [isFocused, setIsFocused] = useState(false);
  const [internalError, setInternalError] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  const error = externalError || internalError;
  const hasError = !!error;
  const isPasswordField = initialSecureTextEntry;
  const secureTextEntry = isPasswordField && !showPassword;

  /**
   * Handle focus
   */
  const handleFocus = () => {
    setIsFocused(true);
  };

  /**
   * Handle blur
   */
  const handleBlur = () => {
    setIsFocused(false);

    // Validate on blur if enabled
    if (validateOnBlur && validate) {
      const validationError = validate(value);
      setInternalError(validationError);
    }
  };

  /**
   * Handle text change
   */
  const handleChangeText = (text: string) => {
    onChangeText(text);

    // Clear error on change
    if (internalError) {
      setInternalError(null);
    }
  };

  /**
   * Toggle password visibility
   */
  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <View style={styles.container}>
      {/* Label */}
      <View style={styles.labelContainer}>
        <Text style={styles.label}>
          {label}
          {required && <Text style={styles.required}> *</Text>}
        </Text>
      </View>

      {/* Input Container */}
      <View
        style={[
          styles.inputContainer,
          isFocused && styles.inputContainerFocused,
          hasError && styles.inputContainerError,
          disabled && styles.inputContainerDisabled,
          multiline && styles.inputContainerMultiline,
        ]}
      >
        {/* Left Icon */}
        {icon && (
          <View style={styles.leftIcon}>
            <Icon
              name={icon}
              size={20}
              color={
                hasError
                  ? colors.error
                  : isFocused
                  ? colors.primary
                  : colors.textTertiary
              }
            />
          </View>
        )}

        {/* Text Input */}
        <TextInput
          style={[
            styles.input,
            multiline && styles.inputMultiline,
            icon && styles.inputWithLeftIcon,
            (rightIcon || isPasswordField) && styles.inputWithRightIcon,
          ]}
          value={value}
          onChangeText={handleChangeText}
          onFocus={handleFocus}
          onBlur={handleBlur}
          editable={!disabled}
          multiline={multiline}
          numberOfLines={multiline ? numberOfLines : 1}
          textAlignVertical={multiline ? 'top' : 'center'}
          placeholderTextColor={colors.textTertiary}
          secureTextEntry={secureTextEntry}
          {...textInputProps}
        />

        {/* Right Icon or Password Toggle */}
        {(rightIcon || isPasswordField) && (
          <TouchableOpacity
            style={styles.rightIcon}
            onPress={isPasswordField ? togglePasswordVisibility : onRightIconPress}
            disabled={disabled}
          >
            <Icon
              name={
                isPasswordField
                  ? showPassword
                    ? 'eye-off'
                    : 'eye'
                  : rightIcon || 'information'
              }
              size={20}
              color={
                hasError
                  ? colors.error
                  : isFocused
                  ? colors.primary
                  : colors.textTertiary
              }
            />
          </TouchableOpacity>
        )}
      </View>

      {/* Helper Text or Error */}
      {(helperText || error) && (
        <View style={styles.helperContainer}>
          {error ? (
            <View style={styles.errorContainer}>
              <Icon name="alert-circle" size={14} color={colors.error} />
              <Text style={styles.errorText}>{error}</Text>
            </View>
          ) : (
            <Text style={styles.helperText}>{helperText}</Text>
          )}
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.lg,
  },
  labelContainer: {
    marginBottom: spacing.xs,
  },
  label: {
    fontSize: fonts.sm,
    fontWeight: '600',
    color: colors.text,
  },
  required: {
    color: colors.error,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    paddingHorizontal: spacing.md,
    minHeight: 48,
  },
  inputContainerFocused: {
    borderColor: colors.primary,
    borderWidth: 2,
  },
  inputContainerError: {
    borderColor: colors.error,
  },
  inputContainerDisabled: {
    backgroundColor: colors.backgroundSecondary,
    opacity: 0.6,
  },
  inputContainerMultiline: {
    minHeight: 100,
    alignItems: 'flex-start',
    paddingVertical: spacing.md,
  },
  leftIcon: {
    marginRight: spacing.sm,
  },
  rightIcon: {
    marginLeft: spacing.sm,
    padding: spacing.xs,
  },
  input: {
    flex: 1,
    fontSize: fonts.base,
    color: colors.text,
    paddingVertical: 0,
    minHeight: 48,
  },
  inputMultiline: {
    minHeight: 100,
    paddingTop: spacing.xs,
  },
  inputWithLeftIcon: {
    marginLeft: 0,
  },
  inputWithRightIcon: {
    marginRight: 0,
  },
  helperContainer: {
    marginTop: spacing.xs,
    paddingHorizontal: spacing.sm,
  },
  helperText: {
    fontSize: fonts.xs,
    color: colors.textTertiary,
    lineHeight: 16,
  },
  errorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  errorText: {
    fontSize: fonts.xs,
    color: colors.error,
    marginLeft: spacing.xs,
    flex: 1,
    lineHeight: 16,
  },
});

export default FormInput;
