<script lang="ts">
	import { resolve } from '$app/paths';
	import { setCookie } from '$lib/cookies';
	import type { AuthResponse, AuthUserResponse } from '$lib/types/export_types';
	import { auth, loadAuthFromCookies } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';

	let loading = $state(true);
	let data: AuthUserResponse[] = $state([]);
	let error: string | null = $state(null);

	async function load() {
		try {
			const res = await fetch('/api/auth/users/list');

			if (!res.ok) {
				throw new Error('Request failed');
			}

			data = await res.json();
			console.log('data:', data);
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	async function clickUser(username: string) {
		console.log('user clicked ', username);
		try {
			const res = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username: username })
			});

			if (!res.ok) {
				throw new Error('Request failed');
			}

			var response = (await res.json()) as AuthResponse;
			console.log(response);

			setCookie('bearer', response.token, 30);
			loadAuthFromCookies();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		load();
	});

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.loggedIn) return;

		goto(resolve('/(protected)'));
	});
</script>

<main
	class="grid min-h-dvh place-items-center px-4 py-6"
>
	<section
		class="w-full max-w-sm rounded-2xl border border-neutral-500 bg-neutral-800 p-6 shadow-lg shadow-neutral-300/30"
	>
		<h1 class="m-0 text-2xl font-semibold">Choose Profile</h1>
		{#if auth.loggedIn}
			<div>Already logged in</div>
		{:else}
			<div>Not logged in</div>
		{/if}
		{#if loading}
			<p>Loading</p>
		{:else if error}
			<p>Error: {error}</p>
		{:else}
			<div class="my-8 flex items-center justify-center gap-2">
				{#each data as item (item.id)}
					<button
						onclick={() => clickUser(item.username)}
						class="h-22 w-22 cursor-pointer place-content-center place-items-center rounded border border-neutral-500 bg-neutral-700 text-center shadow shadow-neutral-600"
					>
						{item.username}
					</button>
				{/each}
			</div>
		{/if}

		<p class="text-sm">
			No profile yet? <a class="cursor-pointer text-blue-500" href={resolve('/(auth)/register')}
				>create a new one</a
			>
		</p>
	</section>
</main>
