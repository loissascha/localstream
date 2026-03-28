<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'User' }: Props = $props();
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
	<path d="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z" />
	<path d="M4.5 20.5a7.5 7.5 0 0 1 15 0 .5.5 0 0 1-.5.5H5a.5.5 0 0 1-.5-.5Z" />
</svg>
