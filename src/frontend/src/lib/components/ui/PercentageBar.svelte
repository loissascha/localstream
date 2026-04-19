<script lang="ts">
	interface Props {
		percentage: number;
	}
	let { percentage }: Props = $props();
</script>

<div class="bg-neutral-600">
	<div style={`width: ${percentage}%;`} class={`h-2 bg-brand text-sm`}></div>
</div>
