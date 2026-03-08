<script lang="ts">
	let loadingShowData = $state(true);
	let loadingSeasons = $state(true);
</script>

<main class="min-h-dvh px-5 py-5">
	{#if loadingShowData}
		<p>Loading stuff...</p>
	{:else}
		<p>Show:</p>
	{/if}
</main>
