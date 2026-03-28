<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let { children } = $props();

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) {
			goto(resolve('/(auth)/login'));
			return;
		}
		if (!auth.isAdmin) {
			goto(resolve('/(protected)'));
		}
	});
</script>

<div>
	<div>Layout</div>
	<div>
		{@render children()}
	</div>
</div>
