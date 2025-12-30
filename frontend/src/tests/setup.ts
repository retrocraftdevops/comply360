import { expect, afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/svelte';
import * as matchers from '@testing-library/jest-dom/matchers';

// Extend Vitest's expect with jest-dom matchers
if (matchers && typeof matchers === 'object') {
    expect.extend(matchers);
}

// Cleanup after each test
afterEach(() => {
	cleanup();
});

// Mock SvelteKit modules
vi.mock('$app/environment', () => ({
	browser: false,
	dev: true,
	building: false,
	version: 'test',
}));

vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	invalidate: vi.fn(),
	invalidateAll: vi.fn(),
	preloadData: vi.fn(),
	preloadCode: vi.fn(),
	beforeNavigate: vi.fn(),
	afterNavigate: vi.fn(),
}));

vi.mock('$app/stores', () => {
	const getStores = () => {
		const navigating = { subscribe: vi.fn() };
		const page = {
			subscribe: vi.fn(() => () => {}),
		};
		const session = { subscribe: vi.fn() };
		const updated = { subscribe: vi.fn(), check: vi.fn() };

		return { navigating, page, session, updated };
	};

	const page = {
		subscribe: vi.fn(() => () => {}),
	};
	const navigating = { subscribe: vi.fn() };
	const updated = { subscribe: vi.fn(), check: vi.fn() };

	return { getStores, navigating, page, updated };
});

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
	writable: true,
	value: vi.fn().mockImplementation((query) => ({
		matches: false,
		media: query,
		onchange: null,
		addListener: vi.fn(),
		removeListener: vi.fn(),
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		dispatchEvent: vi.fn(),
	})),
});

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
	constructor() {}
	disconnect() {}
	observe() {}
	takeRecords() {
		return [];
	}
	unobserve() {}
} as any;

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
	constructor() {}
	disconnect() {}
	observe() {}
	unobserve() {}
} as any;
