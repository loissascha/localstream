<script lang="ts">
	import ShowMetadataSearchOverlay from './overlays/ShowMetadataSearchOverlay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';

	let { children, show } = $props();

	let showMetadataOverlayOpen = $state(false);
</script>

<ContextMenu closeOnItemClick={true}>
	{@render children()}
	{#snippet items(closeMenu)}
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
