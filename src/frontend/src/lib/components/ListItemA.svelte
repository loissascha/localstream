<script lang="ts">
	let { href, children } = $props();
</script>

<a
	{href}
	class="flex cursor-pointer flex-col rounded-lg transition-all duration-200 hover:scale-105 justify-between gap-2"
>
	{@render children()}
</a>
