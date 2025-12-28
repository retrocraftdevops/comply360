/**
 * Network Status Service
 * Detects and monitors network connectivity
 */

import NetInfo, { NetInfoState } from '@react-native-community/netinfo';

type NetworkStatusCallback = (isOnline: boolean) => void;

class NetworkStatus {
  private isOnline = true;
  private listeners: NetworkStatusCallback[] = [];
  private unsubscribe: (() => void) | null = null;

  /**
   * Initialize network status monitoring
   */
  async initialize(): Promise<void> {
    // Get initial state
    const state = await NetInfo.fetch();
    this.isOnline = state.isConnected ?? false;
    console.log(`[NetworkStatus] Initial state: ${this.isOnline ? 'online' : 'offline'}`);

    // Listen for changes
    this.unsubscribe = NetInfo.addEventListener(this.handleNetworkChange.bind(this));
  }

  /**
   * Handle network state changes
   */
  private handleNetworkChange(state: NetInfoState): void {
    const wasOnline = this.isOnline;
    this.isOnline = state.isConnected ?? false;

    if (wasOnline !== this.isOnline) {
      console.log(`[NetworkStatus] State changed: ${this.isOnline ? 'online' : 'offline'}`);
      this.notifyListeners();
    }
  }

  /**
   * Check if currently online
   */
  getIsOnline(): boolean {
    return this.isOnline;
  }

  /**
   * Subscribe to network status changes
   */
  subscribe(callback: NetworkStatusCallback): () => void {
    this.listeners.push(callback);

    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter((cb) => cb !== callback);
    };
  }

  /**
   * Notify all listeners of status change
   */
  private notifyListeners(): void {
    this.listeners.forEach((callback) => {
      try {
        callback(this.isOnline);
      } catch (error) {
        console.error('[NetworkStatus] Listener error:', error);
      }
    });
  }

  /**
   * Clean up listeners
   */
  cleanup(): void {
    if (this.unsubscribe) {
      this.unsubscribe();
      this.unsubscribe = null;
    }
    this.listeners = [];
    console.log('[NetworkStatus] Cleaned up');
  }
}

export default new NetworkStatus();
