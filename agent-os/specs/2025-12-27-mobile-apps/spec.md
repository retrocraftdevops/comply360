# Mobile Applications (iOS/Android) - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Priority:** P0 - Critical  
**Estimated Duration:** 8-10 weeks  

---

## Executive Summary

Native mobile applications for iOS and Android to enable agents and clients to access Comply360 services on-the-go. Critical for market competitiveness as 40% of users prefer mobile platforms.

---

## Business Case

### Market Need
- 40% of users prefer mobile-first experience
- Agents work in the field and need mobile access
- Clients expect modern mobile capabilities
- Zero competitors have native mobile apps in SADC

### Business Impact
- **User Acquisition**: 40% larger addressable market
- **Agent Productivity**: 3x increase with field access
- **Client Satisfaction**: 50% improvement in NPS
- **Competitive Advantage**: First-mover in SADC mobile

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    MOBILE APPS                           │
│                                                          │
│  ┌──────────────┐              ┌──────────────┐        │
│  │  iOS App     │              │ Android App  │        │
│  │  (Swift UI)  │              │  (Kotlin)    │        │
│  └──────┬───────┘              └──────┬───────┘        │
│         │                              │                │
│         └──────────────┬───────────────┘                │
└────────────────────────┼─────────────────────────────────┘
                         │
                         │ HTTPS/REST API
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  API GATEWAY                             │
│  (Existing infrastructure)                               │
└─────────────────────────────────────────────────────────┘
```

**Technology Decision:**
- **Option A**: Native (Swift + Kotlin) - Best performance
- **Option B**: React Native - Single codebase
- **Option C**: Flutter - Fast development
- **Recommendation**: Flutter for speed-to-market and single codebase

---

## Core Features

### 1. Authentication & Onboarding

**Features:**
- Biometric login (Face ID, Touch ID, fingerprint)
- PIN code fallback
- OAuth social login (Google, Apple)
- Secure token storage (Keychain/Keystore)
- Offline authentication capability

**User Flow:**
```
Launch App
  ↓
Biometric Auth
  ↓
Dashboard
```

---

### 2. Dashboard (Agent View)

**Widgets:**
- Active registrations count
- Pending documents
- Commission this month
- Recent activity feed
- Quick actions (Start registration, Upload document)

**Features:**
- Real-time updates
- Pull-to-refresh
- Push notifications
- Offline viewing

---

### 3. Registration Wizard (Mobile-Optimized)

**Features:**
- Step-by-step wizard
- Progress indicator
- Auto-save to cloud
- Offline draft capability
- Camera integration for documents
- Form validation
- Smart suggestions

**Screens:**
1. Service selection
2. Company details
3. Directors/Shareholders
4. Documents upload
5. Review & submit
6. Payment

---

### 4. Document Management

**Camera Features:**
- Auto-capture with edge detection
- Multiple page scanning
- Image enhancement (contrast, brightness)
- PDF conversion
- OCR text extraction
- Document type classification

**Document Actions:**
- Capture from camera
- Select from gallery
- Sign documents
- Download for offline
- Share securely

---

### 5. Notifications

**Push Notifications:**
- Registration status updates
- Document approval/rejection
- Payment confirmations
- New commission earned
- System announcements

**In-App:**
- Notification center
- Read/unread status
- Action buttons
- Deep linking

---

### 6. Client Portal (Mobile)

**Features:**
- View registration status
- Upload documents
- Track progress
- Make payments
- Contact support
- View certificates

---

### 7. Offline Capability

**Offline Features:**
- View cached data
- Draft registrations
- Capture documents
- Queue for sync
- Offline indicator

**Sync:**
- Background sync
- Conflict resolution
- Progress indicator
- Retry failed syncs

---

## Technical Specifications

### Flutter Implementation

**Framework:** Flutter 3.16+  
**Language:** Dart 3.0+

**Key Packages:**
```yaml
dependencies:
  flutter: sdk: flutter
  
  # State Management
  riverpod: ^2.4.0
  
  # Networking
  dio: ^5.4.0
  retrofit: ^4.0.0
  
  # Storage
  hive: ^2.2.3
  hive_flutter: ^1.1.0
  secure_storage: ^9.0.0
  
  # Authentication
  local_auth: ^2.1.8
  flutter_secure_storage: ^9.0.0
  
  # Camera & Images
  camera: ^0.10.5
  image_picker: ^1.0.5
  image_cropper: ^5.0.1
  
  # Document Scanning
  edge_detection: ^1.1.0
  pdf: ^3.10.7
  
  # Biometrics
  local_auth: ^2.1.8
  
  # Push Notifications
  firebase_messaging: ^14.7.9
  flutter_local_notifications: ^16.3.0
  
  # QR/Barcode
  qr_code_scanner: ^1.0.1
  
  # Maps
  google_maps_flutter: ^2.5.3
  
  # Forms
  flutter_form_builder: ^9.1.1
  
  # UI Components
  flutter_svg: ^2.0.9
  cached_network_image: ^3.3.0
  shimmer: ^3.0.0
