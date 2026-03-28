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
	});
</script>

<div>
	<section id="header" class="flex justify-between bg-neutral-900 px-4 py-4">
		<div>Layout User</div>
		<div>
			{#if auth.isAdmin}
				<a
					href={resolve('/(protected)/(admin)/admin')}
					type="submit"
					class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
				>
					Admin
				</a>
			{/if}
			<a
				href={resolve('/logout')}
				type="submit"
				class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
			>
				Log out
			</a>
		</div>
	</section>
	<div>
		{@render children()}
	</div>
</div>
