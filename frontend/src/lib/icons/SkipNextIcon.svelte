<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Skip Next' }: Props = $props();
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
		d="M6.75 6.75c0-1.18 1.28-1.9 2.28-1.28l6.75 4.5a1.5 1.5 0 0 1 0 2.56l-6.75 4.5a1.5 1.5 0 0 1-2.28-1.28V6.75Z"
	/>
	<rect x="16.5" y="5.25" width="1.5" height="13.5" rx="0.75" />
</svg>