```

---

### Project Structure

```
lib/
├── main.dart
├── app.dart
├── core/
│   ├── config/
│   │   ├── app_config.dart
│   │   ├── api_config.dart
│   │   └── theme_config.dart
│   ├── network/
│   │   ├── api_client.dart
│   │   ├── interceptors.dart
│   │   └── error_handler.dart
│   ├── storage/
│   │   ├── secure_storage.dart
│   │   └── local_storage.dart
│   └── utils/
│       ├── validators.dart
│       ├── formatters.dart
│       └── helpers.dart
├── features/
│   ├── auth/
│   │   ├── data/
│   │   │   ├── models/
│   │   │   ├── repositories/
│   │   │   └── sources/
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   ├── repositories/
│   │   │   └── usecases/
│   │   └── presentation/
│   │       ├── pages/
│   │       ├── widgets/
│   │       └── providers/
│   ├── dashboard/
│   ├── registration/
│   ├── documents/
│   ├── payments/
│   └── profile/
├── shared/
│   ├── widgets/
│   ├── models/
│   └── providers/
└── router/
    └── app_router.dart
```

---

### State Management

**Pattern:** Riverpod (Provider-based)

```dart
// Example: Registration Provider
@riverpod
class RegistrationNotifier extends _$RegistrationNotifier {
  @override
  FutureOr<List<Registration>> build() async {
    return await ref.read(registrationRepositoryProvider)
        .getRegistrations();
  }

  Future<void> createRegistration(RegistrationData data) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await ref.read(registrationRepositoryProvider)
          .createRegistration(data);
      return await ref.read(registrationRepositoryProvider)
          .getRegistrations();
    });
  }
}
```

---

### API Integration

**REST Client:**
```dart
@RestApi(baseUrl: "https://api.comply360.com/v1")
abstract class ApiClient {
  factory ApiClient(Dio dio, {String baseUrl}) = _ApiClient;

  @GET("/registrations")
  Future<List<Registration>> getRegistrations(
    @Header("Authorization") String token,
    @Query("tenantId") String tenantId,
  );

  @POST("/registrations")
  Future<Registration> createRegistration(
    @Header("Authorization") String token,
    @Body() RegistrationData data,
  );

  @POST("/documents/upload")
  @MultiPart()
  Future<Document> uploadDocument(
    @Header("Authorization") String token,
    @Part(name: "file") File file,
    @Part(name: "type") String type,
  );
}
```

---

### Offline Storage

**Hive Database:**
```dart
// Models
@HiveType(typeId: 0)
class Registration extends HiveObject {
  @HiveField(0)
  final String id;
  
  @HiveField(1)
  final String companyName;
  
  @HiveField(2)
  final String status;
  
  @HiveField(3)
  final DateTime createdAt;
  
  @HiveField(4)
  final bool synced;
}

// Repository
class RegistrationRepository {
  final Box<Registration> _box;
  
  Future<void> saveOffline(Registration registration) async {
    registration.synced = false;
    await _box.put(registration.id, registration);
  }
  
