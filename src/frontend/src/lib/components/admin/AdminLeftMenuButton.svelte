<script lang="ts">
	let { children, href } = $props();
</script>

<a class="cursor-pointer bg-neutral-800/70 hover:bg-neutral-800 rounded py-2 px-4" href={href}>{@render children()}</a>
