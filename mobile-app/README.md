# Comply360 Mobile App

**SADC Corporate Gateway Platform - Mobile Application**

A production-ready React Native mobile application for iOS and Android, providing corporate registration, document management, and commission tracking services across Southern Africa.

---

## 📱 Features

### Authentication
- ✅ Email/password login
- ✅ Biometric authentication (Touch ID / Face ID)
- ✅ Secure credential storage (Keychain)
- ✅ Password reset flow
- ✅ Remember me functionality
- ✅ Automatic token refresh

### Dashboard
- ✅ Real-time statistics from API
- ✅ Registration metrics
- ✅ Commission tracking
- ✅ Document status overview
- ✅ Quick action buttons
- ✅ Pull-to-refresh

### Registration Management (Sprint 2)
- ✅ Multi-step registration form (5 steps)
- ✅ Company information capture
- ✅ Contact details validation
- ✅ Business details collection
- ✅ Document upload placeholder
- ✅ Review before submit
- ✅ Draft saving capability
- ✅ Registration list with search
- ✅ Status filtering (Draft, Pending, In Progress, Completed, Rejected)
- ✅ Pull-to-refresh

### Commission Management (Sprint 2)
- ✅ Commission list with summary cards
- ✅ Real-time totals (Pending, Approved, Paid)
- ✅ Search and filter commissions
- ✅ Payout request flow
- ✅ Commission details view
- ✅ Status tracking

### Document Management (Sprint 2)
- ✅ Document list with filtering
- ✅ Status badges (Pending, Verified, Rejected)
- ✅ File type filtering (PDF, Images, Documents, Spreadsheets)
- ✅ Document download
- ✅ Upload options (Camera, Gallery, Files)
- ✅ Search documents
- ✅ File size display

### Profile Management (Sprint 3)
- ✅ User profile viewing with avatar
- ✅ Profile editing with validation
- ✅ Avatar upload (camera/gallery placeholders)
- ✅ Account statistics display
- ✅ Logout functionality

### Settings & Preferences (Sprint 3)
- ✅ Dark mode theme support
- ✅ Theme persistence (light/dark/system)
- ✅ Notification preferences (push, email, SMS)
- ✅ Biometric login toggle
- ✅ Language selection (placeholder)
- ✅ Clear cache option
- ✅ About and legal information

### Notifications (Sprint 3)
- ✅ Notification list with filtering
- ✅ Unread count badges
- ✅ Mark as read functionality
- ✅ Mark all as read
- ✅ Clear all notifications
- ✅ Filter by: ALL, UNREAD, READ
- ✅ Push notification infrastructure (ready)

### Offline Sync (Sprint 3)
- ✅ Offline queue system
- ✅ Network status detection
- ✅ Automatic retry logic
- ✅ Queue persistence (AsyncStorage)
- ✅ Auto-sync when online
- ✅ Sync status indicators

### Performance & UX (Sprint 3)
- ✅ Loading skeleton components
- ✅ Animated shimmer effects
- ✅ Skeleton variants for all card types
- ✅ useMemo optimizations
- ✅ Smooth theme transitions

### Core Features
- ✅ User profile management
- ✅ Offline data persistence
- ✅ Push notifications (ready)

---

## 🏗️ Architecture

### Tech Stack
- **Framework**: React Native 0.73.2
- **Language**: TypeScript 5.3
- **State Management**: Redux Toolkit 2.0 + RTK Query
- **Navigation**: React Navigation 6
- **UI Library**: React Native Paper 5.12
- **Icons**: Material Community Icons
- **Storage**: AsyncStorage + Redux Persist
- **Security**: React Native Keychain + Biometrics

