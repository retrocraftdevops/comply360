/**
 * Document Slice
 * Manages document state and filters
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface DocumentState {
  currentCategory: string | null;
  searchQuery: string;
  selectedDocumentId: string | null;
  uploadProgress: Record<string, number>; // Track upload progress for multiple files
}

const initialState: DocumentState = {
  currentCategory: null,
  searchQuery: '',
  selectedDocumentId: null,
  uploadProgress: {},
};

const documentSlice = createSlice({
  name: 'document',
  initialState,
  reducers: {
    setCategory: (state, action: PayloadAction<string | null>) => {
      state.currentCategory = action.payload;
    },
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
    selectDocument: (state, action: PayloadAction<string | null>) => {
      state.selectedDocumentId = action.payload;
    },
    setUploadProgress: (state, action: PayloadAction<{ fileId: string; progress: number }>) => {
      state.uploadProgress[action.payload.fileId] = action.payload.progress;
    },
    removeUploadProgress: (state, action: PayloadAction<string>) => {
      delete state.uploadProgress[action.payload];
    },
    clearFilters: (state) => {
      state.currentCategory = null;
      state.searchQuery = '';
    },
  },
});

export const {
  setCategory,
  setSearchQuery,
  selectDocument,
  setUploadProgress,
  removeUploadProgress,
  clearFilters,
} = documentSlice.actions;

export default documentSlice.reducer;
