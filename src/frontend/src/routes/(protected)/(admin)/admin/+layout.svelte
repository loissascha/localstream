<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import AdminLeftMenuButton from '$lib/components/admin/AdminLeftMenuButton.svelte';
	import Logo from '$lib/components/ui/Logo.svelte';
	import LayoutUserDropdown from '$lib/components/ui/LayoutUserDropdown.svelte';
	import MenuIcon from '$lib/icons/MenuIcon.svelte';
	import DropdownMelt from '$lib/components/ui/DropdownMelt.svelte';
	import DropdownItem from '$lib/components/ui/DropdownItem.svelte';

	let { children } = $props();

	$effect(() => {
		if (!auth.initialized) {
			return;
		}
	});
</script>

<div class="flex h-dvh flex-col">
	<section id="header" class="flex items-center justify-between px-4 py-4">
		<div class="flex items-center gap-4">
			<a
				href={resolve('/(protected)/(admin)/admin')}
				class="flex items-center gap-2 text-2xl font-semibold tracking-wider select-none"
				><Logo />ocalstream</a
			>
			<DropdownMelt className="lg:hidden">
				<MenuIcon />
				{#snippet items()}
					<DropdownItem href={resolve('/(protected)/(user)')}>
						<div class="flex items-center gap-1">Home</div>
					</DropdownItem>
					{#if auth.isAdmin}
						<div class="flex flex-col gap-2">
							<span class="px-2 tracking-wide text-neutral-300 uppercase">Admin</span>
							<DropdownItem href={resolve('/(protected)/(admin)/admin')}>Dashboard</DropdownItem>
							<DropdownItem href={resolve('/(protected)/(admin)/admin/libraries')}
								>Libraries</DropdownItem
							>
							<DropdownItem href={resolve('/(protected)/(admin)/admin/users')}>Users</DropdownItem>
							<DropdownItem href={resolve('/(protected)/(admin)/admin/metadata')}
								>Metadata</DropdownItem
							>
						</div>
					{/if}
				{/snippet}
			</DropdownMelt>
		</div>
		<div class="flex items-center gap-2">
			<LayoutUserDropdown />
		</div>
	</section>
	<div class="block min-h-0 grow lg:flex">
		<section
			id="left"
			class="hidden w-80 shrink-0 grow-0 overflow-y-auto border-r border-r-neutral-600 p-4 lg:block"
		>
			<div class="mb-4 flex flex-col gap-2 border-b border-neutral-600 pb-4">
				<AdminLeftMenuButton href={resolve('/(protected)/(user)')}>Home</AdminLeftMenuButton>
			</div>
			{#if auth.isAdmin}
				<div class="flex flex-col gap-2">
					<span class="tracking-wide text-neutral-300 uppercase">Admin</span>
					<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin')}
						>Dashboard</AdminLeftMenuButton
					>
					<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin/libraries')}
						>Libraries</AdminLeftMenuButton
					>
					<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin/users')}
						>Users</AdminLeftMenuButton
					>
					<AdminLeftMenuButton href={resolve('/(protected)/(admin)/admin/metadata')}
						>Metadata</AdminLeftMenuButton
					>
				</div>
			{/if}
		</section>
		<section id="content" class="grow overflow-y-auto p-4">
			{@render children()}
		</section>
	</div>
</div>
