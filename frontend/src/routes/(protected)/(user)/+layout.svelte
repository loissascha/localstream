<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';

	let { children } = $props();
</script>

<div>
	<section id="header" class="flex justify-between bg-neutral-900 px-4 py-4">
		<div>
			<a href={resolve('/(protected)/(user)')}>Localstream</a>
		</div>
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
