/**
 * Registration Slice
 * Manages registration state and filters
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface RegistrationState {
  currentFilter: 'all' | 'draft' | 'pending' | 'in_progress' | 'completed';
  searchQuery: string;
  selectedRegistrationId: string | null;
}

const initialState: RegistrationState = {
  currentFilter: 'all',
  searchQuery: '',
  selectedRegistrationId: null,
};

const registrationSlice = createSlice({
  name: 'registration',
  initialState,
  reducers: {
    setFilter: (state, action: PayloadAction<RegistrationState['currentFilter']>) => {
      state.currentFilter = action.payload;
    },
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
    selectRegistration: (state, action: PayloadAction<string | null>) => {
      state.selectedRegistrationId = action.payload;
    },
    clearFilters: (state) => {
      state.currentFilter = 'all';
      state.searchQuery = '';
    },
  },
});

export const {
  setFilter,
  setSearchQuery,
  selectRegistration,
  clearFilters,
} = registrationSlice.actions;

export default registrationSlice.reducer;
