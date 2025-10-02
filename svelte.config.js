import adapter from '@sveltejs/adapter-auto';
import { sveltekit } from '@sveltejs/kit/vite';

const config = {
  kit: { adapter: adapter() },
  vitePlugin: { inspector: false },
  plugins: [sveltekit()]
};

export default config;
