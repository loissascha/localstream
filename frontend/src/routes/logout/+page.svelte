<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { clearAuth } from '$lib/auth.svelte';
	import { deleteCookie } from '$lib/cookies';

	$effect(() => {
		deleteCookie('bearer');
		clearAuth();

		goto(resolve('/(auth)/login'));
	});
</script>
