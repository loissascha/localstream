<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import AdminLeftMenuButton from '$lib/components/AdminLeftMenuButton.svelte';

	let { children } = $props();

	$effect(() => {
		if (!auth.initialized) return;
		if (!auth.isAdmin) {
			goto(resolve('/(protected)'));
		}
	});
</script>

<div class="flex h-dvh flex-col">
	<section id="header" class="flex items-center justify-between bg-neutral-900 px-4 py-4">
		<div>
			<a href={resolve('/(protected)/(user)')} class="text-2xl font-semibold tracking-wider"
				>Localstream Admin</a
			>
		</div>
		<div>
			<a
				href={resolve('/logout')}
				type="submit"
				class="cursor-pointer rounded-md border border-neutral-500 bg-neutral-600 px-3 py-1.5 text-sm hover:bg-neutral-500"
			>
				Log out
			</a>
		</div>
	</section>
	<div class="flex min-h-0 grow">
		<section id="left" class="w-80 shrink-0 grow-0 overflow-y-auto bg-neutral-800 p-4">
			<div class="flex flex-col gap-2">
				<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin')}
					>Dashboard</AdminLeftMenuButton
				>
				<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin/libraries')}
					>Libraries</AdminLeftMenuButton
				>
				<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin/users')}>Users</AdminLeftMenuButton
				>
			</div>
		</section>
		<section id="content" class="grow overflow-y-auto p-4">
			{@render children()}
		</section>
	</div>
</div>
