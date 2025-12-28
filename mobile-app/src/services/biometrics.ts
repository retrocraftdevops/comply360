/**
 * Biometric Authentication Service
 * Handles fingerprint/Face ID authentication and secure credential storage
 */

import ReactNativeBiometrics, { BiometryTypes } from 'react-native-biometrics';
import * as Keychain from 'react-native-keychain';
import { Platform } from 'react-native';

// Service identifier for Keychain
const KEYCHAIN_SERVICE = 'com.comply360.auth';
const BIOMETRIC_PROMPT_TITLE = 'Authenticate';
const BIOMETRIC_PROMPT_DESCRIPTION = 'Use biometrics to sign in to Comply360';

// Stored credentials interface
export interface StoredCredentials {
  email: string;
  password: string;
}

/**
 * Biometric Service Class
 */
class BiometricAuthService {
  private rnBiometrics: ReactNativeBiometrics;

  constructor() {
    this.rnBiometrics = new ReactNativeBiometrics({
      allowDeviceCredentials: false, // Only allow biometrics, not PIN/pattern
    });
  }

  /**
   * Check if biometric authentication is available on the device
   * @returns True if biometrics are available
   */
  async isBiometricAvailable(): Promise<boolean> {
    try {
      const { available, biometryType } = await this.rnBiometrics.isSensorAvailable();

      if (available && biometryType) {
        console.log('[BiometricService] Biometric type available:', biometryType);
        return true;
      }

      console.log('[BiometricService] No biometric sensor available');
      return false;
    } catch (error: any) {
      console.error('[BiometricService] Error checking biometric availability:', error);
      return false;
    }
  }

  /**
   * Get the type of biometric authentication available
   * @returns Biometric type string (TouchID, FaceID, Biometrics, or null)
   */
  async getBiometricType(): Promise<string | null> {
    try {
      const { available, biometryType } = await this.rnBiometrics.isSensorAvailable();

      if (!available || !biometryType) {
        return null;
      }

      switch (biometryType) {
        case BiometryTypes.TouchID:
          return 'Touch ID';
        case BiometryTypes.FaceID:
          return 'Face ID';
        case BiometryTypes.Biometrics:
          return Platform.OS === 'android' ? 'Fingerprint' : 'Biometrics';
        default:
          return 'Biometrics';
      }
    } catch (error: any) {
      console.error('[BiometricService] Error getting biometric type:', error);
      return null;
    }
  }

  /**
   * Prompt user for biometric authentication
   * @param promptTitle Optional custom title
   * @param promptDescription Optional custom description
   * @returns True if authentication was successful
   */
  async authenticate(
    promptTitle: string = BIOMETRIC_PROMPT_TITLE,
    promptDescription: string = BIOMETRIC_PROMPT_DESCRIPTION
  ): Promise<boolean> {
    try {
      const { success } = await this.rnBiometrics.simplePrompt({
        promptMessage: promptTitle,
        ...(Platform.OS === 'android' && {
          cancelButtonText: 'Cancel',
        }),
      });

      if (success) {
        console.log('[BiometricService] Biometric authentication successful');
        return true;
      } else {
        console.log('[BiometricService] Biometric authentication failed or cancelled');
        return false;
      }
    } catch (error: any) {
      console.error('[BiometricService] Biometric authentication error:', error);
      return false;
    }
  }

  /**
   * Store user credentials securely in device keychain
   * @param email User email
   * @param password User password
   * @returns True if credentials were stored successfully
   */
  async storeCredentials(email: string, password: string): Promise<boolean> {
    try {
      await Keychain.setGenericPassword(email, password, {
        service: KEYCHAIN_SERVICE,
        accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED_THIS_DEVICE_ONLY,
        accessControl: Keychain.ACCESS_CONTROL.BIOMETRY_ANY,
        securityLevel: Keychain.SECURITY_LEVEL.SECURE_HARDWARE,
      });

      console.log('[BiometricService] Credentials stored successfully');
      return true;
    } catch (error: any) {
      console.error('[BiometricService] Error storing credentials:', error);
      return false;
    }
  }

