/**
 * Splash Screen
 * Displayed while app is initializing and loading persisted state
 */

import React from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';

const SplashScreen = () => {
  return (
    <View style={styles.container}>
      {/* Logo */}
      <View style={styles.logoContainer}>
        <View style={styles.logoCircle}>
          <Icon name="domain" size={80} color="#7c3aed" />
        </View>
        <Text style={styles.title}>Comply360</Text>
        <Text style={styles.subtitle}>SADC Corporate Gateway</Text>
      </View>

      {/* Loading Indicator */}
      <ActivityIndicator size="large" color="#7c3aed" style={styles.loader} />

      {/* Footer */}
      <Text style={styles.footer}>Initializing...</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#FFFFFF',
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: 24,
  },
  logoContainer: {
    alignItems: 'center',
    marginBottom: 60,
  },
  logoCircle: {
    width: 140,
    height: 140,
    borderRadius: 70,
    backgroundColor: '#f3e8ff',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 20,
  },
  title: {
    fontSize: 36,
    fontWeight: '700',
    color: '#111827',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 16,
    color: '#6b7280',
  },
  loader: {
    marginTop: 40,
  },
  footer: {
    position: 'absolute',
    bottom: 40,
    fontSize: 14,
    color: '#6b7280',
  },
});

export default SplashScreen;
