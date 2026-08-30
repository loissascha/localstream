<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { loadMoviesDatabase, movies } from '$lib/movies.svelte';
	import { loadShowsDatabase, shows } from '$lib/shows.svelte';

	let { children } = $props();

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		if (!movies.initialized) {
			loadMoviesDatabase().catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
		}
		if (!shows.initialized) {
			loadShowsDatabase().catch((e) => {
				const m = (e as Error).message;
				alert(m);
			});
		}
	});
</script>

{@render children()}
