<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ShowInfo } from '$lib/types/export_types';
	import ListItemA from './ListItemA.svelte';
	import ShowInfoDisplay from './ShowInfoDisplay.svelte';
	import ContextMenu from './ui/ContextMenu.svelte';
	import ShowMetadataSearchOverlay from './overlays/ShowMetadataSearchOverlay.svelte';

	interface Props {
		show: ShowInfo;
		selectable?: boolean;
		selected?: boolean;
		showFinished?: boolean;
	}

	let {
		show,
		selectable = false,
		selected = $bindable(false),
		showFinished = false
	}: Props = $props();

	let showMetadataOverlayOpen = $state(false);
</script>

<div class="relative">
	<ContextMenu closeOnItemClick={true}>
		<ListItemA href={resolve('/(protected)/(user)/shows/[showID]', { showID: show.id })}>
			<ShowInfoDisplay {show} />
		</ListItemA>
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

	{#if showFinished}
		<div class="absolute top-2 right-2 z-10">
			<span
				class={`flex h-7 w-7 items-center justify-center rounded-full border bg-neutral-950/80 text-brand shadow-sm transition-all duration-150`}
			>
				<svg aria-hidden="true" class="h-4 w-4" viewBox="0 0 16 16" fill="none">
					<path
						d="M3.5 8.5L6.5 11.5L12.5 4.5"
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
					/>
				</svg>
			</span>
		</div>
	{/if}

	{#if selectable}
		<button
			type="button"
			class="absolute top-2 left-2 z-10"
			role="checkbox"
			aria-checked={selected}
			aria-label={`Select ${show.name}`}
			onclick={(event: MouseEvent) => {
				event.preventDefault();
				event.stopPropagation();
				selected = !selected;
			}}
		>
			<span
				class={`flex h-7 w-7 items-center justify-center rounded-full border shadow-sm transition-all duration-150 ${selected ? 'border-brand bg-brand text-white' : 'border-neutral-500/80 bg-neutral-950/85 text-transparent hover:border-neutral-300 hover:bg-neutral-900'}`}
			>
				<svg aria-hidden="true" class="h-4 w-4" viewBox="0 0 16 16" fill="none">
					<path
						d="M3.5 8.5L6.5 11.5L12.5 4.5"
						stroke="currentColor"
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
					/>
				</svg>
			</span>
		</button>
	{/if}
</div>

{#if showMetadataOverlayOpen}
	<ShowMetadataSearchOverlay
		{show}
		close={() => {
			showMetadataOverlayOpen = false;
		}}
	/>
{/if}
