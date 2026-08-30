<script lang="ts">
	let { children } = $props();
</script>

<div
	class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
>
	{@render children()}
</div>
