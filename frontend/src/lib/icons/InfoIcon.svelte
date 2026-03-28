<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Info' }: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 24 24"
	fill="currentColor"
	width={size}
	height={size}
	class={className}
	role="img"
	aria-hidden={title ? undefined : 'true'}
>
	{#if title}
		<title>{title}</title>
	{/if}
	<path
		fill-rule="evenodd"
		clip-rule="evenodd"
		d="M12 3.5a8.5 8.5 0 1 1 0 17 8.5 8.5 0 0 1 0-17Zm0 3a1 1 0 1 0 0 2 1 1 0 0 0 0-2Zm-1 4a1 1 0 1 1 2 0v5a1 1 0 1 1-2 0v-5Z"
	/>
</svg>
