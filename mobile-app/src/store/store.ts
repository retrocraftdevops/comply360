/**
 * Redux Store Configuration
 * Centralized state management for Comply360 Mobile
 */

import { configureStore, combineReducers } from '@reduxjs/toolkit';
import {
  persistStore,
  persistReducer,
  FLUSH,
  REHYDRATE,
  PAUSE,
  PERSIST,
  PURGE,
  REGISTER,
} from 'redux-persist';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Import slices
import authReducer from './slices/authSlice';
import registrationReducer from './slices/registrationSlice';
import documentReducer from './slices/documentSlice';
import commissionReducer from './slices/commissionSlice';
import notificationReducer from './slices/notificationSlice';
import uiReducer from './slices/uiSlice';

// Import API slices
import { authApi } from './api/authApi';
import { registrationApi } from './api/registrationApi';
import { documentApi } from './api/documentApi';
import { commissionApi } from './api/commissionApi';

const rootReducer = combineReducers({
  auth: authReducer,
  registration: registrationReducer,
  document: documentReducer,
  commission: commissionReducer,
  notification: notificationReducer,
  ui: uiReducer,
  // RTK Query API reducers
  [authApi.reducerPath]: authApi.reducer,
  [registrationApi.reducerPath]: registrationApi.reducer,
  [documentApi.reducerPath]: documentApi.reducer,
  [commissionApi.reducerPath]: commissionApi.reducer,
});

const persistConfig = {
  key: 'root',
  version: 1,
  storage: AsyncStorage,
  whitelist: ['auth', 'ui'], // Only persist these reducers
  blacklist: [
    authApi.reducerPath,
    registrationApi.reducerPath,
    documentApi.reducerPath,
    commissionApi.reducerPath,
  ],
};

const persistedReducer = persistReducer(persistConfig, rootReducer);

export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }).concat(
      authApi.middleware,
      registrationApi.middleware,
      documentApi.middleware,
      commissionApi.middleware,
    ),
});

export const persistor = persistStore(store);

// Export types
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Export typed hooks
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
