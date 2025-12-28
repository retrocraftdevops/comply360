/**
 * SettingsScreen
 * App settings and preferences management
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Switch,
  TouchableOpacity,
  Alert,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { Card } from '@/lib/components';
import { useAppDispatch } from '@/store/store';
import { showToast } from '@/store/slices/uiSlice';
import { useTheme } from '@/contexts/ThemeContext';

const SettingsScreen: React.FC = () => {
  const dispatch = useAppDispatch();
  const { theme, isDark, toggleTheme } = useTheme();

  // Settings state
  const [pushNotifications, setPushNotifications] = useState(true);
  const [emailNotifications, setEmailNotifications] = useState(true);
  const [smsNotifications, setSmsNotifications] = useState(false);
  const [biometricLogin, setBiometricLogin] = useState(false);

  const handleDarkModeToggle = () => {
    toggleTheme();
    dispatch(
      showToast({
        message: isDark ? 'Light mode enabled' : 'Dark mode enabled',
        type: 'info',
      })
    );
  };

  const handlePushNotificationsToggle = (value: boolean) => {
    setPushNotifications(value);
    dispatch(
      showToast({
        message: value ? 'Push notifications enabled' : 'Push notifications disabled',
        type: 'info',
      })
    );
  };

  const handleEmailNotificationsToggle = (value: boolean) => {
    setEmailNotifications(value);
    dispatch(
      showToast({
        message: value ? 'Email notifications enabled' : 'Email notifications disabled',
        type: 'info',
      })
    );
  };

  const handleSmsNotificationsToggle = (value: boolean) => {
    setSmsNotifications(value);
    dispatch(
      showToast({
        message: value ? 'SMS notifications enabled' : 'SMS notifications disabled',
        type: 'info',
      })
    );
  };

  const handleBiometricLoginToggle = (value: boolean) => {
    setBiometricLogin(value);
    dispatch(
      showToast({
        message: value ? 'Biometric login enabled' : 'Biometric login disabled',
        type: 'info',
      })
    );
  };

  const handleLanguageChange = () => {
    dispatch(showToast({ message: 'Language selection coming soon!', type: 'info' }));
  };

  const handleClearCache = () => {
    Alert.alert(
      'Clear Cache',
      'This will clear all cached data. Are you sure?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Clear',
          style: 'destructive',
          onPress: () => {
            dispatch(showToast({ message: 'Cache cleared successfully', type: 'success' }));
          },
        },
      ],
      { cancelable: true }
    );
  };

  const handleAbout = () => {
    Alert.alert(
      'About Comply360',
      'Version: 1.0.0\n\n' +
        'SADC Corporate Gateway Platform\n' +
        'Mobile Application for iOS and Android\n\n' +
        'Copyright © 2025 Comply360\n' +
        'All rights reserved.',
      [{ text: 'OK' }]
    );
  };

  const handlePrivacyPolicy = () => {
    dispatch(showToast({ message: 'Privacy policy coming soon!', type: 'info' }));
  };

  const handleTermsOfService = () => {
    dispatch(showToast({ message: 'Terms of service coming soon!', type: 'info' }));
  };

  const handleHelp = () => {
    dispatch(showToast({ message: 'Help center coming soon!', type: 'info' }));
  };

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      {/* Appearance Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Appearance</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="theme-light-dark"
            iconColor="#9f7aea"
            label="Dark Mode"
            subtitle="Switch between light and dark theme"
            rightComponent={
              <Switch
                value={isDark}
                onValueChange={handleDarkModeToggle}
                trackColor={{ false: '#d1d5db', true: theme.colors.primary }}
                thumbColor="#FFFFFF"
              />
            }
          />
        </Card>
      </View>

      {/* Notifications Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Notifications</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="bell"
            iconColor="#3b82f6"
            label="Push Notifications"
            subtitle="Receive push notifications"
            rightComponent={
              <Switch
                value={pushNotifications}
                onValueChange={handlePushNotificationsToggle}
                trackColor={{ false: '#d1d5db', true: theme.colors.primary }}
                thumbColor="#FFFFFF"
              />
            }
          />
          <SettingRow
            icon="email"
            iconColor="#10b981"
            label="Email Notifications"
            subtitle="Receive email updates"
            rightComponent={
              <Switch
                value={emailNotifications}
                onValueChange={handleEmailNotificationsToggle}
                trackColor={{ false: '#d1d5db', true: theme.colors.primary }}
                thumbColor="#FFFFFF"
              />
            }
          />
          <SettingRow
            icon="message-text"
            iconColor="#f59e0b"
            label="SMS Notifications"
            subtitle="Receive SMS alerts"
            showBorder={false}
            rightComponent={
              <Switch
                value={smsNotifications}
                onValueChange={handleSmsNotificationsToggle}
                trackColor={{ false: '#d1d5db', true: theme.colors.primary }}
                thumbColor="#FFFFFF"
              />
            }
          />
        </Card>
      </View>

      {/* Security Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Security</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="fingerprint"
            iconColor="#ef4444"
            label="Biometric Login"
            subtitle="Use Touch ID or Face ID"
            rightComponent={
              <Switch
                value={biometricLogin}
                onValueChange={handleBiometricLoginToggle}
                trackColor={{ false: '#d1d5db', true: theme.colors.primary }}
                thumbColor="#FFFFFF"
              />
            }
            showBorder={false}
          />
        </Card>
      </View>

      {/* Preferences Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Preferences</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="translate"
            iconColor="#8b5cf6"
            label="Language"
            subtitle="English"
            onPress={handleLanguageChange}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
            showBorder={false}
          />
        </Card>
      </View>

      {/* Data Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Data & Storage</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="delete-sweep"
            iconColor="#f97316"
            label="Clear Cache"
            subtitle="Free up storage space"
            onPress={handleClearCache}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
            showBorder={false}
          />
        </Card>
      </View>

      {/* About Section */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>About</Text>
        <Card variant="outlined" padding="none">
          <SettingRow
            icon="information"
            iconColor="#06b6d4"
            label="About Comply360"
            subtitle="Version 1.0.0"
            onPress={handleAbout}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
          />
          <SettingRow
            icon="shield-check"
            iconColor="#10b981"
            label="Privacy Policy"
            onPress={handlePrivacyPolicy}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
          />
          <SettingRow
            icon="file-document"
            iconColor="#f59e0b"
            label="Terms of Service"
            onPress={handleTermsOfService}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
          />
          <SettingRow
            icon="help-circle"
            iconColor="#3b82f6"
            label="Help & Support"
            onPress={handleHelp}
            rightComponent={<Icon name="chevron-right" size={24} color="#9ca3af" />}
            showBorder={false}
          />
        </Card>
      </View>

      {/* Footer */}
      <View style={styles.footer}>
        <Text style={styles.footerText}>Comply360 Mobile v1.0.0</Text>
        <Text style={styles.footerSubtext}>© 2025 Comply360. All rights reserved.</Text>
      </View>
    </ScrollView>
  );
};