  Future<void> syncToCloud() async {
    final unsynced = _box.values.where((r) => !r.synced);
    for (final registration in unsynced) {
      await apiClient.createRegistration(registration);
      registration.synced = true;
      await registration.save();
    }
  }
}
```

---

### Security

**Secure Storage:**
```dart
class SecureStorageService {
  final FlutterSecureStorage _storage = FlutterSecureStorage();
  
  Future<void> saveToken(String token) async {
    await _storage.write(
      key: 'auth_token',
      value: token,
      iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock),
      aOptions: AndroidOptions(encryptedSharedPreferences: true),
    );
  }
  
  Future<String?> getToken() async {
    return await _storage.read(key: 'auth_token');
  }
}
```

**Biometric Authentication:**
```dart
class BiometricService {
  final LocalAuthentication _auth = LocalAuthentication();
  
  Future<bool> authenticate() async {
    try {
      final canCheck = await _auth.canCheckBiometrics;
      if (!canCheck) return false;
      
      return await _auth.authenticate(
        localizedReason: 'Please authenticate to access Comply360',
        options: const AuthenticationOptions(
          stickyAuth: true,
          biometricOnly: true,
        ),
      );
    } catch (e) {
      return false;
    }
  }
}
```

---

### Camera & Document Scanning

**Document Scanner:**
```dart
class DocumentScannerService {
  Future<List<File>> scanDocument() async {
    try {
      // Use edge detection for automatic capture
      final result = await EdgeDetection.detectEdge;
      
      if (result != null) {
        // Enhance image
        final enhanced = await _enhanceImage(result);
        
        // Convert to PDF
        final pdf = await _convertToPDF(enhanced);
        
        return [pdf];
      }
      return [];
    } catch (e) {
      throw DocumentScanException(e.toString());
    }
  }
  
  Future<File> _enhanceImage(String imagePath) async {
    final image = img.decodeImage(File(imagePath).readAsBytesSync())!;
    
    // Apply enhancements
    final enhanced = img.contrast(image, 120);
    final sharpened = img.convolution(enhanced, [0, -1, 0, -1, 5, -1, 0, -1, 0]);
    
    return File(imagePath)..writeAsBytesSync(img.encodePng(sharpened));
  }
}
```

---

### Push Notifications

**Firebase Cloud Messaging:**
```dart
class NotificationService {
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;
  
  Future<void> initialize() async {
    // Request permissions
    await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );
    
    // Get token
    final token = await _messaging.getToken();
    await _saveTokenToServer(token);
    
    // Handle foreground messages
    FirebaseMessaging.onMessage.listen(_handleMessage);
    
    // Handle background messages
    FirebaseMessaging.onBackgroundMessage(_handleBackgroundMessage);
  }
  
  void _handleMessage(RemoteMessage message) {
    final notification = message.notification;
    if (notification != null) {
      _showLocalNotification(
        notification.title ?? '',
        notification.body ?? '',
        message.data,
      );
    }
  }
}
```

---

## UI/UX Design

### Design System

**Colors:**
```dart
class AppColors {
  static const primary = Color(0xFF0066CC);
  static const secondary = Color(0xFF00C853);
  static const background = Color(0xFFF5F5F5);
  static const surface = Color(0xFFFFFFFF);
  static const error = Color(0xFFD32F2F);
  static const success = Color(0xFF388E3C);
  static const warning = Color(0xFFF57C00);
}
```

**Typography:**
```dart
class AppTextStyles {
  static const heading1 = TextStyle(
    fontSize: 32,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.5,
  );
  
