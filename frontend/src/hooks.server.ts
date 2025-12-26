import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// Pass through all /api requests to the proxy (don't let SvelteKit handle them)
	if (event.url.pathname.startsWith('/api')) {
		// In dev mode, Vite proxy will handle this
		// In production, this will be handled by the reverse proxy (nginx/caddy)
		return resolve(event);
	}

	return resolve(event);
};

