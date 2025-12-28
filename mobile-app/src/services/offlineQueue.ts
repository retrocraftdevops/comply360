/**
 * Offline Queue Service
 * Manages offline data synchronization queue
 */

import AsyncStorage from '@react-native-async-storage/async-storage';

export interface QueueItem {
  id: string;
  type: 'CREATE' | 'UPDATE' | 'DELETE';
  endpoint: string;
  data: any;
  timestamp: number;
  retryCount: number;
  maxRetries: number;
  status: 'pending' | 'processing' | 'failed' | 'completed';
}

const QUEUE_STORAGE_KEY = '@comply360_offline_queue';
const MAX_RETRIES = 3;

class OfflineQueue {
  private queue: QueueItem[] = [];
  private isProcessing = false;

  /**
   * Initialize queue from storage
   */
  async initialize(): Promise<void> {
    try {
      const storedQueue = await AsyncStorage.getItem(QUEUE_STORAGE_KEY);
      if (storedQueue) {
        this.queue = JSON.parse(storedQueue);
        console.log(`[OfflineQueue] Loaded ${this.queue.length} items from storage`);
      }
    } catch (error) {
      console.error('[OfflineQueue] Failed to load queue:', error);
    }
  }

  /**
   * Add item to queue
   */
  async add(item: Omit<QueueItem, 'id' | 'timestamp' | 'retryCount' | 'status'>): Promise<void> {
    const queueItem: QueueItem = {
      ...item,
      id: this.generateId(),
      timestamp: Date.now(),
      retryCount: 0,
      status: 'pending',
      maxRetries: item.maxRetries || MAX_RETRIES,
    };

    this.queue.push(queueItem);
    await this.persistQueue();
    console.log(`[OfflineQueue] Added item: ${queueItem.id} (${queueItem.type} ${queueItem.endpoint})`);
  }

  /**
   * Process all pending items in queue
   */
  async process(): Promise<void> {
    if (this.isProcessing) {
      console.log('[OfflineQueue] Already processing, skipping');
      return;
    }

    if (this.queue.length === 0) {
      console.log('[OfflineQueue] Queue is empty');
      return;
    }

    this.isProcessing = true;
    console.log(`[OfflineQueue] Processing ${this.queue.length} items`);

    const pendingItems = this.queue.filter(
      (item) => item.status === 'pending' || item.status === 'failed'
    );

    for (const item of pendingItems) {
      try {
        await this.processItem(item);
      } catch (error) {
        console.error(`[OfflineQueue] Failed to process item ${item.id}:`, error);
      }
    }

    this.isProcessing = false;
    await this.persistQueue();
    console.log('[OfflineQueue] Processing complete');
  }

  /**
   * Process a single queue item
   */
  private async processItem(item: QueueItem): Promise<void> {
    console.log(`[OfflineQueue] Processing item: ${item.id}`);

    item.status = 'processing';
    item.retryCount++;

    try {
      // Simulate API call (replace with actual API client)
      await this.sendToAPI(item);

      item.status = 'completed';
      console.log(`[OfflineQueue] Item ${item.id} completed successfully`);

      // Remove completed items from queue
      this.queue = this.queue.filter((i) => i.id !== item.id);
    } catch (error) {
      console.error(`[OfflineQueue] Item ${item.id} failed (attempt ${item.retryCount}):`, error);

      if (item.retryCount >= item.maxRetries) {
        item.status = 'failed';
        console.log(`[OfflineQueue] Item ${item.id} exceeded max retries, marking as failed`);
      } else {
        item.status = 'pending';
        console.log(`[OfflineQueue] Item ${item.id} will be retried`);
      }
    }
  }

  /**
   * Send queue item to API (placeholder)
   */
  private async sendToAPI(item: QueueItem): Promise<void> {
    // TODO: Replace with actual API client
    // Example:
    // const response = await apiClient.request({
    //   method: item.type === 'DELETE' ? 'DELETE' : item.type === 'CREATE' ? 'POST' : 'PUT',
    //   url: item.endpoint,
    //   data: item.data,
    // });

    // Simulate API call
    await new Promise((resolve, reject) => {
      setTimeout(() => {
        // Simulate 90% success rate
        if (Math.random() > 0.1) {
          resolve(true);
        } else {
          reject(new Error('Simulated API error'));
        }
      }, 1000);
    });
  }

  /**
   * Retry failed items
   */
  async retryFailed(): Promise<void> {
    const failedItems = this.queue.filter((item) => item.status === 'failed');
    console.log(`[OfflineQueue] Retrying ${failedItems.length} failed items`);

    for (const item of failedItems) {
      item.status = 'pending';
      item.retryCount = 0;
    }

    await this.persistQueue();
    await this.process();
  }

  /**
   * Clear all items from queue
   */
  async clear(): Promise<void> {
    this.queue = [];
    await this.persistQueue();
    console.log('[OfflineQueue] Queue cleared');
  }

  /**
   * Get queue status
   */
  getStatus(): {
    total: number;
    pending: number;
    processing: number;
    failed: number;
    completed: number;
  } {
    return {
      total: this.queue.length,
      pending: this.queue.filter((i) => i.status === 'pending').length,
      processing: this.queue.filter((i) => i.status === 'processing').length,
      failed: this.queue.filter((i) => i.status === 'failed').length,
      completed: this.queue.filter((i) => i.status === 'completed').length,
    };
  }

  /**
   * Get all queue items
   */
  getQueue(): QueueItem[] {
    return [...this.queue];
  }

  /**
   * Persist queue to storage
   */
  private async persistQueue(): Promise<void> {
    try {
      await AsyncStorage.setItem(QUEUE_STORAGE_KEY, JSON.stringify(this.queue));
    } catch (error) {
      console.error('[OfflineQueue] Failed to persist queue:', error);
    }
  }

  /**
   * Generate unique ID
   */
  private generateId(): string {
    return `queue_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

export default new OfflineQueue();
