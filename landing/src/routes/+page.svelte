<script>
	import { onMount } from 'svelte';

	// Define the general description and features, to be used in one of the outputs
	const generalDescription =
		'A simple CLI tool to chat with Gemini 2.5 Flash Lite, supporting text and multimodal file analysis directly from your terminal.';
	const generalFeatures = [
		'Analyze various file types (PDF, audio, images, video)',
		'Ask direct questions about files or general topics',
		'Get concise, direct answers, even extracting terminal commands to clipboard',
		'Easy API key setup',
		'Seamless integration with your terminal workflow'
	];

	// Function to format the general help output
	function getHelpOutput() {
		let output = '';
		output += `--- YoChat: Your Most Reachable LLM Interface ---\n\n`;
		output += `${generalDescription}\n\n`;
		output += `--- Key Features ---\n`;
		generalFeatures.forEach((feature) => {
			output += `- ${feature}\n`;
		});
		output += `\nFor more specific help, try: chat help`;
		return output;
	}

	// Array of commands with their respective simulated outputs
	const commandData = [
		{
			command: 'chat help',
			output: getHelpOutput()
		},
		{
			command: 'chat "How do i revert to a previous commit"',
			output:
				'To revert back to a previous commit, use the following command: \n\ngit revert <commit_hash> \nReplace `<commit_hash>` with the hash of the commit you want to revert to.\n\nExtracted commands copied to clipboard!'
		},
		{
			command: 'chat --file ./my_document.pdf "Summarize this document."',
			output:
				'Summary: This document details the architectural evolution of the Svelte framework, focusing on its compile-time approach and reactivity model.'
		},
		{
			command: 'chat --file ./image.jpg "ocr this image"',
			output:
				'Thirdpen.app\n\nthirdpen.app is an AI tool that helps you learn any concept interactively, visit thirdpen to get started. twitter link, facebook'
		},
		{
			command: 'chat "how do i compress a video with ffmpeg"',
			output:
				'Use this command - ffmpeg -i input.mp4 -vf scale=1280:-1 output.mp4 \n\nCommand copied to clipboard!'
		},

		{
			command: `chat "I'm bored"`,
			output: 'Tough luck dude!'
		}
	];

	let currentCommand = '';
	let currentOutputText = '';
	let commandIndex = 0; // Index of the character being typed for the current command
	let currentCommandDataIndex = 0; // Index of the command in the commandData array
	let typingDone = false;
	let showingOutput = false;

	const TYPING_SPEED = 50; // milliseconds per character
	const PAUSE_BEFORE_OUTPUT = 1000; // milliseconds before output appears
	const OUTPUT_DISPLAY_TIME = 3500; // milliseconds output is displayed
	const PAUSE_BEFORE_NEXT_COMMAND = 500; // milliseconds before next command starts typing

	async function typeAndDisplayCycle() {
		// Reset for the new command cycle
		typingDone = false;
		showingOutput = false;
		currentCommand = '';
		currentOutputText = '';
		commandIndex = 0;

		const currentData = commandData[currentCommandDataIndex];
		const commandToType = currentData.command;
		const outputToDisplay = currentData.output;

		// 1. Type the command
		while (commandIndex < commandToType.length) {
			currentCommand += commandToType.charAt(commandIndex);
			commandIndex++;
			await new Promise((resolve) => setTimeout(resolve, TYPING_SPEED));
		}
		typingDone = true; // Stop cursor blinking

		// 2. Pause before showing output
		await new Promise((resolve) => setTimeout(resolve, PAUSE_BEFORE_OUTPUT));
		showingOutput = true;
		currentOutputText = outputToDisplay;

		// 3. Display output for a duration
		await new Promise((resolve) => setTimeout(resolve, OUTPUT_DISPLAY_TIME));

		// 4. Increment to the next command in the array, looping if at the end
		currentCommandDataIndex = (currentCommandDataIndex + 1) % commandData.length;

		// 5. Pause before starting the next command (effectively clearing the screen for the user)
		await new Promise((resolve) => setTimeout(resolve, PAUSE_BEFORE_NEXT_COMMAND));

		// 6. Start the cycle again for the next command
		typeAndDisplayCycle();
	}

	onMount(() => {
		typeAndDisplayCycle(); // Start the command typing and display loop
	});
</script>

<svelte:head>
	<title>YouChat CLI Tool</title>
</svelte:head>

<div class="w-screen mx-auto flex flex-col items-center justify-center min-h-screen py-10">
	<div class="mb-12 flex flex-col gap-4 text-center">
		<h1 class="text-4xl font-bold text-white">YoChat</h1>
		<p class="text-lg text-gray-300">your most reachable llm interface</p>
	</div>
	<div class="terminal-container bg-stone-800">
		<div class="terminal-header">
			<div class="dot red"></div>
			<div class="dot yellow"></div>
			<div class="dot green"></div>
			<span class="title">yochat_cli - bash</span>
		</div>
		<pre class="terminal-body">
<span class="prompt">user@yochat:~# </span><span class="command-text">{currentCommand}</span><span
				class="cursor"
				class:blinking={!typingDone}>|</span
			>
{#if showingOutput}<span class="output-text"
					>{currentOutputText}
</span>
			{/if}
    </pre>
	</div>

	<div class="bg-stone-800 rounded-full divide-x-2 divide-stone-500 *:p-8 grid grid-cols-2 mt-20">
		<div class="flex items-center justify-center">
			<a href="https://github.com/Abdulmumin1/yochat">github</a>
		</div>
		<div class="flex items-center justify-center">
			<a href="">installation</a>
		</div>
	</div>
</div>

<style>
	/* Your existing CSS styles go here */
	:global(body) {
		background-color: #1a1a1a;
		display: flex;
		justify-content: center;
		align-items: center;
		margin: 0;
		font-family: 'Hack', 'Fira Code', 'Roboto Mono', monospace;
		color: #e0e0e0;
	}

	.terminal-container {
		border-radius: 8px;
		width: 90%;
		max-width: 800px;
		height: 500px;
		overflow: hidden;
		position: relative;
	}

	.terminal-header {
		padding: 8px 15px;
		display: flex;
		align-items: center;
		border-top-left-radius: 8px;
		border-top-right-radius: 8px;
	}

	.dot {
		width: 12px;
		height: 12px;
		border-radius: 50%;
		margin-right: 8px;
	}

	.red {
		background-color: #ff5f56;
	}
	.yellow {
		background-color: #ffbd2e;
	}
	.green {
		background-color: #27c93f;
	}

	.title {
		color: #c0c0c0;
		font-size: 0.9em;
		margin-left: auto;
		margin-right: auto;
	}

	.terminal-body {
		padding: 20px;
		min-height: 300px; /* Adjust as needed */
		white-space: pre-wrap; /* Preserve whitespace and wrap text */
		overflow-y: auto;
		font-size: 1em;
		line-height: 1.5;
	}

	.prompt {
		color: #7aff2e; /* Green for prompt */
	}

	.command-text {
		color: #e0e0e0; /* Light gray for command */
	}

	.cursor {
		display: inline-block;
		background-color: #e0e0e0;
		width: 8px;
		height: 1.2em; /* Matches line-height */
		vertical-align: middle;
		margin-left: 2px;
	}

	.cursor.blinking {
		animation: blink-caret 0.75s step-end infinite;
	}

	.output-text {
		color: #e0e0e0;
	}

	@keyframes blink-caret {
		from,
		to {
			background-color: transparent;
		}
		50% {
			background-color: #e0e0e0;
		}
	}
</style>
