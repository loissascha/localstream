<script lang="ts">
	import type { ShowInfo } from '$lib/types/export_types';
	import type { Snippet } from 'svelte';
	import ShowMetadataSearchOverlay from './overlays/ShowMetadataSearchOverlay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';
	import { auth } from '$lib/auth.svelte';

	interface Props {
		children: Snippet;
		show: ShowInfo;
	}
	let { children, show }: Props = $props();

	let showMetadataOverlayOpen = $state(false);
</script>

<ContextMenu closeOnItemClick={true}>
	{@render children()}
	{#snippet items(closeMenu)}
		{#if auth.isAdmin}
			<button
				role="menuitem"
				onclick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					closeMenu();
					showMetadataOverlayOpen = true;
				}}
				class="cursor-pointer px-4 py-2 hover:bg-neutral-700">Update Metadata</button
			>
		{/if}
	{/snippet}
</ContextMenu>

{#if showMetadataOverlayOpen}
	<ShowMetadataSearchOverlay
		{show}
		close={() => {
			showMetadataOverlayOpen = false;
		}}
	/>
{/if}