  /**
   * Retrieve stored credentials from device keychain
   * Requires biometric authentication
   * @returns Stored credentials or null if not found/authentication failed
   */
  async getStoredCredentials(): Promise<StoredCredentials | null> {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: KEYCHAIN_SERVICE,
        authenticationPrompt: {
          title: 'Authenticate',
          subtitle: 'Access your stored credentials',
          description: 'Use biometrics to retrieve your login credentials',
          cancel: 'Cancel',
        },
      });

      if (credentials && credentials.username && credentials.password) {
        console.log('[BiometricService] Credentials retrieved successfully');
        return {
          email: credentials.username,
          password: credentials.password,
        };
      }

      console.log('[BiometricService] No credentials found');
      return null;
    } catch (error: any) {
      console.error('[BiometricService] Error retrieving credentials:', error);
      return null;
    }
  }

  /**
   * Delete stored credentials from device keychain
   * @returns True if credentials were deleted successfully
   */
  async deleteStoredCredentials(): Promise<boolean> {
    try {
      const result = await Keychain.resetGenericPassword({
        service: KEYCHAIN_SERVICE,
      });

      if (result) {
        console.log('[BiometricService] Credentials deleted successfully');
        return true;
      } else {
        console.log('[BiometricService] No credentials to delete');
        return true; // Not an error, just no credentials existed
      }
    } catch (error: any) {
      console.error('[BiometricService] Error deleting credentials:', error);
      return false;
    }
  }

  /**
   * Check if credentials are currently stored
   * @returns True if credentials exist in keychain
   */
  async hasStoredCredentials(): Promise<boolean> {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: KEYCHAIN_SERVICE,
      });

      return credentials !== false;
    } catch (error: any) {
      console.error('[BiometricService] Error checking stored credentials:', error);
      return false;
    }
  }

  /**
   * Create biometric keys for signature-based authentication
   * This is for advanced use cases requiring cryptographic signatures
   * @returns True if keys were created successfully
   */
  async createBiometricKeys(): Promise<boolean> {
    try {
      const { publicKey } = await this.rnBiometrics.createKeys();

      if (publicKey) {
        console.log('[BiometricService] Biometric keys created successfully');
        return true;
      }

      return false;
    } catch (error: any) {
      console.error('[BiometricService] Error creating biometric keys:', error);
      return false;
    }
  }

  /**
   * Delete biometric keys
   * @returns True if keys were deleted successfully
   */
  async deleteBiometricKeys(): Promise<boolean> {
    try {
      const { keysDeleted } = await this.rnBiometrics.deleteKeys();
      console.log('[BiometricService] Biometric keys deleted:', keysDeleted);
      return keysDeleted;
    } catch (error: any) {
      console.error('[BiometricService] Error deleting biometric keys:', error);
      return false;
    }
  }

  /**
   * Check if biometric keys exist
   * @returns True if keys exist
   */
  async biometricKeysExist(): Promise<boolean> {
    try {
      const { keysExist } = await this.rnBiometrics.biometricKeysExist();
      return keysExist;
    } catch (error: any) {
      console.error('[BiometricService] Error checking biometric keys:', error);
      return false;
    }
  }

  /**
   * Create cryptographic signature with biometric authentication
   * Advanced feature for enhanced security
   * @param payload String to sign
   * @returns Signature string or null if failed
   */
  async createSignature(payload: string): Promise<string | null> {
    try {
      const { success, signature } = await this.rnBiometrics.createSignature({
        promptMessage: 'Sign in',
        payload,
      });

      if (success && signature) {
        console.log('[BiometricService] Signature created successfully');
        return signature;
      }

      return null;
    } catch (error: any) {
      console.error('[BiometricService] Error creating signature:', error);
      return null;
    }
  }

  /**
   * Get detailed biometric capability information
   * @returns Object with biometric capabilities
   */
  async getBiometricCapabilities(): Promise<{
    available: boolean;
    biometryType: string | null;
    hasEnrolledBiometrics: boolean;
  }> {
    try {
      const { available, biometryType } = await this.rnBiometrics.isSensorAvailable();

      return {
        available,
        biometryType: await this.getBiometricType(),
        hasEnrolledBiometrics: available && biometryType !== null,
      };
    } catch (error: any) {
      console.error('[BiometricService] Error getting biometric capabilities:', error);
      return {
        available: false,
        biometryType: null,
        hasEnrolledBiometrics: false,
      };
    }
  }

  /**
   * Disable biometric authentication for the user
   * Deletes both keys and stored credentials
   * @returns True if biometrics were disabled successfully
   */
  async disableBiometricAuth(): Promise<boolean> {
    try {
      const credentialsDeleted = await this.deleteStoredCredentials();
      const keysDeleted = await this.deleteBiometricKeys();

      const success = credentialsDeleted && keysDeleted;
      console.log('[BiometricService] Biometric auth disabled:', success);
      return success;
    } catch (error: any) {
      console.error('[BiometricService] Error disabling biometric auth:', error);
      return false;
    }
  }

  /**
   * Enable biometric authentication for the user
   * Stores credentials and creates keys
   * @param email User email
   * @param password User password
   * @returns True if biometrics were enabled successfully
   */
  async enableBiometricAuth(email: string, password: string): Promise<boolean> {
    try {
      // First verify biometrics are available
      const available = await this.isBiometricAvailable();
      if (!available) {
        throw new Error('Biometric authentication is not available on this device');
      }

      // Authenticate user with biometrics
      const authenticated = await this.authenticate(
        'Enable Biometric Login',
        'Authenticate to enable biometric login'
      );

      if (!authenticated) {
        throw new Error('Biometric authentication failed');
      }

      // Store credentials
      const credentialsStored = await this.storeCredentials(email, password);
      if (!credentialsStored) {
        throw new Error('Failed to store credentials securely');
      }

      // Create biometric keys (optional, for advanced security)
      await this.createBiometricKeys();

      console.log('[BiometricService] Biometric auth enabled successfully');
      return true;
    } catch (error: any) {
      console.error('[BiometricService] Error enabling biometric auth:', error);
      // Cleanup on failure
      await this.disableBiometricAuth();
      throw error;
    }
  }
}

// Export singleton instance
export const BiometricService = new BiometricAuthService();
