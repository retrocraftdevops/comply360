/**
 * ProfileScreen
 * User profile display with account information and actions
 */

import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { Avatar, Card, Button } from '@/lib/components';
import { useAppDispatch, useAppSelector } from '@/store/store';
import { logout } from '@/store/slices/authSlice';
import { showToast } from '@/store/slices/uiSlice';
import { colors, spacing, fonts } from '@/lib/utils/theme';

const ProfileScreen: React.FC = () => {
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.auth.user);

  const handleEditProfile = () => {
    dispatch(showToast({ message: 'Edit profile coming soon!', type: 'info' }));
  };

  const handleSettings = () => {
    dispatch(showToast({ message: 'Settings coming soon!', type: 'info' }));
  };

  const handleLogout = () => {
    Alert.alert('Logout', 'Are you sure you want to logout?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Logout',
        style: 'destructive',
        onPress: () => {
          dispatch(logout());
          dispatch(showToast({ message: 'Logged out successfully', type: 'success' }));
        },
      },
    ]);
  };

  const stats = {
    registrations: 12,
    commissions: 'R 25,430',
    documents: 34,
  };

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <View style={styles.header}>
        <Avatar name={user?.name || 'User'} size="xlarge" showBadge editable />
        <Text style={styles.name}>{user?.name || 'Unknown User'}</Text>
        <Text style={styles.email}>{user?.email || 'No email'}</Text>

        <View style={styles.actionButtons}>
          <Button
            title="Edit Profile"
            onPress={handleEditProfile}
            variant="primary"
            icon="account-edit"
            size="medium"
          />
          <View style={{ width: 12 }} />
          <Button
            title="Settings"
            onPress={handleSettings}
            variant="outline"
            icon="cog"
            size="medium"
          />
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Account Statistics</Text>
        <View style={styles.statsContainer}>
          <StatCard icon="file-document-multiple" label="Registrations" value="12" color="#7c3aed" />
          <StatCard icon="currency-usd" label="Commissions" value="R 25,430" color="#10b981" />
          <StatCard icon="file-pdf-box" label="Documents" value="34" color="#3b82f6" />
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Account Information</Text>
        <Card variant="outlined" padding="none">
          <InfoRow icon="account" label="Full Name" value={user?.name || 'Not set'} />
          <InfoRow icon="email" label="Email Address" value={user?.email || 'Not set'} />
          <InfoRow icon="phone" label="Phone Number" value={user?.phone || 'Not set'} />
          <InfoRow icon="domain" label="Company" value={user?.company || 'Not set'} />
          <InfoRow icon="shield-account" label="Role" value="Agent" showBorder={false} />
        </Card>
      </View>

      <View style={styles.section}>
        <Button title="Logout" onPress={handleLogout} variant="danger" icon="logout" fullWidth />
      </View>

      <View style={styles.footer}>
        <Text style={styles.version}>Comply360 Mobile v1.0.0</Text>
        <Text style={styles.copyright}>© 2025 Comply360. All rights reserved.</Text>
      </View>
    </ScrollView>
  );
};

interface StatCardProps {
  icon: string;
  label: string;
  value: string;
  color: string;
}

const StatCard: React.FC<StatCardProps> = ({ icon, label, value, color }) => (
  <View style={styles.statCard}>
    <View style={[styles.statIcon, { backgroundColor: `${color}15` }]}>
      <Icon name={icon} size={24} color={color} />
    </View>
    <Text style={styles.statValue}>{value}</Text>
    <Text style={styles.statLabel}>{label}</Text>
  </View>
);

interface InfoRowProps {
  icon: string;
  label: string;
  value: string;
  showBorder?: boolean;
}

const InfoRow: React.FC<InfoRowProps> = ({ icon, label, value, showBorder = true }) => (
  <View style={[styles.infoRow, !showBorder && styles.infoRowNoBorder]}>
    <View style={styles.infoLeft}>
      <Icon name={icon} size={20} color="#6b7280" />
      <Text style={styles.infoLabel}>{label}</Text>
    </View>
    <Text style={styles.infoValue}>{value}</Text>
  </View>
);

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f5f5f5' },
  header: { backgroundColor: '#FFFFFF', alignItems: 'center', paddingTop: 32, paddingBottom: 20, paddingHorizontal: 20 },
  name: { fontSize: 24, fontWeight: '700', color: '#111827', marginTop: 12 },
  email: { fontSize: 16, color: '#6b7280', marginTop: 4 },
  actionButtons: { flexDirection: 'row', marginTop: 20 },
  section: { paddingHorizontal: 20, paddingTop: 32 },
  sectionTitle: { fontSize: 18, fontWeight: '700', color: '#111827', marginBottom: 12 },
  statsContainer: { flexDirection: 'row', gap: 12 },
  statCard: { flex: 1, backgroundColor: '#FFFFFF', borderRadius: 12, padding: 12, alignItems: 'center' },
  statIcon: { width: 48, height: 48, borderRadius: 24, alignItems: 'center', justifyContent: 'center', marginBottom: 8 },
  statValue: { fontSize: 20, fontWeight: '700', color: '#111827', marginBottom: 4 },
  statLabel: { fontSize: 12, color: '#6b7280', textAlign: 'center' },
  infoRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: 12, paddingHorizontal: 12, borderBottomWidth: 1, borderBottomColor: '#e5e7eb' },
  infoRowNoBorder: { borderBottomWidth: 0 },
  infoLeft: { flexDirection: 'row', alignItems: 'center', flex: 1 },
  infoLabel: { fontSize: 16, color: '#111827', marginLeft: 12 },
  infoValue: { fontSize: 16, color: '#6b7280', fontWeight: '500' },
  footer: { alignItems: 'center', paddingVertical: 32 },
  version: { fontSize: 14, color: '#6b7280', marginBottom: 4 },
  copyright: { fontSize: 12, color: '#9ca3af' },
});

export default ProfileScreen;