### Project Structure
```
mobile-app/
├── src/
│   ├── navigation/              # Navigation controllers
│   │   ├── AppNavigator.tsx     # Root navigation
│   │   ├── AuthNavigator.tsx    # Auth flow
│   │   └── TabNavigator.tsx     # Main app tabs
│   │
│   ├── screens/                 # Screen components
│   │   ├── Auth/                # Authentication screens
│   │   │   ├── LoginScreen.tsx
│   │   │   ├── BiometricSetupScreen.tsx
│   │   │   └── ForgotPasswordScreen.tsx
│   │   ├── Dashboard/           # Dashboard
│   │   ├── Registrations/       # Registration screens
│   │   ├── Documents/           # Document screens
│   │   ├── Commissions/         # Commission screens
│   │   ├── Profile/             # Profile screens
│   │   └── SplashScreen.tsx     # Splash screen
│   │
│   ├── lib/                     # Reusable code
│   │   ├── components/          # UI components (20 total + 9 skeleton variants)
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── LoadingSpinner.tsx
│   │   │   ├── EmptyState.tsx
│   │   │   ├── ErrorBoundary.tsx
│   │   │   ├── Toast.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── BottomSheet.tsx
│   │   │   ├── FormWizard.tsx       # Sprint 2
│   │   │   ├── StepIndicator.tsx    # Sprint 2
│   │   │   ├── FormInput.tsx        # Sprint 2
│   │   │   ├── SearchBar.tsx        # Sprint 2
│   │   │   ├── RegistrationCard.tsx # Sprint 2
│   │   │   ├── CommissionCard.tsx   # Sprint 2
│   │   │   ├── DocumentCard.tsx     # Sprint 2
│   │   │   ├── Avatar.tsx           # Sprint 3
│   │   │   ├── NotificationCard.tsx # Sprint 3
│   │   │   ├── SyncIndicator.tsx    # Sprint 3
│   │   │   └── LoadingSkeleton.tsx  # Sprint 3 (+ 9 variants)
│   │   │
│   │   └── utils/               # Utilities
│   │       ├── theme.ts         # Light theme
│   │       ├── theme-dark.ts    # Dark theme (Sprint 3)
│   │       ├── validation.ts    # Validators
│   │       ├── formatting.ts    # Formatters
│   │       └── constants.ts     # Constants
│   │
│   ├── store/                   # Redux store
│   │   ├── store.ts             # Store configuration
│   │   ├── slices/              # State slices
│   │   │   ├── authSlice.ts
│   │   │   ├── registrationSlice.ts
│   │   │   ├── documentSlice.ts
│   │   │   ├── commissionSlice.ts
│   │   │   ├── notificationSlice.ts
│   │   │   └── uiSlice.ts
│   │   │
│   │   └── api/                 # RTK Query APIs
│   │       ├── authApi.ts
│   │       ├── registrationApi.ts
│   │       ├── documentApi.ts
│   │       └── commissionApi.ts
│   │
│   ├── contexts/                # React contexts
│   │   └── ThemeContext.tsx     # Theme management (Sprint 3)
│   │
│   └── services/                # Business logic
│       ├── auth.ts              # Auth service
│       ├── biometrics.ts        # Biometric service
│       ├── offlineQueue.ts      # Offline queue (Sprint 3)
│       └── networkStatus.ts     # Network detection (Sprint 3)
│
├── android/                     # Android native code
├── ios/                         # iOS native code
├── App.tsx                      # App entry point
├── package.json                 # Dependencies
└── tsconfig.json                # TypeScript config
```

---

## 🚀 Getting Started

### Prerequisites
- Node.js >= 18
- npm or yarn
- React Native CLI
- Xcode (for iOS)
- Android Studio (for Android)

### Installation

```bash
# Clone the repository
cd mobile-app

# Install dependencies
npm install

# Install iOS pods (macOS only)
cd ios && pod install && cd ..
```

### Configuration

Create a `.env` file in the root:

```env
API_URL=http://localhost:8080/api/v1
```

For production, update the API URL to your production backend.

### Running the App

**iOS:**
```bash
npm run ios
# or
npx react-native run-ios
```

**Android:**
```bash
npm run android
# or
npx react-native run-android
```

**Development:**
```bash
# Start Metro bundler
npm start

# Clear cache and start
npm start -- --reset-cache
```

---

## 📚 Usage Guide

### Authentication

```typescript
import { useAppDispatch } from '@/store/store';
import { loginSuccess } from '@/store/slices/authSlice';
import { AuthService } from '@/services/auth';

// Login
const handleLogin = async (email: string, password: string) => {
  try {
    const response = await AuthService.login(email, password);
    dispatch(loginSuccess({
      user: response.user,
      token: response.token,
      refreshToken: response.refreshToken,
    }));
  } catch (error) {
    console.error('Login failed:', error);
  }
};
```

### API Calls with RTK Query

```typescript
import { useGetRegistrationsQuery } from '@/store/api/registrationApi';

const MyComponent = () => {
  const { data, isLoading, error, refetch } = useGetRegistrationsQuery({
    page: 1,
    limit: 20,
  });

  return (
    <ScrollView refreshControl={
      <RefreshControl refreshing={isLoading} onRefresh={refetch} />
    }>
      {data?.registrations.map(reg => (
        <Text key={reg.id}>{reg.company_name}</Text>
      ))}
    </ScrollView>
  );
};
```

### Using Components

```typescript
import { Button, Card, Modal, Toast } from '@/lib/components';

// Button
<Button
  title="Submit"
  onPress={handleSubmit}
  variant="primary"
  loading={isLoading}
  icon="check"
/>

// Card
<Card variant="elevated" padding="large">
  <Text>Card content</Text>
</Card>

// Modal
<Modal
  visible={showModal}
  onClose={() => setShowModal(false)}
  title="Confirmation"
  primaryButton={{
    title: 'Confirm',
    onPress: handleConfirm,
  }}
>
  <Text>Are you sure?</Text>
</Modal>
```

### Validation

