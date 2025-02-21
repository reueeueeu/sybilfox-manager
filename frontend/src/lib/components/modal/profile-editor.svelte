<script lang="ts">
	import { profile } from '$lib/wailsjs/go/models';
	import type { HTMLAttributes } from 'svelte/elements';
	import { Button } from '../ui/button';
	import * as Card from '../ui/card';
	import CardContent from '../ui/card/card-content.svelte';
	import { Checkbox } from '../ui/checkbox';
	import { Input } from '../ui/input';
	import { Label } from '../ui/label';
	import { Separator } from '../ui/separator';

	let {
		request,
		onSubmit
	}: {
		request: profile.Request;
		onSubmit: (config: profile.Request) => Promise<void>;
	} = $props();
	console.log(request);
	let editing: boolean = $derived(request.access_code > 0);
	let name: string = $state(request.name);
	let host: string = $state(request.proxy.Host);
	let user: string = $state(request.proxy.User);
	let password: string = $state(request.proxy.Password);
	$effect(() => {
		request.name = name;
		request.proxy.Host = host;
		request.proxy.User = user;
		request.proxy.Password = password;
	});
	let requesting = $state(false);
</script>

<Card.Root class="p-4">
	<Card.Title class="text-xl">
		{editing ? 'Edit Profile' : 'Add Profile'}
	</Card.Title>
	<Card.Description>Fingerprint will be generated using your data.</Card.Description>
	<CardContent class="mt-4 p-0">
		<form>
			<div class="grid gap-4">
				<div class="grid gap-2">
					<Label for="name">Name</Label>
					<Input id="name" bind:value={name} />
				</div>
				<div class="grid gap-2">
					<Label for="proxyAddress">Proxy Address</Label>
					<Input id="proxyAddress" bind:value={host} />
				</div>
				<div class="grid gap-2">
					<Label for="proxyUsername">Proxy Username</Label>
					<Input id="proxyUsername" bind:value={user} />
				</div>
				<div class="grid gap-2">
					<Label for="proxyPassword">Proxy Password</Label>
					<Input id="proxyPassword" bind:value={password} />
				</div>
				<Separator></Separator>
				<div class="flex items-center gap-2">
					<Checkbox checked={true} disabled></Checkbox>
					<div class="grid gap-1.5 leading-none">
						<Label
							class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
						>
							Enable GeoIP
						</Label>
						<p class="text-muted-foreground text-sm">Adjust fingerprint to your proxy location.</p>
					</div>
				</div>
				<Button
					disabled={requesting}
					onclick={() => {
						requesting = true;
						onSubmit(request).then(() => {
							requesting = false;
						});
					}}>{editing ? 'Save' : 'Generate'}</Button
				>
			</div>
		</form>
	</CardContent>
</Card.Root>
