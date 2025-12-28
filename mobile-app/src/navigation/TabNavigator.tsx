/**
 * Tab Navigator
 * Bottom tab navigation for main app screens
 */

import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';

import DashboardScreen from '../screens/Dashboard/DashboardScreen';
import RegistrationListScreen from '../screens/Registrations/RegistrationListScreen';
import DocumentListScreen from '../screens/Documents/DocumentListScreen';
import CommissionListScreen from '../screens/Commissions/CommissionListScreen';
import ProfileScreen from '../screens/Profile/ProfileScreen';

export type TabParamList = {
  Dashboard: undefined;
  Registrations: undefined;
  Documents: undefined;
  Commissions: undefined;
  Profile: undefined;
};

const Tab = createBottomTabNavigator<TabParamList>();

const TabNavigator = () => {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: string;

          switch (route.name) {
            case 'Dashboard':
              iconName = focused ? 'view-dashboard' : 'view-dashboard-outline';
              break;
            case 'Registrations':
              iconName = focused ? 'file-document' : 'file-document-outline';
              break;
            case 'Documents':
              iconName = focused ? 'folder' : 'folder-outline';
              break;
            case 'Commissions':
              iconName = focused ? 'cash-multiple' : 'cash';
              break;
            case 'Profile':
              iconName = focused ? 'account' : 'account-outline';
              break;
            default:
              iconName = 'circle';
          }

          return <Icon name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: '#7c3aed',
        tabBarInactiveTintColor: '#6b7280',
        tabBarStyle: {
          backgroundColor: '#ffffff',
          borderTopWidth: 1,
          borderTopColor: '#e5e7eb',
          height: 60,
          paddingBottom: 8,
          paddingTop: 8,
        },
        tabBarLabelStyle: {
          fontSize: 12,
          fontWeight: '600',
        },
        headerShown: false,
      })}
    >
      <Tab.Screen
        name="Dashboard"
        component={DashboardScreen}
        options={{ title: 'Dashboard' }}
      />
      <Tab.Screen
        name="Registrations"
        component={RegistrationListScreen}
        options={{ title: 'Registrations' }}
      />
      <Tab.Screen
        name="Documents"
        component={DocumentListScreen}
        options={{ title: 'Documents' }}
      />
      <Tab.Screen
        name="Commissions"
        component={CommissionListScreen}
        options={{ title: 'Commissions' }}
      />
      <Tab.Screen
        name="Profile"
        component={ProfileScreen}
        options={{ title: 'Profile' }}
      />
    </Tab.Navigator>
  );
};

export default TabNavigator;
