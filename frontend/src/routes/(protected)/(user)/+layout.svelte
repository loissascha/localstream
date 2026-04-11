<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import InstallAppButton from '$lib/components/InstallAppButton.svelte';
	import LogoutIcon from '$lib/icons/LogoutIcon.svelte';
	import SettingsIcon from '$lib/icons/SettingsIcon.svelte';

	let { children } = $props();
</script>

<section id="header" class="flex items-center justify-between px-4 py-4">
	<div class="flex items-center gap-2">
		<a href={resolve('/(protected)/(user)')} class="text-2xl font-semibold tracking-wider"
			>Localstream</a
		>
	</div>
	<div class="flex items-center gap-2">
		<InstallAppButton />
		{#if auth.isAdmin}
			<a
				href={resolve('/(protected)/(admin)/admin')}
				type="submit"
				class="cursor-pointer px-3 py-1.5 text-sm hover:text-brand"
			>
				<SettingsIcon />
			</a>
		{/if}
		<a
			href={resolve('/logout')}
			type="submit"
			class="cursor-pointer px-3 py-1.5 text-sm hover:text-brand"
		>
			<LogoutIcon />
		</a>
	</div>
</section>
<section id="seletions" class="flex items-center justify-center gap-4 p-4">
	<div class="flex gap-2 rounded-full bg-neutral-700 px-3 py-2">
		<a
			href={resolve('/(protected)/(user)')}
			class="cursor-pointer rounded-full bg-neutral-600 px-4 py-2 text-sm hover:bg-neutral-600"
		>
			Home
		</a>
		<a
			href={resolve('/(protected)/(user)/shows')}
			class="cursor-pointer rounded-full px-4 py-2 text-sm hover:bg-neutral-600"
		>
			Shows
		</a>
		<a
			href={resolve('/(protected)/(user)/movies')}
			class="cursor-pointer rounded-full px-4 py-2 text-sm hover:bg-neutral-600"
		>
			Movies
		</a>
	</div>
</section>
<section id="content" class="p-4">
	{@render children()}
</section>
