<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Auth } from '$lib/wailsjs/go/main/App';
	import { BrowserOpenURL } from '$lib/wailsjs/runtime/runtime';
	
	import { toast } from "svelte-sonner";

	let {
		isAuthed = $bindable(),
		...props
	}: {isAuthed: boolean} = $props();
	let isAuthing = $state(false)
	let authCode = $state('')
</script>

<Card.Root class="mx-auto max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Auth</Card.Title>	
		<Card.Description>Enter access code to access interface</Card.Description>
	</Card.Header>
	<Card.Content>
		<form >
			<div class="grid gap-4">
				<div class="grid gap-2">
					<Label for="accessCode">Access Code</Label>
					<Input id="accessCode" bind:value={authCode} type="text" required />
				</div>
				<Button type="submit" class="w-full" disabled={isAuthing} onclick={async ()=>{
					isAuthing = true
					try {
						await Auth(Number(authCode))
						isAuthed = true
					} catch(e) {
						console.log(e)
						toast("Login error: " + e)
					}
					isAuthing = false;
				}}>Proceed</Button>
				<Button variant="outline" class="w-full"  onclick={()=>{
					BrowserOpenURL("https://t.me/sybilfox_bot")
				}}>Get Access Code</Button>
			</div>
		</form>
	</Card.Content>
</Card.Root>
