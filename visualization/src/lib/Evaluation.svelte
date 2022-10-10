<script lang="ts">
	import {fade} from 'svelte/transition';
	import type { Evaluation } from "src/wasmInterface/types";

	export let evaluation: Evaluation;
	export let numOfPlanes: number;
	$: ({
		Accuracy,
		TruePositives,
		FalsePositives,
		FalseNegatives,
		ComputedPlanes,
	} = evaluation);
	$: precision = TruePositives / (TruePositives + FalsePositives);
	$: recall = TruePositives / (TruePositives + FalseNegatives);
	$: f1 = (2 * precision * recall) / (precision + recall);

	interface Content {
		title: string;
		f: () => string;
	}

	$: content = [
		{
			title: "Accuracy",
			f: () => `${Math.round(Accuracy * 10000) / 100}%`,
		},
		{
			title: "Number of computed planes / actual planes",
			f: () => `${ComputedPlanes.length} / ${numOfPlanes}`,
		},
		{
			title: "Precision",
			f: () => {
				if (Object.is(precision, NaN)) {
					return "-";
				}
				return `${Math.round(precision * 1000) / 1000}`;
			},
		},
		{
			title: "Recall",
			f: () => {
				if (Object.is(recall, NaN)) {
					return "-";
				}
				return `${Math.round(recall * 1000) / 1000}`;
			},
		},
		{
			title: "F1-Score",
			f: () => {
				if (Object.is(f1, NaN)) {
					return "-";
				}
				return `${Math.round(f1 * 1000) / 1000}`;
			},
		},
	];
</script>

<div transition:fade>
	<div class="title">Evaluation</div>
	<div class="evaluation">
		{#each content as stat, i}
			<div class="key">{stat.title}:</div>
			<div class="value">{stat.f()}</div>
			{#if i < content.length - 1}
			<hr class="line" />
			{/if}
		{/each}
	</div>
</div>

<style>
	.evaluation {
		display: grid;
		grid-template-columns: 2fr 1fr;
		border: 1px solid black;
		border-radius: 20px;
		padding: 10px;
		column-gap: 5px;
	}
	.key {
		justify-self: start;
		align-self: center;
	}
	.value {
		justify-self: start;
		align-self: center;
	}
	.line {
		grid-column: 1/3;
		width: 100%;
	}
</style>
