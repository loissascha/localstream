<script lang="ts">
	import { searchShowMetadata, setPrimaryShowMetadataByFetchID } from '$lib/api/show_metadata';
	import { loadShows } from '$lib/api/shows';
	import { auth } from '$lib/auth.svelte';
	import SearchIcon from '$lib/icons/SearchIcon.svelte';
	import type { ShowInfo, ShowSearchResult } from '$lib/types/export_types';
	import Overlay from './Overlay.svelte';

	interface Props {
		close: () => void;
		show: ShowInfo;
	}
	let { close, show }: Props = $props();

	let searchingMetadata = $state(false);
	let error_message = $state('');
	let success_message = $state('');
	let searchQuery = $state('');
	let searchResults = $state<ShowSearchResult[]>([]);

	async function submitMetadataSearchForm() {
		try {
			if (!auth.token) return;
			searchingMetadata = true;
			success_message = '';
			error_message = '';
			if (searchQuery == '') {
				alert('No query');
				return;
			}
			searchResults = await searchShowMetadata(auth.token, searchQuery);
			searchingMetadata = false;
		} catch (e) {
			const m = (e as Error).message;
			error_message = m;
		}
	}
</script>

<Overlay {close}>
	<h1 class="text-2xl font-bold tracking-wide">Metadata</h1>
	<h2 class="text-lg font-semibold tracking-wide">{show.name}</h2>
	<div class="flex">
		<div class="grow">
			{show.description}
		</div>
		<div>
			<img src={show.medium_image_url} class="w-120" alt={show.name} />
		</div>
	</div>
	{#if success_message != ''}
		<div class="mt-4 text-green-500">
			{success_message}
		</div>
	{/if}
	{#if error_message != ''}
		<div class="mt-4 text-red-500">
			{error_message}
		</div>
	{/if}

	<form
		onsubmit={(e) => {
			e.preventDefault();
			if (searchingMetadata) return;
			submitMetadataSearchForm();
		}}
		class="flex items-center gap-2"
	>
		<input
			bind:value={searchQuery}
			type="text"
			class="my-4 w-full rounded bg-neutral-700 px-4 py-2"
			placeholder="Search for new Metadata"
		/>
		<button class="flex cursor-pointer gap-2 rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600">
			<SearchIcon /> Search
		</button>
	</form>

	{#if searchingMetadata}
		Searching...
	{:else}
		{#each searchResults as result (result.show.id)}
			<div class="my-4 p-4">
				<div class="grid grid-cols-1 md:grid-cols-2">
					<div>
						<div class="font-bold">{result.show.name} ({result.show.premiered})</div>
						<div class="my-2">
							{result.show.summary}
						</div>
						<div>
							<button
								onclick={() => {
									if (!auth.token) return;
									setPrimaryShowMetadataByFetchID(auth.token, show.id, result.show.id)
										.then(() => {
											searchResults = [];
											searchQuery = '';
											success_message = 'Updated primary metadata!';

											// TODO: invent reloadSingleShow
											// reloadSingleShow(show.id);
											if (auth.token) {
												loadShows(auth.token);
											}
										})
										.catch((e) => {
											const m = (e as Error).message;
											alert(m);
										});
								}}
								class="mt-4 mb-4 cursor-pointer rounded bg-neutral-700 px-4 py-2 hover:bg-neutral-600"
								>Select as Primary</button
							>
						</div>
					</div>
					<div>
						{#if result.show.image && result.show.image.medium != ''}
							<img alt={'show ' + result.show.id} src={`${result.show.image.medium}`} />
						{/if}
					</div>
				</div>
			</div>
		{/each}
	{/if}
</Overlay>
