<script lang="ts">
	import DarkModel from '$lib/components/dark-model.svelte';
	import LoginForm from '$lib/components/login-form.svelte';
	import SidebarPage from '$lib/components/sidebar-page.svelte';
	import { IsAuthed } from '$lib/wailsjs/go/main/App';
	import { onMount } from 'svelte';
	import { on } from 'svelte/events';

	let isAuthed = $state(false);
	let stateLoaded = $state(false);
	onMount(async () => {
		isAuthed = await IsAuthed();
		stateLoaded = true;
	});
</script>

{#if stateLoaded}
	{#if !isAuthed}
		<div class="flex h-screen w-full items-center justify-center px-4">
			<div class="flex flex-col items-center gap-4">
				<LoginForm bind:isAuthed={isAuthed} />
				<DarkModel />
			</div>
		</div>
	{:else}
		<SidebarPage onLogout={()=>{isAuthed = false}}/>
	{/if}
{/if}
