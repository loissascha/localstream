import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const backendOrigin = process.env.VITE_API_URL ?? 'http://localhost:42069';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: backendOrigin,
				changeOrigin: true
			},
			'/static': {
				target: backendOrigin,
				changeOrigin: true
			}
		}
	}
});