interface SettingRowProps {
  icon: string;
  iconColor: string;
  label: string;
  subtitle?: string;
  onPress?: () => void;
  rightComponent?: React.ReactNode;
  showBorder?: boolean;
}

const SettingRow: React.FC<SettingRowProps> = ({
  icon,
  iconColor,
  label,
  subtitle,
  onPress,
  rightComponent,
  showBorder = true,
}) => {
  const content = (
    <View style={[styles.settingRow, !showBorder && styles.settingRowNoBorder]}>
      <View style={styles.settingLeft}>
        <View style={[styles.iconContainer, { backgroundColor: `${iconColor}15` }]}>
          <Icon name={icon} size={24} color={iconColor} />
        </View>
        <View style={styles.settingContent}>
          <Text style={styles.settingLabel}>{label}</Text>
          {subtitle && <Text style={styles.settingSubtitle}>{subtitle}</Text>}
        </View>
      </View>
      {rightComponent && <View style={styles.settingRight}>{rightComponent}</View>}
    </View>
  );

  if (onPress) {
    return (
      <TouchableOpacity onPress={onPress} activeOpacity={0.7}>
        {content}
      </TouchableOpacity>
    );
  }

  return content;
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  section: {
    paddingHorizontal: 20,
    paddingTop: 24,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 12,
  },
  settingRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 14,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
  },
  settingRowNoBorder: {
    borderBottomWidth: 0,
  },
  settingLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  iconContainer: {
    width: 44,
    height: 44,
    borderRadius: 22,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 12,
  },
  settingContent: {
    flex: 1,
  },
  settingLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#111827',
  },
  settingSubtitle: {
    fontSize: 13,
    color: '#6b7280',
    marginTop: 2,
  },
  settingRight: {
    marginLeft: 12,
  },
  footer: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  footerText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#6b7280',
    marginBottom: 4,
  },
  footerSubtext: {
    fontSize: 12,
    color: '#9ca3af',
  },
});

export default SettingsScreen;
