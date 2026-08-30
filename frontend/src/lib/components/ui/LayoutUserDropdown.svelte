<script lang="ts">
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth.svelte';
	import { debug, showDebug } from '$lib/debug.svelte';
	import CheckIcon from '$lib/icons/CheckIcon.svelte';
	import LogoutIcon from '$lib/icons/LogoutIcon.svelte';
	import SettingsIcon from '$lib/icons/SettingsIcon.svelte';
	import UserIcon from '$lib/icons/UserIcon.svelte';
	import XIcon from '$lib/icons/XIcon.svelte';
	import DropdownItem from './DropdownItem.svelte';
	import DropdownMelt from './DropdownMelt.svelte';
</script>

<DropdownMelt>
	<UserIcon />
	{#snippet items()}
		{#if auth.isAdmin}
			<DropdownItem href={resolve('/(protected)/(admin)/admin')}>
				<div class="flex items-center gap-1">
					<SettingsIcon /> Settings
				</div>
			</DropdownItem>
			<DropdownItem
				onclick={() => {
					showDebug(!debug.showDebug);
				}}
			>
				<div class="flex items-center gap-1">
					{#if debug.showDebug}
						<CheckIcon />
					{:else}
						<XIcon />
					{/if}
					Debug
				</div>
			</DropdownItem>
		{/if}
		<DropdownItem href={resolve('/logout')}>
			<div class="flex items-center gap-1">
				<LogoutIcon /> Change Account
			</div>
		</DropdownItem>
	{/snippet}
</DropdownMelt>
