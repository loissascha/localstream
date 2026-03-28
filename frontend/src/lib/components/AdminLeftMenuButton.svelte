<script lang="ts">
	let { children, href } = $props();
</script>

<a class="cursor-pointer bg-neutral-700/60 hover:bg-neutral-700 rounded py-2 px-4" href={href}>{@render children()}</a>