```typescript
import { validateEmail, validatePassword } from '@/lib/utils';

const email = 'user@example.com';
const isValid = validateEmail(email); // true

const password = 'MyP@ssw0rd';
const result = validatePassword(password);
// { isValid: true, errors: [] }
```

### Formatting

```typescript
import { formatCurrency, formatDate, formatPhoneNumber } from '@/lib/utils';

formatCurrency(1234.56); // "R 1,234.56"
formatDate(new Date(), 'medium'); // "Dec 28, 2025"
formatPhoneNumber('0123456789'); // "012 345 6789"
```

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Watch mode
npm run test:watch
```

### Manual Testing Checklist

**Authentication:**
- [ ] Login with valid credentials
- [ ] Login with invalid credentials
- [ ] Biometric setup flow
- [ ] Biometric login
- [ ] Forgot password
- [ ] Logout

**Dashboard:**
- [ ] Load stats from API
- [ ] Pull to refresh
- [ ] Quick actions navigate correctly
- [ ] Notification badge shows count

**Offline:**
- [ ] App works offline
- [ ] Data persists after restart
- [ ] Syncs when online

---

## 🏭 Building for Production

### iOS

```bash
# Create release build
cd ios
xcodebuild -workspace Comply360.xcworkspace \
  -scheme Comply360 \
  -configuration Release \
  -archivePath build/Comply360.xcarchive \
  archive

# Export IPA
xcodebuild -exportArchive \
  -archivePath build/Comply360.xcarchive \
  -exportPath build \
  -exportOptionsPlist ExportOptions.plist
```

### Android

```bash
# Create release APK
cd android
./gradlew assembleRelease

# Create release AAB (for Play Store)
./gradlew bundleRelease

# Output: android/app/build/outputs/
```

### Code Signing

**iOS:**
1. Add provisioning profile to Xcode
2. Update signing team in Xcode
3. Configure push notification certificates

**Android:**
1. Generate keystore: `keytool -genkey -v -keystore comply360.keystore`
2. Update `android/gradle.properties` with keystore details
3. Add keystore to `android/app/`

---

## 📦 Dependencies

### Core
- `react`: 18.2.0
- `react-native`: 0.73.2
- `typescript`: 5.3.3

### Navigation
- `@react-navigation/native`: ^6.1.9
- `@react-navigation/bottom-tabs`: ^6.5.11
- `@react-navigation/stack`: ^6.3.20

### State Management
- `@reduxjs/toolkit`: ^2.0.1
- `react-redux`: ^9.0.4
- `redux-persist`: ^6.0.0

### UI
- `react-native-paper`: ^5.12.1
- `react-native-vector-icons`: ^10.0.3
- `react-native-reanimated`: ^3.6.1

### Native Features
- `react-native-biometrics`: ^3.0.1
- `react-native-keychain`: ^8.1.2
- `react-native-camera`: ^4.2.1
- `react-native-document-scanner-plugin`: ^0.4.2

### Utilities
- `axios`: ^1.6.5
- `date-fns`: ^3.0.6
- `react-native-config`: ^1.5.1

---

## 🔒 Security

### Best Practices Implemented

✅ Secure credential storage (Keychain)
✅ Biometric authentication
✅ Token-based auth with auto-refresh
✅ HTTPS-only API calls
✅ Input validation and sanitization
✅ Error boundary for crash protection
✅ No sensitive data in logs (production)

### Security Considerations

- **Never commit** `.env` files
- **Always use** HTTPS in production
- **Enable** certificate pinning for production
- **Implement** root detection (optional)
- **Add** jailbreak detection (optional)
- **Use** ProGuard/R8 for Android obfuscation

---

## 🐛 Troubleshooting

### Common Issues

**Metro bundler won't start:**
```bash
npm start -- --reset-cache
```

**Pod install fails (iOS):**
```bash
cd ios
pod deintegrate
pod install
```

**Android build fails:**
```bash
cd android
./gradlew clean
./gradlew assembleDebug
```

**Biometrics not working:**
- Ensure device has biometrics enabled
- Check permissions in Info.plist (iOS) / AndroidManifest.xml
- Test on physical device (not simulator)

---

## 📄 License

Copyright © 2025 Comply360. All rights reserved.

---

## 📞 Support

**Email:** support@comply360.com
**Phone:** +27 11 123 4567
**Hours:** Mon-Fri: 8:00 AM - 6:00 PM SAST

---

## 🚧 Roadmap

### Sprint 2 (Weeks 3-4)
- [ ] Enhanced registration forms
- [ ] Document scanner integration
- [ ] Camera integration
- [ ] Advanced filtering and search
- [ ] Commission payout requests

### Sprint 3 (Weeks 5-6)
- [ ] Real-time notifications
- [ ] Offline document queue
- [ ] Advanced analytics
- [ ] Multi-language support
- [ ] Dark mode

---

**Version:** 1.0.0
**Last Updated:** December 28, 2025
**Maintained by:** Comply360 Development Team
