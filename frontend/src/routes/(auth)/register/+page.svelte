<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth.svelte';

	let username = $state('');
	let loading = $state(false);
	let error: string | null = null;

	async function submit() {
		loading = true;

		try {
			const res = await fetch('/api/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username: username })
			});

			if (!res.ok) {
				throw new Error('Request failed');
			}

			goto(resolve('/(auth)/login'));
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) return;

		goto(resolve('/(protected)'));
	});
</script>

<main class="grid min-h-dvh place-items-center px-4 py-6">
	<section
		class="w-full max-w-sm rounded-2xl border border-neutral-500 bg-neutral-800 p-6 shadow-lg shadow-neutral-300/30"
	>
		<h1 class="m-0 text-2xl font-semibold">Register</h1>

		{#if loading}
			<p>Loading</p>
		{:else}
			<form on:submit|preventDefault={submit}>
				<input
					type="text"
					name="username"
					bind:value={username}
					class="my-3 rounded border border-neutral-500 p-2"
				/>
				<button class="rounded border border-neutral-500 px-4 py-2">Submit</button>
			</form>
		{/if}

		<p>Want to login? Go to <a href={resolve('/(auth)/login')} class="text-blue-500">login page</a></p>
	</section>
</main>