  static const body1 = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.normal,
    height: 1.5,
  );
}
```

---

### Key Screens

#### 1. Login Screen
- Logo
- Biometric button
- Email/password fields
- "Remember me" checkbox
- Forgot password link
- Sign up link

#### 2. Dashboard
- Header with greeting
- Statistics cards
- Recent activity list
- Quick action buttons
- Bottom navigation

#### 3. Registration Form
- Progress stepper
- Form sections
- Smart validation
- Save draft button
- Next/Previous buttons

#### 4. Document Camera
- Camera viewfinder
- Edge detection overlay
- Capture button
- Flash toggle
- Gallery button

---

## Platform-Specific Features

### iOS

**Features:**
- Face ID integration
- Haptic feedback
- iOS design guidelines
- Apple Wallet integration (certificates)
- Siri shortcuts
- Widgets

**Requirements:**
- Xcode 15+
- iOS 15.0+
- Swift 5.9+

---

### Android

**Features:**
- Fingerprint integration
- Material Design 3
- Android-specific permissions
- Google Wallet integration
- Voice shortcuts
- Widgets

**Requirements:**
- Android Studio Hedgehog+
- Min SDK: 24 (Android 7.0)
- Target SDK: 34 (Android 14)
- Kotlin 1.9+

---

## Performance Requirements

- **App Size**: < 50MB
- **Launch Time**: < 2 seconds (cold start)
- **Screen Load**: < 500ms
- **API Calls**: < 200ms (cached), < 2s (network)
- **Battery Usage**: < 5% per hour active use
- **Offline**: Full functionality for viewing and drafting

---

## Testing Strategy

### Unit Tests
- Business logic
- Data models
- Validators
- Formatters

### Widget Tests
- UI components
- Forms
- Navigation
- State management

### Integration Tests
- API integration
- Database operations
- Authentication flow
- Complete user journeys

### Platform Tests
- iOS device testing (iPhone 12+, iPad)
- Android device testing (Samsung, Pixel, Xiaomi)
- Different screen sizes
- Different OS versions

---

## Deployment

### App Store (iOS)

**Requirements:**
- Apple Developer Account ($99/year)
- App Store guidelines compliance
- Privacy policy
- App Store screenshots
- App Store description

**Process:**
1. Create app in App Store Connect
2. Configure certificates and profiles
3. Build release version
4. Upload to App Store Connect
5. Submit for review
6. Release to users

---

### Google Play (Android)

**Requirements:**
- Google Play Developer Account ($25 one-time)
- Play Store guidelines compliance
- Privacy policy
- Play Store screenshots
- Play Store description

**Process:**
1. Create app in Play Console
2. Configure signing keys
3. Build release bundle
4. Upload to Play Console
5. Submit for review
6. Release to production

---

## Analytics & Monitoring

**Tools:**
- Firebase Analytics
- Crashlytics for crash reporting
- Performance Monitoring
- Remote Config for feature flags

**Metrics:**
- Daily active users (DAU)
- Monthly active users (MAU)
- Session duration
- Screen views
- Conversion rates
- Crash-free rate (target: 99.9%)
- App size
- Network usage

---

## Security Measures

1. **Data Encryption**: All data encrypted at rest and in transit
2. **Certificate Pinning**: Prevent man-in-the-middle attacks
3. **Jailbreak/Root Detection**: Warn users of compromised devices
4. **Secure Storage**: Use platform secure storage (Keychain/Keystore)
5. **Code Obfuscation**: Protect source code
6. **API Security**: Token-based authentication, refresh tokens

---

## Implementation Roadmap

### Phase 1: Foundation (2 weeks)
- Project setup
- Architecture implementation
- API integration
- Authentication flow
- Basic UI components

### Phase 2: Core Features (3 weeks)
- Dashboard
- Registration wizard
- Document management
- Push notifications
- Offline capability

### Phase 3: Advanced Features (2 weeks)
- Camera scanning
- Biometric auth
- Payment integration
- Advanced UI polish
- Performance optimization

### Phase 4: Testing & Launch (3 weeks)
- Comprehensive testing
- Bug fixes
- Beta testing
- App Store submission
- Production launch

---

## Success Metrics

- **Downloads**: 10,000+ in first 3 months
- **Active Users**: 60% of total users
- **Rating**: 4.5+ stars
- **Crash Rate**: < 0.1%
- **User Satisfaction**: 90%+ positive reviews

---

**Next Steps:** See `tasks.md` for detailed implementation tasks

