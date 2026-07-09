<script lang="ts">
	import { page } from '$app/state';
	import { loadLibraries } from '$lib/api/libraries';
	import { auth } from '$lib/auth.svelte';
	import AdminLibraryForm from '$lib/components/admin/AdminLibraryForm.svelte';
	import type { LibraryListItem } from '$lib/types/export_types';

	const libraryId = $derived(page.params.libraryID ?? '');
	let library = $state<LibraryListItem | null>(null);

	let libraries = $state<LibraryListItem[]>([]);

	async function loadLibs() {
		if (!auth.token) {
			return;
		}
		const data = await loadLibraries(auth.token);
		libraries = data.libraries;
	}

	$effect(() => {
		if (!auth.token) return;
		loadLibs();
	});

	$effect(() => {
		libraryId;
		libraries;
		if (libraryId != '') {
			for (var lib of libraries) {
				if (lib.id == libraryId) {
					library = lib;
				}
			}
		}
	});
</script>

{#if library != null}
	<AdminLibraryForm {library} />
{/if}
