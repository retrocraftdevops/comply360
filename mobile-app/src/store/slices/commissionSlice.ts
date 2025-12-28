/**
 * Commission Slice
 * Manages commission state and filters
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface CommissionState {
  currentFilter: 'all' | 'pending' | 'approved' | 'paid' | 'disputed';
  dateRange: {
    start: string | null;
    end: string | null;
  };
  selectedCommissionId: string | null;
}

const initialState: CommissionState = {
  currentFilter: 'all',
  dateRange: {
    start: null,
    end: null,
  },
  selectedCommissionId: null,
};

const commissionSlice = createSlice({
  name: 'commission',
  initialState,
  reducers: {
    setFilter: (state, action: PayloadAction<CommissionState['currentFilter']>) => {
      state.currentFilter = action.payload;
    },
    setDateRange: (state, action: PayloadAction<{ start: string | null; end: string | null }>) => {
      state.dateRange = action.payload;
    },
    selectCommission: (state, action: PayloadAction<string | null>) => {
      state.selectedCommissionId = action.payload;
    },
    clearFilters: (state) => {
      state.currentFilter = 'all';
      state.dateRange = { start: null, end: null };
    },
  },
});

export const {
  setFilter,
  setDateRange,
  selectCommission,
  clearFilters,
} = commissionSlice.actions;

export default commissionSlice.reducer;
