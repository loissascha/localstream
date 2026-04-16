<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Collections' }: Props = $props();
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
	<rect x="3" y="5" width="13" height="14" rx="2" opacity="0.5" />
	<rect x="8" y="3" width="13" height="14" rx="2" />
	<rect x="10.5" y="8" width="8" height="1.5" rx="0.75" opacity="0.5" />
	<rect x="10.5" y="11" width="6" height="1.5" rx="0.75" opacity="0.5" />
</svg>
