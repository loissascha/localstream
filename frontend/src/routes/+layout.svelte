<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { loadAuthFromCookies } from '$lib/auth.svelte';

	let { children } = $props();

	$effect(() => {
		loadAuthFromCookies();
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
