<script lang="ts">
	import type { SettingItems } from "./types";
	import { createEventDispatcher } from "svelte";

	export let title: string;
	export let items: SettingItems;
	export let submissionText: string;
	export let submissionEventName: string;

	let dispatchAfterChange = false;

	const state = items.reduce(
		(prev, setting) =>
			Object.assign(prev, { [setting.field]: setting.defaultValue }),
		{}
	);
	const dispatch = createEventDispatcher();

	$: {
		if (dispatchAfterChange) {
			const _ = state;
			emitEvent();
		}
	}

	function emitEvent() {
		dispatch(submissionEventName, state);
	}
</script>

<div>
<div class="title">{title}</div>
<div class="settings">
	{#each items as item}
		{item.name}:&nbsp;
		{#if "minValue" in item}
		<input type="number" bind:value={state[item.field]} />
		<input
			type="range"
			bind:value={state[item.field]}
			min={item.minValue}
			max={item.maxValue}
			step={item.step}
		/>
		{:else}
		<select bind:value={state[item.field]} style="grid-column: 2 / 4">
			{#each item.options as option}
				<option value={option}>{option}</option>
			{/each}
		</select>
		{/if}
	{/each}
</div>
<div class="bot">
	<label
		title="Automatically update the visualization if you change the values?"
	>
		<input type="checkbox" bind:checked={dispatchAfterChange} />
		Auto update&nbsp;
	</label>
	<button on:click={emitEvent}>{submissionText}</button>
</div>
</div>

<style>
	input[type="number"], select {
		padding: 12px 0 12px 5px;
		margin: 4px 0;
		box-sizing: border-box;
	}

	.settings {
		display: grid;
		grid-template-columns: 2fr 1fr 2fr;
		align-items: center;
	}

	.bot {
		margin-top: 10px;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 30px;
	}

</style>
