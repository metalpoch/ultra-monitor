/** @type {import('tailwindcss').Config} */

export default {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
	theme: {
		colors: {
			gray: {
				'50': '#FFFFFF',
				'55': '#FCFCFC',
				'100': '#EFEFEF',
				'200': '#DCDCDC',
				'300': '#BDBDBD',
				'400': '#989898',
				'500': '#7C7C7C',
				'600': '#656565',
				'700': '#525252',
				'800': '#464646',
				'900': '#3D3D3D',
				'950': '#292929',
			},
			blue: {
				'50': '#F0F8FF',
				'100': '#DFF0FF',
				'200': '#CAE9FF',
				'300': '#79CAFF',
				'400': '#32B0FE',
				'500': '#0796F0',
				'600': '#0077CD',
				'700': '#005EA6',
				'800': '#035089',
				'900': '#094471',
				'950': '#062A4B',
			}
		},
		fontFamily: {
			kanit: ['Kanit', 'sans-serif']
		},
		extend: {}
	},
	plugins: [],
}
