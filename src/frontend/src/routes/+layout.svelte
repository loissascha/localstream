<script lang="ts">
	import { browser } from '$app/environment';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { loadAuthFromCookies } from '$lib/auth.svelte';

	let { children } = $props();

	$effect(() => {
		loadAuthFromCookies();
	});

	$effect(() => {
		if (!browser) return;
		if (!('serviceWorker' in navigator)) return;

		void navigator.serviceWorker.register('/service-worker.js');
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="manifest" href="/manifest.webmanifest" />
	<link rel="apple-touch-icon" href="/icons/icon-192.png" />
	<meta name="theme-color" content="#171717" />
	<meta name="apple-mobile-web-app-capable" content="yes" />
	<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
	<meta name="apple-mobile-web-app-title" content="Localstream" />
	<title>Localstream</title>
</svelte:head>
<div class="min-h-dvh w-dvw bg-neutral-900 text-neutral-100">
	{@render children()}
</div>
