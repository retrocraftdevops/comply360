import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				// Modern, clean color palette inspired by vert.sh
				background: 'hsl(var(--background))',
				foreground: 'hsl(var(--foreground))',
				primary: {
					DEFAULT: 'hsl(221 83% 53%)',
					50: 'hsl(221 83% 98%)',
					100: 'hsl(221 83% 95%)',
					200: 'hsl(221 83% 90%)',
					300: 'hsl(221 83% 80%)',
					400: 'hsl(221 83% 65%)',
					500: 'hsl(221 83% 53%)',
					600: 'hsl(221 83% 45%)',
					700: 'hsl(221 83% 38%)',
					800: 'hsl(221 83% 30%)',
					900: 'hsl(221 83% 20%)',
					950: 'hsl(221 83% 10%)'
				},
				accent: {
					DEFAULT: 'hsl(221 83% 53%)',
					foreground: 'hsl(0 0% 98%)'
				},
				muted: {
					DEFAULT: 'hsl(210 40% 96%)',
					foreground: 'hsl(215 16% 47%)'
				},
				border: 'hsl(214 32% 91%)',
				input: 'hsl(214 32% 91%)',
				ring: 'hsl(221 83% 53%)',
				card: {
					DEFAULT: 'hsl(0 0% 100%)',
					foreground: 'hsl(222 47% 11%)'
				}
			},
			borderRadius: {
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)'
			},
			keyframes: {
				'fade-in': {
					'0%': { opacity: '0', transform: 'translateY(10px)' },
					'100%': { opacity: '1', transform: 'translateY(0)' }
				},
				'fade-out': {
					'0%': { opacity: '1', transform: 'translateY(0)' },
					'100%': { opacity: '0', transform: 'translateY(10px)' }
				},
				'slide-in': {
					'0%': { transform: 'translateX(-100%)' },
					'100%': { transform: 'translateX(0)' }
				},
				'slide-out': {
					'0%': { transform: 'translateX(0)' },
					'100%': { transform: 'translateX(-100%)' }
				},
				'scale-in': {
					'0%': { transform: 'scale(0.95)', opacity: '0' },
					'100%': { transform: 'scale(1)', opacity: '1' }
				},
				'scale-out': {
					'0%': { transform: 'scale(1)', opacity: '1' },
					'100%': { transform: 'scale(0.95)', opacity: '0' }
				},
				'shimmer': {
					'0%': { backgroundPosition: '-1000px 0' },
					'100%': { backgroundPosition: '1000px 0' }
				}
			},
			animation: {
				'fade-in': 'fade-in 0.3s ease-out',
				'fade-out': 'fade-out 0.3s ease-out',
				'slide-in': 'slide-in 0.3s ease-out',
				'slide-out': 'slide-out 0.3s ease-out',
				'scale-in': 'scale-in 0.2s ease-out',
				'scale-out': 'scale-out 0.2s ease-out',
				'shimmer': 'shimmer 2s infinite linear'
			},
			boxShadow: {
				'soft': '0 2px 8px rgba(0, 0, 0, 0.04)',
				'medium': '0 4px 16px rgba(0, 0, 0, 0.08)',
				'large': '0 8px 32px rgba(0, 0, 0, 0.12)',
				'glow': '0 0 20px rgba(59, 130, 246, 0.3)'
			},
			backdropBlur: {
				xs: '2px'
			},
			fontFamily: {
				sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', 'sans-serif']
			}
		}
	},
	plugins: []
} satisfies Config;
