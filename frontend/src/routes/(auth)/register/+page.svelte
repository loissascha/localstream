<script lang="ts">
	import { resolve } from '$app/paths';
	import { API_URL } from '$lib/consts';

	let username = $state('');
	let loading = $state(false);
	let error: string | null = null;

	async function submit() {
		loading = true;

		try {
			const res = await fetch(API_URL + '/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username: username })
			});

			if (!res.ok) {
				throw new Error('Request failed');
			}

			resolve('/(auth)/login');
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

<main
	class="grid min-h-dvh place-items-center bg-[radial-gradient(circle_at_20%_10%,#d7e8ff_0%,transparent_40%),radial-gradient(circle_at_80%_0%,#c8f5e9_0%,transparent_32%),linear-gradient(180deg,#eef2f7_0%,#dce6f2_100%)] px-4 py-6"
>
	<section
		class="w-full max-w-sm rounded-2xl border border-slate-900/10 bg-white/85 p-6 shadow-lg shadow-slate-900/10"
	>
		<h1 class="m-0 text-2xl font-semibold text-slate-900">Register</h1>

		{#if loading}
			<p>Loading</p>
		{:else}
			<form on:submit|preventDefault={submit}>
				<input type="text" name="username" bind:value={username} class="border" />
				<button>Submit</button>
			</form>
		{/if}

		<p>Want to login? Go to <a href={resolve('/(auth)/login')}>login page</a></p>
	</section>
</main>
