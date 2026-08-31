<script lang="ts">
	let { href, children } = $props();
</script>

<a
	{href}
	class="flex cursor-pointer flex-col rounded-lg justify-between"
>
	{@render children()}
</a>
