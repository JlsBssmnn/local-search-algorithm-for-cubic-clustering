<script lang="ts">
	import { tick } from "svelte";
	import type { Evaluation, TestData } from "../wasmInterface/types";
	import { generateData, partitionData } from "../wasmInterface/functions";
	import Cube from "../geometry/cube";
	import Sidebar from "./Sidebar.svelte";
	import VisUtilities from "./VisUtilities.svelte";
	import type { VisualizationStatus } from "./types";
	import type { Vector } from "src/geometry/vector";
	import Visualization from "./Visualization.svelte";
	import { createRandomColorsWithoutWhite } from "../utils/createRandomColors";

	let testdata: TestData;
	let evaluation: Evaluation;
	let allIntersects: Vector[][] = [];
	let visualizationStatus: VisualizationStatus = "ready";
	let showPlanes: boolean = true;
	let colorizePoints: boolean = false;

	let points: any[] = [];
	let planes: any[] = [];
	$: data = [{
			type: "scatter3d",
			mode: "lines",
			x: [-1, 1, 1, -1, -1, -1, 1, 1, -1, -1, -1, -1, 1, 1, 1, 1],
			y: [1, 1, -1, -1, 1, 1, 1, -1, -1, 1, -1, -1, -1, -1, 1, 1],
			z: [-1, -1, -1, -1, -1, 1, 1, 1, 1, 1, 1, -1, -1, 1, 1, -1],
		}, ...points, ...planes];


	const cube = new Cube({ X: -1, Y: -1, Z: -1 }, 2, 2, 2);

	$: if (showPlanes) {
		displayVisPlanes();
	} else {
		planes = [];
	}

	$: {
		points = [];
		if (colorizePoints) {
			displayDataColorized();
		} else {
			displayDataOneColor();
		}
	}

	function clearPlanes() {
		clearVisPlanes();
		allIntersects = [];
		evaluation = null;
	}

	function clearVisPlanes() {
		planes = [];
	}

	function displayVisPlanes() {
		clearVisPlanes();

		if (allIntersects == null || allIntersects.length === 0) {
			return;
		}

		planes = allIntersects.map((intersect) => ({ 
				type: "mesh3d",
				x: intersect.map((vec) => vec.X),
				y: intersect.map((vec) => vec.Y),
				z: intersect.map((vec) => vec.Z),
				i: Array(intersect.length - 2).fill(0),
				j: [...Array(intersect.length - 2).keys()].map((i) => i + 1),
				k: [...Array(intersect.length - 2).keys()].map((i) => i + 2),
				alphahull: 5,
				opacity: 0.4,
				color: "rgb(200,100,300)",
		}));
	}

	function displayData(event) {
		clearPlanes();

		const { numOfPlanes, pointsPerPlane, mean, stddev } = event.detail;
		testdata = generateData(numOfPlanes, pointsPerPlane, mean, stddev);

		if (colorizePoints) {
			displayDataColorized();
		} else {
			displayDataOneColor();
		}
	}

	function displayDataOneColor() {
		if (testdata == null) {
			return;
		}
		points = [{
		x: testdata.Points.map((point) => point.X),
		y: testdata.Points.map((point) => point.Y),
		z: testdata.Points.map((point) => point.Z),
			mode: "markers",
			marker: {
				color: "rgb(0, 0, 225)",
				size: 12,
				line: {
					color: "rgba(217, 217, 217, 0.14)",
					width: 0.5,
				},
				opacity: 0.8,
			},
			type: "scatter3d",
		}]
	}

	function displayDataColorized() {
		if (testdata == null) {
			return;
		}
		const pointsPerPlane = testdata.Points.length / testdata.NumOfPlanes;
		const colors = createRandomColorsWithoutWhite(testdata.NumOfPlanes);

		const newPoints = new Array(testdata.NumOfPlanes);
		for (let i = 0; i < testdata.NumOfPlanes; i++) {
			const relevantPoints = testdata.Points.slice(
				i * pointsPerPlane,
				(i + 1) * pointsPerPlane
			);
			newPoints[i] = {
				x: relevantPoints.map((point) => point.X),
				y: relevantPoints.map((point) => point.Y),
				z: relevantPoints.map((point) => point.Z),
				mode: "markers",
				marker: {
					color: colors[i],
					size: 12,
					line: {
						color: "rgba(217, 217, 217, 0.14)",
						width: 0.5,
					},
					opacity: 0.8,
				},
				type: "scatter3d",
			};
		}
		points = newPoints;
	}

	async function calcPartition(event) {
		if (testdata == null) {
			 return;
		}
		visualizationStatus = "computing";
		const { algorithm, threshold, amplification } = event.detail;

		await tick();
		setTimeout(() => {
			evaluation = partitionData(algorithm, threshold, amplification);
			allIntersects = evaluation.ComputedPlanes.map((plane) =>
				cube.planeIntersectSorted(plane, 0)
			);
			visualizationStatus = "ready";

			if (showPlanes) {
				displayVisPlanes();
			}
		}, 10);
	}
</script>

<VisUtilities bind:visualizationStatus bind:showPlanes bind:colorizePoints />
<div class="container">
	<div>
		<Sidebar
			on:partitionEvent={calcPartition}
			on:display={displayData}
			{evaluation}
		/>
	</div>
	<Visualization {data} />
</div>

<style>
	.container {
		display: flex;
	}
</style>
